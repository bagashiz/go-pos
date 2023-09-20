package handler

import (
	"net/http"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// OrderHandler represents the HTTP handler for order-related requests
type OrderHandler struct {
	svc port.OrderService
}

// NewOrderHandler creates a new OrderHandler instance
func NewOrderHandler(svc port.OrderService) *OrderHandler {
	return &OrderHandler{
		svc,
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

// orderProductRequest represents an order product request body
type orderProductRequest struct {
	ProductID uint64 `json:"product_id" binding:"required,min=1" example:"1"`
	Quantity  int64  `json:"qty" binding:"required,number" example:"1"`
}

// createOrderRequest represents a request body for creating a new order
type createOrderRequest struct {
	PaymentID    uint64                `json:"payment_id" binding:"required" example:"1"`
	CustomerName string                `json:"customer_name" binding:"required" example:"John Doe"`
	TotalPaid    int64                 `json:"total_paid" binding:"required" example:"100000"`
	Products     []orderProductRequest `json:"products" binding:"required"`
}

// CreateOrder godoc
//
//	@Summary		Create a new order
//	@Description	Create a new order and return the order data with purchase details
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			createOrderRequest	body		createOrderRequest	true	"Create order request"
//	@Success		201					{object}	orderResponse		"Order created"
//	@Failure		400					{object}	response			"Validation error"
//	@Failure		404					{object}	response			"Data not found error"
//	@Failure		409					{object}	response			"Data conflict error"
//	@Failure		500					{object}	response			"Internal server error"
//	@Router			/orders [post]
func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req createOrderRequest
	var products []domain.OrderProduct

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	for _, product := range req.Products {
		products = append(products, domain.OrderProduct{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		})
	}

	authPayload := getAuthPayload(ctx, authorizationPayloadKey)

	order := domain.Order{
		UserID:       authPayload.UserID,
		PaymentID:    req.PaymentID,
		CustomerName: req.CustomerName,
		TotalPaid:    float64(req.TotalPaid),
		Products:     products,
	}

	_, err := oh.svc.CreateOrder(ctx, &order)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newOrderResponse(&order)

	successResponse(ctx, http.StatusCreated, rsp)
}

// getOrderRequest represents a request body for retrieving an order
type getOrderRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetOrder godoc
//
//	@Summary		Get an order
//	@Description	Get an order by id and return the order data with purchase details
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Order ID"
//	@Success		200	{object}	orderResponse	"Order displayed"
//	@Failure		400	{object}	response		"Validation error"
//	@Failure		404	{object}	response		"Data not found error"
//	@Failure		500	{object}	response		"Internal server error"
//	@Router			/orders/{id} [get]
func (oh *OrderHandler) GetOrder(ctx *gin.Context) {
	var req getOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	order, err := oh.svc.GetOrder(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newOrderResponse(order)

	successResponse(ctx, http.StatusOK, rsp)
}

// listOrdersRequest represents a request body for listing orders
type listOrdersRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListOrders godoc
//
//	@Summary		List orders
//	@Description	List orders and return an array of order data with purchase details
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64		true	"Skip records"
//	@Param			limit	query		uint64		true	"Limit records"
//	@Success		200		{object}	response	"Orders displayed"
//	@Failure		400		{object}	response	"Validation error"
//	@Failure		401		{object}	response	"Unauthorized error"
//	@Failure		500		{object}	response	"Internal server error"
//	@Router			/orders [get]
//	@Security		BearerAuth
func (oh *OrderHandler) ListOrders(ctx *gin.Context) {
	var req listOrdersRequest
	var ordersList []orderResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	orders, err := oh.svc.ListOrders(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, order := range orders {
		ordersList = append(ordersList, newOrderResponse(&order))
	}

	total := uint64(len(ordersList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, ordersList, "orders")

	successResponse(ctx, http.StatusOK, rsp)
}
