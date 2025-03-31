# Unit Tests Cases

## Handler 單元測試
使用MockService模擬API遇到的各種情境，以此對Handler做單元測試。
每個測試都需要遵循 Given-When-Then 的格式來撰寫。

### 通用測試設定

#### Response Body 讀取
需要確保能夠多次讀取response body，以便進行多重驗證。

#### 測試環境準備
- 設定Gin為測試模式
- 建立Mock Service
- 建立Handler實例

### 1. List Tasks API

#### Case 1: 成功取得任務列表
##### Given
- Mock Service 回傳一筆任務記錄
##### When
- 呼叫List Tasks API
##### Then
- 驗證回傳狀態碼為 200
- 驗證回傳資料包含預期的任務內容

### 2. Create Task API

#### Case 1: 成功建立任務
##### Given
- Mock Service 回傳新建立的任務
##### When
- 發送包含有效任務名稱的請求
##### Then
- 驗證回傳狀態碼為 200
- 驗證回傳資料包含新建立的任務資訊

#### Case 2: 建立任務失敗 - 空名稱
##### Given
- 不需 Mock Service（請求會在參數驗證階段失敗）
##### When
- 發送空任務名稱的請求
##### Then
- 驗證回傳狀態碼為 400
- 驗證錯誤訊息指出名稱為必填

### 3. Update Task API

#### Case 1: 成功更新任務
##### Given
- Mock Service 回傳更新後的任務
##### When
- 發送更新特定任務的請求
##### Then
- 驗證回傳狀態碼為 200
- 驗證回傳資料包含更新後的任務資訊

#### Case 2: 更新任務失敗 - 任務不存在
##### Given
- Mock Service 回傳任務不存在錯誤
##### When
- 嘗試更新不存在的任務
##### Then
- 驗證回傳狀態碼為 400
- 驗證錯誤訊息指出任務不存在

### 4. Delete Task API

#### Case 1: 成功刪除任務
##### Given
- Mock Service 回傳刪除成功
##### When
- 發送刪除特定任務的請求
##### Then
- 驗證回傳狀態碼為 204（無內容）

#### Case 2: 刪除任務失敗 - 任務不存在
##### Given
- Mock Service 回傳任務不存在錯誤
##### When
- 嘗試刪除不存在的任務
##### Then
- 驗證回傳狀態碼為 400
- 驗證錯誤訊息指出任務不存在

## TaskService 單元測試
對 TaskService 的內部邏輯和狀態管理進行測試。
每個測試都需要遵循 Given-When-Then 的格式來撰寫。

### 通用測試設定

#### 測試環境準備
- 建立新的 TaskService 實例
- 初始化必要的測試資料

### 1. ListTasks 方法

#### Case 1: 空任務列表
##### Given
- 新建立的 TaskService
##### When
- 呼叫 ListTasks 方法
##### Then
- 驗證回傳空陣列
- 驗證無錯誤發生

#### Case 2: 含有多筆任務
##### Given
- TaskService 中已存在多筆任務
##### When
- 呼叫 ListTasks 方法
##### Then
- 驗證回傳所有已存在的任務
- 驗證任務順序與建立順序一致
- 驗證無錯誤發生

### 2. CreateTask 方法

#### Case 1: 成功建立任務
##### Given
- 新建立的 TaskService
##### When
- 呼叫 CreateTask 方法並提供有效的任務名稱
##### Then
- 驗證回傳新建立的任務
- 驗證任務 ID 為遞增值
- 驗證任務狀態為預設值（0）
- 驗證無錯誤發生

#### Case 2: 建立任務失敗 - 空名稱
##### Given
- 新建立的 TaskService
##### When
- 呼叫 CreateTask 方法並提供空名稱
##### Then
- 驗證回傳適當的錯誤
- 驗證任務未被建立

### 3. UpdateTask 方法

#### Case 1: 成功更新任務
##### Given
- TaskService 中已存在特定任務
##### When
- 呼叫 UpdateTask 方法更新該任務
##### Then
- 驗證回傳更新後的任務
- 驗證任務資料已正確更新
- 驗證其他任務不受影響
- 驗證無錯誤發生

#### Case 2: 更新失敗 - 任務不存在
##### Given
- TaskService 中不存在目標任務
##### When
- 嘗試更新不存在的任務
##### Then
- 驗證回傳任務不存在錯誤

#### Case 3: 更新失敗 - 無效的狀態值
##### Given
- TaskService 中已存在特定任務
##### When
- 嘗試將任務狀態更新為無效值（非 0 或 1）
##### Then
- 驗證回傳無效狀態錯誤
- 驗證任務狀態未被更新

### 4. DeleteTask 方法

#### Case 1: 成功刪除任務
##### Given
- TaskService 中已存在特定任務
##### When
- 呼叫 DeleteTask 方法刪除該任務
##### Then
- 驗證無錯誤發生
- 驗證任務已被移除
- 驗證其他任務不受影響

#### Case 2: 刪除失敗 - 任務不存在
##### Given
- TaskService 中不存在目標任務
##### When
- 嘗試刪除不存在的任務
##### Then
- 驗證回傳任務不存在錯誤

### 5. 併發安全性測試

#### Case 1: 併發讀取
##### Given
- TaskService 中已存在多筆任務
##### When
- 多個 goroutine 同時呼叫 ListTasks 方法
##### Then
- 驗證所有讀取操作都能正確完成
- 驗證回傳資料一致性
- 驗證無錯誤發生

#### Case 2: 併發寫入
##### Given
- 新建立的 TaskService
##### When
- 多個 goroutine 同時建立任務
##### Then
- 驗證所有任務都被正確建立
- 驗證任務 ID 無重複
- 驗證無錯誤發生

#### Case 3: 併發讀寫
##### Given
- TaskService 中已存在部分任務
##### When
- 同時進行多個讀取和寫入操作
##### Then
- 驗證所有操作都能正確完成
- 驗證資料一致性
- 驗證無錯誤發生