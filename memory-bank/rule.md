# 專案規範文件

本文件定義了專案的強制性規範。任何與規範衝突的修改都必須先提出討論。

## 1. 程式碼規範
必須嚴格遵循以下規範：
### 1.1. Go 程式碼風格
- Uber Go 風格指南：https://github.com/uber-go/guide/blob/master/style.md
- Google Go 決策指南：https://google.github.io/styleguide/go/decisions

### 1.2. Git 提交規範
- Angular 提交訊息格式：https://github.com/angular/angular/blob/22b96b9/CONTRIBUTING.md

## 2. 專案結構
必須遵循 Go 標準專案佈局：https://github.com/golang-standards/project-layout

## 3. 技術棧要求
基礎版本：Go 1.21
### 3.1. 必要依賴
以下依賴必須使用與 Go 1.21 相容的最新穩定版本：
- REST API 框架：github.com/gin-gonic/gin
- 依賴注入框架：go.uber.org/fx
- Mock 生成工具：github.com/golang/mock/mockgen

## 4. 程式碼倉庫
主倉庫位置：https://github.com/randallchang/ggl-be-demo-code

## 5. 系統分層規範

### 5.1. API 層 (internal/api/handler.go)
強制要求：
1. 結構定義：
   - 必須定義 `Handler` 結構體
   - 結構體方法僅處理 HTTP 請求的參數綁定和驗證
2. 方法簽名：
   - 所有方法必須且只能接收 `*gin.Context` 參數
   - 方法不允許有返回值
3. 路由配置：
   - 必須在同一 package 中創建 `routes.go`
   - 必須實現 `/ping` 健康檢查端點，返回 200 OK

### 5.2. 服務層 (internal/service/service.go)
強制要求：
1. 介面定義：
   - 必須定義 `Service` 介面
   - 必須實現 `TaskService` 結構體
2. 職責劃分：
   - 所有業務邏輯必須在此層實現
   - 必須使用 map 作為內存存儲
3. Mock 生成：
   - 必須添加 `go:generate` 註解
   - 必須在同目錄下生成 `mock/mock_service.go`
   - 生成的文件必須包含 `MockService` 結構體

## 6. 依賴注入規範

### 6.1. 介面注入規則
強制要求：
1. 注入方式：
   - 必須使用 `fx.Annotate` 結合 `fx.As` 進行介面注入
   - 禁止使用獨立的 provider 函數進行類型轉換
2. 標準範例：
   ```go
   fx.Provide(
     fx.Annotate(
       service.NewTaskService,
       fx.As(new(service.Service)),
     ),
   )
   ```
3. 實現目標：
   - 確保依賴注入的意圖明確可見
   - 避免不必要的類型轉換函數
   - 保持依賴注入的一致性

注意：本規範文件中的所有規則都是強制性的，除非得到明確的修改許可。

## 7. Makefile 規範

### 7.1. 必要指令
專案必須包含以下 Makefile 指令：

1. 基本操作：
   - `make init`：初始化專案（下載依賴、生成 mock）
   - `make build`：編譯專案
   - `make run`：運行專案
   - `make test`：運行所有測試
   - `make clean`：清理編譯產物和暫存檔案

2. Docker 相關：
   - `make docker-build`：構建 Docker 映像
   - `make docker-run`：運行 Docker 容器

3. 開發工具：
   - `make mock`：生成所有 mock 檔案
   - `make fmt`：格式化所有程式碼
   - `make lint`：運行程式碼檢查
   - `make tidy`：整理並更新依賴

### 7.2. 指令規範
1. 命名規則：
   - 使用小寫字母和連字符
   - 相關指令使用相同前綴（如 docker-*）

2. 文檔要求：
   - Makefile 必須包含每個指令的說明註解
   - 必須在 README.md 中列出所有可用的 make 指令

3. 執行要求：
   - 所有指令必須處理錯誤情況
   - 必須顯示執行進度和結果
   - 關鍵指令必須有檢查點（如檢查必要工具是否安裝）

### 7.3. 標準範例：
```makefile
.PHONY: init build run test clean docker-build docker-run mock fmt lint tidy

# 初始化專案
init:
	go mod download
	go generate ./...

# 編譯專案
build:
	go build -o app

# 運行專案
run:
	go run main.go

# 運行測試
test:
	go test -v ./...

# 清理編譯產物
clean:
	rm -f app
	go clean

# 構建 Docker 映像
docker-build:
	docker build -f build/Dockerfile -t task-service .

# 運行 Docker 容器
docker-run:
	docker run -p 8080:8080 task-service

# 生成 mock 檔案
mock:
	go generate ./...

# 格式化程式碼
fmt:
	go fmt ./...

# 運行程式碼檢查
lint:
	go vet ./...

# 整理依賴
tidy:
	go mod tidy
```

注意：本規範文件中的所有規則都是強制性的，除非得到明確的修改許可。

## 8. 程式碼檢查規範

### 8.1. 強制性檢查項目
1. Context 使用規範：
   - 所有需要 context.Context 的函數呼叫必須正確傳入
   - fx.New() 的 Start/Stop 方法必須傳入 context.Context
   - 禁止使用 context.TODO() 或 context.Background() 作為臨時解決方案

2. 錯誤處理規範：
   - 所有錯誤必須被處理或返回
   - 禁止使用 _ 忽略錯誤
   - 錯誤訊息必須具有描述性

3. 並發安全：
   - 共享資源必須使用互斥鎖保護
   - channel 操作必須考慮關閉情況
   - goroutine 必須有合適的退出機制

4. 資源管理：
   - 檔案操作必須正確關閉
   - 資料庫連接必須正確釋放
   - HTTP 響應體必須正確關閉

### 8.2. 常見錯誤檢查清單
1. fx 框架相關：
   ```go
   // 錯誤示例
   app.Start()  // 缺少 context 參數

   // 正確示例
   ctx := context.Background()
   app.Start(ctx)
   ```

2. HTTP 伺服器相關：
   ```go
   // 錯誤示例
   http.ListenAndServe(":8080", nil)  // 忽略錯誤返回

   // 正確示例
   if err := http.ListenAndServe(":8080", nil); err != nil {
       log.Fatal(err)
   }
   ```

3. 檔案操作相關：
   ```go
   // 錯誤示例
   file.Write(data)  // 忽略錯誤返回

   // 正確示例
   if _, err := file.Write(data); err != nil {
       return err
   }
   ```

### 8.3. 程式碼審查要點
1. 依賴注入：
   - 檢查 fx.New() 的配置是否完整
   - 驗證所有依賴是否正確注入
   - 確保生命週期鉤子正確設置

2. 錯誤處理：
   - 檢查是否有未處理的錯誤
   - 驗證錯誤訊息的準確性
   - 確保錯誤處理邏輯合理

3. 資源管理：
   - 檢查資源是否正確初始化
   - 驗證資源是否正確釋放
   - 確保沒有資源洩漏

注意：以上規範必須在程式碼審查時嚴格執行，任何違反規範的程式碼都不允許合併到主分支。