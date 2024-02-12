FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

RUN apt-get update  \
    && apt-get install -y postgresql-client \
    && apt-get install -y redis-tools

# 提前切換使用者避免權限問題
USER vscode

RUN go install github.com/swaggo/swag/cmd/swag@latest