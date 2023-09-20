package handler

import (
	"net/http"
	"strconv"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// stringToUint64 is a helper function to convert a string to uint64
func stringToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)

	return num, err
}

// getAuthPayload is a helper function to get the auth payload from the context
func getAuthPayload(ctx *gin.Context, key string) *domain.TokenPayload {
	return ctx.MustGet(key).(*domain.TokenPayload)
}

// meta represents metadata for a paginated response
type meta struct {
	Total uint64 `json:"total"`
	Limit uint64 `json:"limit"`
	Skip  uint64 `json:"skip"`
}

// newMeta is a helper function to create metadata for a paginated response
func newMeta(total, limit, skip uint64) meta {
	return meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

// toMap is a helper function to add meta and data to a map
func toMap(m meta, data any, key string) map[string]any {
	return map[string]any{
		"meta": m,
		key:    data,
	}
}

// response represents a response body format
type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"Success || {Error message}"`
	Data    any    `json:"data,omitempty"`
}

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	port.ErrDataNotFound:        http.StatusNotFound,
	port.ErrConflictingData:     http.StatusConflict,
	port.ErrInvalidCredentials:  http.StatusUnauthorized,
	port.ErrNoUpdatedData:       http.StatusBadRequest,
	port.ErrInsufficientStock:   http.StatusBadRequest,
	port.ErrInsufficientPayment: http.StatusBadRequest,
}

// handleError determines the status code of an error and returns a JSON response with the error message and status code
func handleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errorResponse(ctx, statusCode, err)
}

// errorResponse sends an error response with the specified status code and error message
func errorResponse(ctx *gin.Context, statusCode int, err error) {
	rsp := response{
		Success: false,
		Message: err.Error(),
	}
	ctx.JSON(statusCode, rsp)
}

// abortResponse sends an error response and aborts the request with the specified status code and error message
func abortResponse(ctx *gin.Context, statusCode int, err error) {
	rsp := response{
		Success: false,
		Message: err.Error(),
	}
	ctx.AbortWithStatusJSON(statusCode, rsp)
}

// successResponse sends a success response with the specified status code and optional data
func successResponse(ctx *gin.Context, statusCode int, data interface{}) {
	rsp := response{
		Success: true,
		Message: "Success",
		Data:    data,
	}
	ctx.JSON(statusCode, rsp)
}
