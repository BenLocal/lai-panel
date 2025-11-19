package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func ErrorHandlerMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)

		// 检查是否有错误
		if len(c.Errors) == 0 {
			return
		}

		lastError := c.Errors.Last()
		if lastError == nil {
			return
		}

		statusCode := determineStatusCode(lastError)
		message := lastError.Error()
		if message == "" {
			message = "Internal server error"
		}

		c.JSON(http.StatusOK, ErrorResponse(statusCode, message))
		c.Abort()
	}
}

func determineStatusCode(err interface{}) int {
	// 尝试使用 IsType 方法检查错误类型
	// ErrorTypeBind = 1 << 0, ErrorTypePublic = 1 << 3
	if errorWithType, ok := err.(interface{ IsType(uint32) bool }); ok {
		if errorWithType.IsType(1<<0) || errorWithType.IsType(1<<3) {
			return http.StatusBadRequest
		}
	}

	var errorMsg string
	if errorWithMsg, ok := err.(interface{ Error() string }); ok {
		errorMsg = strings.ToLower(errorWithMsg.Error())
		if isBindingError(errorMsg) {
			return http.StatusBadRequest
		}
	}

	return http.StatusInternalServerError
}

func isBindingError(msg string) bool {
	keywords := []string{"bind", "validate", "invalid", "required"}
	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}
	return false
}
