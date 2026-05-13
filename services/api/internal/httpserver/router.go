package httpserver

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/nexus/api/internal/identity"
	githubintegration "github.com/nexus/api/internal/integrations/github"
	"github.com/nexus/api/internal/tenancy"
)

type Server struct {
	logger         *slog.Logger
	identity       *identity.Repository
	tenancy        *tenancy.Repository
	githubClient   *githubintegration.Client
	githubIndexer  *githubintegration.Indexer
	githubStore    *githubintegration.RepositoryStore
	allowedOrigins string
	sessionTTL     time.Duration
}

type healthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type registerRequest struct {
	Email            string `json:"email"`
	Name             string `json:"name"`
	Password         string `json:"password"`
	OrganizationName string `json:"organizationName"`
	WorkspaceName    string `json:"workspaceName"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Session      identity.Session         `json:"session"`
	WorkspaceSet tenancy.WorkspaceSummary `json:"workspaceSet"`
}

type ingestGitHubRepositoryRequest struct {
	WorkspaceID string `json:"workspaceId"`
	Owner       string `json:"owner"`
	Repository  string `json:"repository"`
	Token       string `json:"token"`
}

type syncGitHubFilesRequest struct {
	WorkspaceID  string `json:"workspaceId"`
	RepositoryID string `json:"repositoryId"`
	Token        string `json:"token"`
	MaxFiles     int    `json:"maxFiles"`
}

func NewRouter(logger *slog.Logger, identityRepository *identity.Repository, tenancyRepository *tenancy.Repository, githubClient *githubintegration.Client, githubIndexer *githubintegration.Indexer, githubStore *githubintegration.RepositoryStore, allowedOrigins string, sessionTTL time.Duration) http.Handler {
	server := &Server{
		logger:         logger,
		identity:       identityRepository,
		tenancy:        tenancyRepository,
		githubClient:   githubClient,
		githubIndexer:  githubIndexer,
		githubStore:    githubStore,
		allowedOrigins: allowedOrigins,
		sessionTTL:     sessionTTL,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", server.healthHandler)
	mux.HandleFunc("POST /v1/auth/register", server.registerHandler)
	mux.HandleFunc("POST /v1/auth/login", server.loginHandler)
	mux.HandleFunc("GET /v1/me", server.meHandler)
	mux.HandleFunc("GET /v1/workspaces", server.workspacesHandler)
	mux.HandleFunc("POST /v1/integrations/github/repositories", server.ingestGitHubRepositoryHandler)
	mux.HandleFunc("GET /v1/integrations/github/repositories", server.githubRepositoriesHandler)
	mux.HandleFunc("POST /v1/integrations/github/files/sync", server.syncGitHubFilesHandler)
	mux.HandleFunc("GET /v1/integrations/github/files", server.githubFilesHandler)

	return server.cors(server.requestLogger(mux))
}

func (server *Server) healthHandler(response http.ResponseWriter, request *http.Request) {
	writeJSON(response, http.StatusOK, healthResponse{
		Service: "api",
		Status:  "ok",
		Time:    time.Now().UTC().Format(time.RFC3339),
	})
}

func (server *Server) registerHandler(response http.ResponseWriter, request *http.Request) {
	var payload registerRequest
	if !decodeJSON(response, request, &payload) {
		return
	}

	if strings.TrimSpace(payload.Email) == "" || strings.TrimSpace(payload.Name) == "" || len(payload.Password) < 8 {
		writeError(response, http.StatusBadRequest, "email, name, and password with at least 8 characters are required")
		return
	}
	if strings.TrimSpace(payload.OrganizationName) == "" || strings.TrimSpace(payload.WorkspaceName) == "" {
		writeError(response, http.StatusBadRequest, "organizationName and workspaceName are required")
		return
	}

	session, err := server.identity.Register(request.Context(), payload.Email, payload.Name, payload.Password, server.sessionTTL)
	if err != nil {
		server.logger.Warn("user registration failed", slog.String("error", err.Error()))
		writeError(response, http.StatusConflict, "user could not be registered")
		return
	}

	workspaceSet, err := server.tenancy.CreateOrganizationWithWorkspace(request.Context(), session.User.ID, payload.OrganizationName, payload.WorkspaceName)
	if err != nil {
		server.logger.Error("workspace creation failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "workspace could not be created")
		return
	}

	writeJSON(response, http.StatusCreated, authResponse{Session: session, WorkspaceSet: workspaceSet})
}

func (server *Server) loginHandler(response http.ResponseWriter, request *http.Request) {
	var payload loginRequest
	if !decodeJSON(response, request, &payload) {
		return
	}

	session, err := server.identity.Login(request.Context(), payload.Email, payload.Password, server.sessionTTL)
	if errors.Is(err, identity.ErrInvalidCredentials) {
		writeError(response, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		server.logger.Error("login failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "login failed")
		return
	}

	workspaceSets, err := server.tenancy.ListForUser(request.Context(), session.User.ID)
	if err != nil {
		server.logger.Error("workspace lookup failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "workspace lookup failed")
		return
	}

	var firstWorkspaceSet tenancy.WorkspaceSummary
	if len(workspaceSets) > 0 {
		firstWorkspaceSet = workspaceSets[0]
	}

	writeJSON(response, http.StatusOK, authResponse{Session: session, WorkspaceSet: firstWorkspaceSet})
}

func (server *Server) meHandler(response http.ResponseWriter, request *http.Request) {
	user, ok := server.requireUser(response, request)
	if !ok {
		return
	}

	writeJSON(response, http.StatusOK, map[string]identity.User{"user": user})
}

func (server *Server) workspacesHandler(response http.ResponseWriter, request *http.Request) {
	user, ok := server.requireUser(response, request)
	if !ok {
		return
	}

	workspaceSets, err := server.tenancy.ListForUser(request.Context(), user.ID)
	if err != nil {
		server.logger.Error("workspace list failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "workspace list failed")
		return
	}

	writeJSON(response, http.StatusOK, map[string][]tenancy.WorkspaceSummary{"items": workspaceSets})
}

func (server *Server) ingestGitHubRepositoryHandler(response http.ResponseWriter, request *http.Request) {
	user, ok := server.requireUser(response, request)
	if !ok {
		return
	}

	var payload ingestGitHubRepositoryRequest
	if !decodeJSON(response, request, &payload) {
		return
	}

	if strings.TrimSpace(payload.WorkspaceID) == "" || strings.TrimSpace(payload.Owner) == "" || strings.TrimSpace(payload.Repository) == "" || strings.TrimSpace(payload.Token) == "" {
		writeError(response, http.StatusBadRequest, "workspaceId, owner, repository, and token are required")
		return
	}

	allowed, err := server.tenancy.UserCanAccessWorkspace(request.Context(), user.ID, payload.WorkspaceID)
	if err != nil {
		server.logger.Error("workspace access check failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "workspace access check failed")
		return
	}
	if !allowed {
		writeError(response, http.StatusForbidden, "workspace access denied")
		return
	}

	fetchedRepository, err := server.githubClient.FetchRepository(request.Context(), payload.Token, payload.Owner, payload.Repository)
	if err != nil {
		server.logger.Warn("github repository fetch failed", slog.String("error", err.Error()))
		writeError(response, http.StatusBadGateway, "github repository could not be fetched")
		return
	}

	savedRepository, err := server.githubStore.UpsertRepository(request.Context(), payload.WorkspaceID, fetchedRepository)
	if err != nil {
		server.logger.Error("github repository save failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "github repository could not be saved")
		return
	}

	writeJSON(response, http.StatusCreated, map[string]githubintegration.Repository{"repository": savedRepository})
}

func (server *Server) githubRepositoriesHandler(response http.ResponseWriter, request *http.Request) {
	user, ok := server.requireUser(response, request)
	if !ok {
		return
	}

	workspaceID := strings.TrimSpace(request.URL.Query().Get("workspaceId"))
	if workspaceID == "" {
		writeError(response, http.StatusBadRequest, "workspaceId is required")
		return
	}

	allowed, err := server.tenancy.UserCanAccessWorkspace(request.Context(), user.ID, workspaceID)
	if err != nil {
		server.logger.Error("workspace access check failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "workspace access check failed")
		return
	}
	if !allowed {
		writeError(response, http.StatusForbidden, "workspace access denied")
		return
	}

	repositories, err := server.githubStore.ListRepositories(request.Context(), workspaceID)
	if err != nil {
		server.logger.Error("github repository list failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "github repositories could not be listed")
		return
	}

	writeJSON(response, http.StatusOK, map[string][]githubintegration.Repository{"items": repositories})
}

func (server *Server) syncGitHubFilesHandler(response http.ResponseWriter, request *http.Request) {
	user, ok := server.requireUser(response, request)
	if !ok {
		return
	}

	var payload syncGitHubFilesRequest
	if !decodeJSON(response, request, &payload) {
		return
	}

	if strings.TrimSpace(payload.WorkspaceID) == "" || strings.TrimSpace(payload.RepositoryID) == "" || strings.TrimSpace(payload.Token) == "" {
		writeError(response, http.StatusBadRequest, "workspaceId, repositoryId, and token are required")
		return
	}

	repository, allowed := server.requireRepositoryAccess(response, request, user.ID, payload.WorkspaceID, payload.RepositoryID)
	if !allowed {
		return
	}

	result, err := server.githubIndexer.SyncRepositoryFiles(request.Context(), payload.Token, repository, payload.MaxFiles)
	if err != nil {
		server.logger.Warn("github file sync failed", slog.String("error", err.Error()))
		writeError(response, http.StatusBadGateway, "github files could not be synced")
		return
	}

	writeJSON(response, http.StatusCreated, map[string]githubintegration.FileSyncResult{"result": result})
}

func (server *Server) githubFilesHandler(response http.ResponseWriter, request *http.Request) {
	user, ok := server.requireUser(response, request)
	if !ok {
		return
	}

	workspaceID := strings.TrimSpace(request.URL.Query().Get("workspaceId"))
	repositoryID := strings.TrimSpace(request.URL.Query().Get("repositoryId"))
	if workspaceID == "" || repositoryID == "" {
		writeError(response, http.StatusBadRequest, "workspaceId and repositoryId are required")
		return
	}

	repository, allowed := server.requireRepositoryAccess(response, request, user.ID, workspaceID, repositoryID)
	if !allowed {
		return
	}

	files, err := server.githubStore.ListFiles(request.Context(), repository.ID, true, 100)
	if err != nil {
		server.logger.Error("github file list failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "github files could not be listed")
		return
	}

	writeJSON(response, http.StatusOK, map[string][]githubintegration.RepositoryFile{"items": files})
}

func (server *Server) requireRepositoryAccess(response http.ResponseWriter, request *http.Request, userID string, workspaceID string, repositoryID string) (githubintegration.Repository, bool) {
	allowed, err := server.tenancy.UserCanAccessWorkspace(request.Context(), userID, workspaceID)
	if err != nil {
		server.logger.Error("workspace access check failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "workspace access check failed")
		return githubintegration.Repository{}, false
	}
	if !allowed {
		writeError(response, http.StatusForbidden, "workspace access denied")
		return githubintegration.Repository{}, false
	}

	repository, err := server.githubStore.GetRepository(request.Context(), workspaceID, repositoryID)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(response, http.StatusNotFound, "repository not found")
		return githubintegration.Repository{}, false
	}
	if err != nil {
		server.logger.Error("github repository lookup failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "repository lookup failed")
		return githubintegration.Repository{}, false
	}

	return repository, true
}

func (server *Server) requireUser(response http.ResponseWriter, request *http.Request) (identity.User, bool) {
	token := bearerToken(request.Header.Get("Authorization"))
	if token == "" {
		writeError(response, http.StatusUnauthorized, "missing bearer token")
		return identity.User{}, false
	}

	user, err := server.identity.UserByToken(request.Context(), token)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(response, http.StatusUnauthorized, "invalid bearer token")
		return identity.User{}, false
	}
	if err != nil {
		server.logger.Error("session lookup failed", slog.String("error", err.Error()))
		writeError(response, http.StatusInternalServerError, "session lookup failed")
		return identity.User{}, false
	}

	return user, true
}

func bearerToken(value string) string {
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func (server *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		started := time.Now()
		next.ServeHTTP(response, request)
		server.logger.Info(
			"http request handled",
			slog.String("method", request.Method),
			slog.String("path", request.URL.Path),
			slog.Duration("duration", time.Since(started)),
		)
	})
}

func (server *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		if origin != "" && server.originAllowed(origin) {
			response.Header().Set("Access-Control-Allow-Origin", origin)
			response.Header().Set("Vary", "Origin")
		}

		response.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if request.Method == http.MethodOptions {
			response.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(response, request)
	})
}

func (server *Server) originAllowed(origin string) bool {
	for _, allowed := range strings.Split(server.allowedOrigins, ",") {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}

	return false
}

func decodeJSON(response http.ResponseWriter, request *http.Request, target any) bool {
	defer request.Body.Close()
	if err := json.NewDecoder(request.Body).Decode(target); err != nil {
		writeError(response, http.StatusBadRequest, "invalid json body")
		return false
	}

	return true
}

func writeJSON(response http.ResponseWriter, status int, body any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)

	if err := json.NewEncoder(response).Encode(body); err != nil {
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func writeError(response http.ResponseWriter, status int, message string) {
	writeJSON(response, status, errorResponse{Error: message})
}
