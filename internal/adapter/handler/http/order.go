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
func NewOrderHandler(OrderService port.OrderService) *OrderHandler {
	return &OrderHandler{
		svc: OrderService,
	}
}

// orderResponse represents an order response body
type orderResponse struct {
	ID           uint64                 `json:"id"`
	UserID       uint64                 `json:"user_id"`
	PaymentID    uint64                 `json:"payment_type_id"`
	CustomerName string                 `json:"customer_name"`
	TotalPrice   float64                `json:"total_price"`
	TotalPaid    float64                `json:"total_paid"`
	TotalReturn  float64                `json:"total_return"`
	ReceiptCode  string                 `json:"receipt_id"`
	Products     []orderProductResponse `json:"products"`
	PaymentType  paymentResponse        `json:"payment_type"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// orderProductResponse represents an order product response body
type orderProductResponse struct {
	ID               uint64          `json:"id"`
	OrderID          uint64          `json:"order_id"`
	ProductID        uint64          `json:"product_id"`
	Quantity         int64           `json:"qty"`
	Price            float64         `json:"price"`
	TotalNormalPrice float64         `json:"total_normal_price"`
	TotalFinalPrice  float64         `json:"total_final_price"`
	Product          productResponse `json:"product"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
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
	ProductID uint64 `json:"product_id" binding:"required,min=1"`
	Quantity  int64  `json:"qty" binding:"required,number"`
}

// createOrderRequest represents a request body for creating a new order
type createOrderRequest struct {
	PaymentID    uint64                `json:"payment_id" binding:"required"`
	CustomerName string                `json:"customer_name" binding:"required"`
	TotalPaid    int64                 `json:"total_paid" binding:"required"`
	Products     []orderProductRequest `json:"products" binding:"required"`
}

// CreateOrder creates a new order
func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req createOrderRequest
	var products []domain.OrderProduct

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	userID := 1 // TODO: get user ID from JWT

	for _, product := range req.Products {
		products = append(products, domain.OrderProduct{
			ProductID: product.ProductID,
			Quantity:  product.Quantity,
		})
	}

	order := domain.Order{
		UserID:       uint64(userID),
		PaymentID:    req.PaymentID,
		CustomerName: req.CustomerName,
		TotalPaid:    float64(req.TotalPaid),
		Products:     products,
	}

	_, err := oh.svc.CreateOrder(ctx, &order)
	if err != nil {
		if err.Error() == "Product stock is not enough" {
			errorResponse(ctx, http.StatusBadRequest, err)
			return
		}

		if err.Error() == "Total paid is less than total price" {
			errorResponse(ctx, http.StatusBadRequest, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newOrderResponse(&order)

	successResponse(ctx, http.StatusCreated, rsp)
}

// getOrderRequest represents a request body for retrieving an order
type getOrderRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// GetOrder gets an order by ID
func (oh *OrderHandler) GetOrder(ctx *gin.Context) {
	var req getOrderRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	order, err := oh.svc.GetOrder(ctx, req.ID)
	if err != nil {
		if err.Error() == "order not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newOrderResponse(order)

	successResponse(ctx, http.StatusOK, rsp)
}

// listOrdersRequest represents a request body for listing orders
type listOrdersRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0"`
	Limit uint64 `form:"limit" binding:"required,min=5"`
}

// ListOrders lists all orders with pagination
func (oh *OrderHandler) ListOrders(ctx *gin.Context) {
	var req listOrdersRequest
	var ordersList []orderResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	orders, err := oh.svc.ListOrders(ctx, req.Skip, req.Limit)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
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
