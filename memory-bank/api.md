# Task API Document
A task including name and status.

## 1. List tasks
List tasks.
### Maethod:
GET
### Path:
/tasks
### Response Status:
200: OK
### Response Payload:
| name   | type   | required | description                                         |
|--------|--------|----------|-----------------------------------------------------|
| id     | int    | true     | task's ID                                           |
| name   | string | true     | task's name                                         |
| status | int    | true     | enum[0,1], task's status; 0:incomplete, 1:completed |

## 2. Create a task
Create a task.
### Maethod:
POST
### Path:
/tasks
### Response Status:
200: OK
400: If the name field is empty.
### Request Payload:
| name | type   | required | description |
|------|--------|----------|-------------|
| name | string | true     | task name   |
### Response Payload:
| name   | type   | required | description                                         |
|--------|--------|----------|-----------------------------------------------------|
| id     | int    | true     | task's ID                                           |
| name   | string | true     | task's name                                         |
| status | int    | true     | enum[0,1], task's status; 0:incomplete, 1:completed |

## 3. Update a task
Update a task's name or status.
### Maethod:
PUT
### Path:
/tasks/{id}
### Path parameters:
| name | type | required | description |
|------|------|----------|-------------|
| id   | int  | true     | task's id   |
### Request Payload:
| name   | type   | required | description                                         |
|--------|--------|----------|-----------------------------------------------------|
| name   | string | true     | task's name                                         |
| status | int    | true     | enum[0,1], task's status; 0:incomplete, 1:completed |
### Response Status:
200: OK
400: If the name or status field is empty, the status is neither 0 nor 1, or the ID is not an integer.
### Response Payload:
| name   | type   | required | description                                         |
|--------|--------|----------|-----------------------------------------------------|
| id     | int    | true     | task's ID                                           |
| name   | string | true     | task's name                                         |
| status | int    | true     | enum[0,1], task's status; 0:incomplete, 1:completed |

## 3. Delete a task
Update a task's name or status.
### Maethod:
DELETE
### Path:
/tasks/{id}
### Path parameters:
| name | type | required | description |
|------|------|----------|-------------|
| id   | int  | true     | task's id   |
### Response Status:
204: No Content
400: If the name or status field is empty, the status is neither 0 nor 1, or the ID is not an integer.
### Response Payload:
| name   | type   | required | description                                         |
|--------|--------|----------|-----------------------------------------------------|
| id     | int    | true     | task's ID                                           |
| name   | string | true     | task's name                                         |
| status | int    | true     | enum[0,1], task's status; 0:incomplete, 1:completed |
