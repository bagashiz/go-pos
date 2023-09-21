package handler

import (
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
//	@Success		200					{object}	orderResponse		"Order created"
//	@Failure		400					{object}	errorResponse		"Validation error"
//	@Failure		404					{object}	errorResponse		"Data not found error"
//	@Failure		409					{object}	errorResponse		"Data conflict error"
//	@Failure		500					{object}	errorResponse		"Internal server error"
//	@Router			/orders [post]
func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req createOrderRequest
	var products []domain.OrderProduct

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
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

	handleSuccess(ctx, rsp)
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
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/orders/{id} [get]
func (oh *OrderHandler) GetOrder(ctx *gin.Context) {
	var req getOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	order, err := oh.svc.GetOrder(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newOrderResponse(order)

	handleSuccess(ctx, rsp)
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
//	@Param			skip	query		uint64			true	"Skip records"
//	@Param			limit	query		uint64			true	"Limit records"
//	@Success		200		{object}	meta			"Orders displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		401		{object}	errorResponse	"Unauthorized error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/orders [get]
//	@Security		BearerAuth
func (oh *OrderHandler) ListOrders(ctx *gin.Context) {
	var req listOrdersRequest
	var ordersList []orderResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
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

	handleSuccess(ctx, rsp)
}
