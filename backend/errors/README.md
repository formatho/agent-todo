# Error Handling

This document describes the custom error handling pattern used in the Agent Todo backend.

## Custom Error Types

We use custom error types to provide consistent, structured error responses across the API.

### AppError Structure

```go
type AppError struct {
    Code       int                    `json:"code"`
    Message    string                 `json:"message"`
    Details    map[string]interface{} `json:"details,omitempty"`
    Internal   error                  `json:"-"`
    HTTPStatus int                    `json:"-"`
}
```

### Using Custom Errors

#### In Handlers

```go
func (h *Handler) SomeMethod(c *gin.Context) {
    // Bad request with validation details
    if err := validate(req); err != nil {
        middleware.HandleError(c, apperrors.BadRequest("Invalid input").
            WithDetails(map[string]interface{}{
                "field": "email",
                "error": err.Error(),
            }))
        return
    }

    // Not found
    user, err := h.service.GetByID(id)
    if err != nil {
        middleware.HandleError(c, apperrors.NotFound("User not found"))
        return
    }

    // Internal server error with wrapped error
    if err := someOperation(); err != nil {
        middleware.HandleError(c, apperrors.InternalServerError("Operation failed").
            WithInternal(err))
        return
    }
}
```

### Predefined Errors

```go
// Authentication errors
apperrors.ErrInvalidCredentials  // 401 - Invalid email or password
apperrors.ErrUnauthorized        // 401 - Unauthorized
apperrors.ErrInvalidToken        // 401 - Invalid token
apperrors.ErrTokenExpired        // 401 - Token expired

// Not found errors
apperrors.ErrUserNotFound        // 404 - User not found
apperrors.ErrAgentNotFound       // 404 - Agent not found
apperrors.ErrTaskNotFound        // 404 - Task not found
apperrors.ErrProjectNotFound     // 404 - Project not found
```

### Creating Custom Errors

```go
// Simple error
err := apperrors.BadRequest("Invalid input")

// With details
err := apperrors.BadRequest("Invalid input").
    WithDetails(map[string]interface{}{
        "field": "email",
        "reason": "must be a valid email address",
    })

// With internal error (for logging)
err := apperrors.InternalServerError("Database error").
    WithInternal(dbErr)
```

### Error Response Format

```json
{
  "error": "Invalid input",
  "code": 400,
  "details": {
    "field": "email",
    "reason": "must be a valid email address"
  }
}
```

### Benefits

1. **Consistent error responses** across all endpoints
2. **Structured error details** for debugging
3. **Proper HTTP status codes** for each error type
4. **Internal error wrapping** for logging while hiding internals from clients
5. **Easy to use** helper functions

### Migration Guide

Replace old error handling:

```go
// Old
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

// New
middleware.HandleError(c, apperrors.BadRequest(err.Error()))
```

Replace common errors:

```go
// Old
c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})

// New
middleware.HandleError(c, apperrors.ErrUserNotFound)
```
