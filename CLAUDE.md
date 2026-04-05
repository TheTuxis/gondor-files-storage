# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview
Files & Storage microservice for Gondor platform. Go/Gin service handling file upload/download with S3/MinIO backend, folder hierarchy, versioning, and file sharing.

## Commands
- `make build` — compile to bin/server
- `make run` — run locally (needs PostgreSQL + MinIO/S3)
- `make test` — run all tests with race detector
- `make lint` — golangci-lint
- `make docker` — build Docker image

## Architecture
- `cmd/server/main.go` — entry point, dependency injection
- `internal/config/` — env-based configuration
- `internal/model/` — GORM models (File, Folder, FileVersion, FileShare)
- `internal/repository/` — database access layer
- `internal/service/` — business logic + S3 storage abstraction
- `internal/handler/` — HTTP handlers (Gin)
- `internal/middleware/` — JWT auth, logging

## Key Decisions
- S3-compatible storage (works with AWS S3 and MinIO)
- Files stored with generated S3 keys, original names in DB
- Folder hierarchy via self-referential parent_id
- File versioning tracked in separate table
- Soft delete for files and folders
- Max upload size configurable (default 100MB)
- All data scoped by company_id
