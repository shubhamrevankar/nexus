$ErrorActionPreference = "Stop"

Write-Host "Checking Go API..."
Push-Location services/api
go test ./...
Pop-Location

Write-Host "Checking Python AI service..."
Push-Location services/ai
python -m unittest discover -s tests
Pop-Location

Write-Host "Checking workspace if dependencies are installed..."
if (Test-Path node_modules) {
  pnpm lint
  pnpm typecheck
  pnpm test
} else {
  Write-Host "Skipping pnpm checks because node_modules is missing. Run pnpm install first."
}

