package identity

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) Register(ctx context.Context, email string, name string, password string, ttl time.Duration) (Session, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	cleanName := strings.TrimSpace(name)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Session{}, err
	}

	var user User
	err = repository.db.QueryRowContext(ctx, `
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING id::text, email, name, created_at`, cleanEmail, cleanName, string(passwordHash)).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
	)
	if err != nil {
		return Session{}, err
	}

	return repository.createSession(ctx, user, ttl)
}

func (repository *Repository) Login(ctx context.Context, email string, password string, ttl time.Duration) (Session, error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(email))

	var user User
	var passwordHash string
	err := repository.db.QueryRowContext(ctx, `
SELECT id::text, email, name, password_hash, created_at
FROM users
WHERE lower(email) = $1`, cleanEmail).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&passwordHash,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Session{}, ErrInvalidCredentials
	}
	if err != nil {
		return Session{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return Session{}, ErrInvalidCredentials
	}

	return repository.createSession(ctx, user, ttl)
}

func (repository *Repository) UserByToken(ctx context.Context, token string) (User, error) {
	tokenHash := hashToken(token)

	var user User
	err := repository.db.QueryRowContext(ctx, `
SELECT users.id::text, users.email, users.name, users.created_at
FROM sessions
JOIN users ON users.id = sessions.user_id
WHERE sessions.token_hash = $1 AND sessions.expires_at > now()`, tokenHash).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
	)
	return user, err
}

func (repository *Repository) createSession(ctx context.Context, user User, ttl time.Duration) (Session, error) {
	token, err := randomToken()
	if err != nil {
		return Session{}, err
	}
	expiresAt := time.Now().UTC().Add(ttl)

	_, err = repository.db.ExecContext(ctx, `
INSERT INTO sessions (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)`, user.ID, hashToken(token), expiresAt)
	if err != nil {
		return Session{}, err
	}

	return Session{Token: token, ExpiresAt: expiresAt, User: user}, nil
}

func randomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
