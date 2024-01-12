package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// response represents a response body format
type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// meta represents metadata for a paginated response
type meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit uint64 `json:"limit" example:"10"`
	Skip  uint64 `json:"skip" example:"0"`
}

// newMeta is a helper function to create metadata for a paginated response
func newMeta(total, limit, skip uint64) meta {
	return meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

// authResponse represents an authentication response body
type authResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

// newAuthResponse is a helper function to create a response body for handling authentication data
func newAuthResponse(token string) authResponse {
	return authResponse{
		AccessToken: token,
	}
}

// userResponse represents a user response body
type userResponse struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// newUserResponse is a helper function to create a response body for handling user data
func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// paymentResponse represents a payment response body
type paymentResponse struct {
	ID   uint64             `json:"id" example:"1"`
	Name string             `json:"name" example:"Tunai"`
	Type domain.PaymentType `json:"type" example:"CASH"`
	Logo string             `json:"logo" example:"https://example.com/cash.png"`
}

// newPaymentResponse is a helper function to create a response body for handling payment data
func newPaymentResponse(payment *domain.Payment) paymentResponse {
	return paymentResponse{
		ID:   payment.ID,
		Name: payment.Name,
		Type: payment.Type,
		Logo: payment.Logo,
	}
}

// categoryResponse represents a category response body
type categoryResponse struct {
	ID   uint64 `json:"id" example:"1"`
	Name string `json:"name" example:"Foods"`
}

// newCategoryResponse is a helper function to create a response body for handling category data
func newCategoryResponse(category *domain.Category) categoryResponse {
	return categoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

// productResponse represents a product response body
type productResponse struct {
	ID        uint64           `json:"id" example:"1"`
	SKU       string           `json:"sku" example:"9a4c25d3-9786-492c-b084-85cb75c1ee3e"`
	Name      string           `json:"name" example:"Chiki Ball"`
	Stock     int64            `json:"stock" example:"100"`
	Price     float64          `json:"price" example:"5000"`
	Image     string           `json:"image" example:"https://example.com/chiki-ball.png"`
	Category  categoryResponse `json:"category"`
	CreatedAt time.Time        `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time        `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// newProductResponse is a helper function to create a response body for handling product data
func newProductResponse(product *domain.Product) productResponse {
	return productResponse{
		ID:        product.ID,
		SKU:       product.SKU.String(),
		Name:      product.Name,
		Stock:     product.Stock,
		Price:     product.Price,
		Image:     product.Image,
		Category:  newCategoryResponse(product.Category),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

// orderResponse represents an order response body
type orderResponse struct {
	ID           uint64                 `json:"id" example:"1"`
	UserID       uint64                 `json:"user_id" example:"1"`
	PaymentID    uint64                 `json:"payment_type_id" example:"1"`
	CustomerName string                 `json:"customer_name" example:"John Doe"`
	TotalPrice   float64                `json:"total_price" example:"100000"`
	TotalPaid    float64                `json:"total_paid" example:"100000"`
	TotalReturn  float64                `json:"total_return" example:"0"`
	ReceiptCode  string                 `json:"receipt_id" example:"4979cf6e-d215-4ff8-9d0d-b3e99bcc7750"`
	Products     []orderProductResponse `json:"products"`
	PaymentType  paymentResponse        `json:"payment_type"`
	CreatedAt    time.Time              `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt    time.Time              `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// newOrderResponse is a helper function to create a response body for handling order data
func newOrderResponse(order *domain.Order) orderResponse {
	return orderResponse{
		ID:           order.ID,
		UserID:       order.UserID,
		PaymentID:    order.PaymentID,
		CustomerName: order.CustomerName,
		TotalPrice:   order.TotalPrice,
		TotalPaid:    order.TotalPaid,
		TotalReturn:  order.TotalReturn,
		ReceiptCode:  order.ReceiptCode.String(),
		Products:     newOrderProductResponse(order.Products),
		PaymentType:  newPaymentResponse(order.Payment),
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}
}

// orderProductResponse represents an order product response body
type orderProductResponse struct {
	ID               uint64          `json:"id" example:"1"`
	OrderID          uint64          `json:"order_id" example:"1"`
	ProductID        uint64          `json:"product_id" example:"1"`
	Quantity         int64           `json:"qty" example:"1"`
	Price            float64         `json:"price" example:"100000"`
	TotalNormalPrice float64         `json:"total_normal_price" example:"100000"`
	TotalFinalPrice  float64         `json:"total_final_price" example:"100000"`
	Product          productResponse `json:"product"`
	CreatedAt        time.Time       `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt        time.Time       `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// newOrderProductResponse is a helper function to create a response body for handling order product data
func newOrderProductResponse(orderProduct []domain.OrderProduct) []orderProductResponse {
	var orderProductResponses []orderProductResponse

	for _, orderProduct := range orderProduct {
		orderProductResponses = append(orderProductResponses, orderProductResponse{
			ID:               orderProduct.ID,
			OrderID:          orderProduct.OrderID,
			ProductID:        orderProduct.ProductID,
			Quantity:         orderProduct.Quantity,
			Price:            orderProduct.Product.Price,
			TotalNormalPrice: orderProduct.TotalPrice,
			TotalFinalPrice:  orderProduct.TotalPrice,
			Product:          newProductResponse(orderProduct.Product),
			CreatedAt:        orderProduct.CreatedAt,
			UpdatedAt:        orderProduct.UpdatedAt,
		})
	}

	return orderProductResponses
}

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	port.ErrDataNotFound:               http.StatusNotFound,
	port.ErrConflictingData:            http.StatusConflict,
	port.ErrInvalidCredentials:         http.StatusUnauthorized,
	port.ErrUnauthorized:               http.StatusUnauthorized,
	port.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	port.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	port.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	port.ErrInvalidToken:               http.StatusUnauthorized,
	port.ErrExpiredToken:               http.StatusUnauthorized,
	port.ErrForbidden:                  http.StatusForbidden,
	port.ErrNoUpdatedData:              http.StatusBadRequest,
	port.ErrInsufficientStock:          http.StatusBadRequest,
	port.ErrInsufficientPayment:        http.StatusBadRequest,
}

// validationError sends an error response for some specific request validation error
func validationError(ctx *gin.Context, err error) {
	errMsgs := parseError(err)
	errRsp := newErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// handleError determines the status code of an error and returns a JSON response with the error message and status code
func handleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// handleAbort sends an error response and aborts the request with the specified status code and error message
func handleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

// parseError parses error messages from the error object and returns a slice of error messages
func parseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}
