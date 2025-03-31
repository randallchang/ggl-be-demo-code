# Build Requirements

## Docker Build Specifications

### 1. Dockerfile Location and Naming
- 必須在 `build/Dockerfile` 位置創建 Dockerfile
- 禁止在其他位置創建 Dockerfile

### 2. Base Image Requirements
- 必須使用官方 golang alpine 鏡像作為基礎鏡像
- 版本必須與 go.mod 中指定的 Go 版本一致
- 禁止使用 latest 標籤

### 3. Multi-stage Build Requirements
必須使用多階段構建，包含以下階段：
1. Builder 階段：
   - 用於編譯 Go 程式
   - 必須設置 GOOS=linux 和 GOARCH=amd64
   - 必須啟用 CGO_ENABLED=0
   - 必須使用 `-ldflags="-s -w"` 減小二進制檔案大小

2. Final 階段：
   - 必須使用 alpine:latest 作為基礎鏡像
   - 必須創建非 root 用戶運行應用
   - 必須設置適當的工作目錄
   - 必須複製編譯後的二進制檔案
   - 必須設置 EXPOSE 8080

### 4. 安全要求
- 必須設置非 root 用戶
- 必須移除不必要的系統工具和檔案
- 必須設置適當的檔案權限

### 5. 標準 Dockerfile 範例
```dockerfile
# Builder stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/server ./cmd/main.go

# Final stage
FROM alpine:latest
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /app/server .
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["./server"]
```

### 6. 構建指令
必須在 Makefile 中包含以下 Docker 相關指令：
```makefile
.PHONY: docker-build docker-run

# Build Docker image
docker-build:
	docker build -f build/Dockerfile -t task-service .

# Run Docker container
docker-run:
	docker run -p 8080:8080 task-service
```

### 7. 自動化檢查清單
每次構建前必須檢查：
1. Dockerfile 是否位於正確位置
2. 基礎鏡像版本是否與 go.mod 一致
3. 是否包含所有必要的安全設置
4. 是否正確設置了多階段構建
5. 是否包含必要的 EXPOSE 指令