package github

import (
	"context"
	"path/filepath"
	"strings"
)

const defaultMaxFiles = 50
const maxFileSizeBytes = 100 * 1024

type Indexer struct {
	client *Client
	store  *RepositoryStore
}

func NewIndexer(client *Client, store *RepositoryStore) *Indexer {
	return &Indexer{client: client, store: store}
}

func (indexer *Indexer) SyncRepositoryFiles(ctx context.Context, token string, repository Repository, maxFiles int) (FileSyncResult, error) {
	if maxFiles <= 0 || maxFiles > defaultMaxFiles {
		maxFiles = defaultMaxFiles
	}

	treeFiles, err := indexer.client.FetchTreeFiles(ctx, token, repository.Owner, repository.Name, repository.DefaultBranch)
	if err != nil {
		return FileSyncResult{}, err
	}

	result := FileSyncResult{RepositoryID: repository.ID, Files: make([]RepositoryFile, 0)}
	indexedCount := 0

	for _, treeFile := range treeFiles {
		if indexedCount >= maxFiles {
			break
		}

		file := RepositoryFile{RepositoryID: repository.ID, Path: treeFile.Path, SHA: treeFile.SHA, Size: treeFile.Size}
		if reason := skipReason(treeFile); reason != "" {
			file.SkippedReason = reason
			result.SkippedCount++
		} else {
			content, err := indexer.client.FetchRawFile(ctx, token, repository.Owner, repository.Name, treeFile.Path, repository.DefaultBranch)
			if err != nil {
				file.SkippedReason = "fetch_failed"
				result.SkippedCount++
			} else {
				file.ContentText = content
				file.Indexed = true
				indexedCount++
				result.IndexedCount++
			}
		}

		savedFile, err := indexer.store.UpsertFile(ctx, repository.ID, file)
		if err != nil {
			return FileSyncResult{}, err
		}
		result.Files = append(result.Files, savedFile)
	}

	return result, nil
}

func skipReason(file TreeFile) string {
	path := strings.ToLower(file.Path)
	base := strings.ToLower(filepath.Base(path))
	extension := strings.ToLower(filepath.Ext(path))

	if file.Size <= 0 {
		return "empty_file"
	}
	if file.Size > maxFileSizeBytes {
		return "file_too_large"
	}
	if strings.Contains(path, "node_modules/") || strings.Contains(path, "vendor/") || strings.Contains(path, "dist/") || strings.Contains(path, "build/") {
		return "generated_or_dependency_path"
	}
	if strings.Contains(path, ".env") || strings.Contains(path, "secret") || strings.Contains(path, "credential") || strings.Contains(path, "private_key") {
		return "sensitive_path"
	}
	if strings.HasSuffix(base, ".lock") || base == "package-lock.json" || base == "pnpm-lock.yaml" || base == "yarn.lock" {
		return "lockfile"
	}
	if strings.HasSuffix(base, ".min.js") || strings.HasSuffix(base, ".min.css") {
		return "minified_asset"
	}
	if binaryOrMediaExtension(extension) {
		return "binary_or_media_file"
	}

	return ""
}

func binaryOrMediaExtension(extension string) bool {
	switch extension {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf", ".zip", ".gz", ".tar", ".mp4", ".mov", ".mp3", ".wav", ".exe", ".dll", ".so", ".dylib", ".class", ".jar", ".woff", ".woff2", ".ttf", ".otf":
		return true
	default:
		return false
	}
}
