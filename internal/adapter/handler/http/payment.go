package handler

import (
	"net/http"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// PaymentHandler represents the HTTP handler for payment-related requests
type PaymentHandler struct {
	svc port.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler instance
func NewPaymentHandler(PaymentService port.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		svc: PaymentService,
	}
}

// paymentResponse represents a payment response body
type paymentResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Logo string `json:"logo"`
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

// createPaymentRequest represents a request body for creating a new payment
type createPaymentRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
	Logo string `json:"logo" binding:"omitempty,required"`
}

// CreatePayment creates a new payment
func (ph *PaymentHandler) CreatePayment(ctx *gin.Context) {
	var req createPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	payment := domain.Payment{
		Name: req.Name,
		Type: req.Type,
		Logo: req.Logo,
	}

	_, err := ph.svc.CreatePayment(ctx, &payment)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newPaymentResponse(&payment)

	successResponse(ctx, http.StatusCreated, rsp)
}

// getPaymentRequest represents a request body for retrieving a payment
type getPaymentRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// GetPayment retrieves a payment by id
func (ph *PaymentHandler) GetPayment(ctx *gin.Context) {
	var req getPaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	payment, err := ph.svc.GetPayment(ctx, req.ID)
	if err != nil {
		if err.Error() == "payment not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newPaymentResponse(payment)

	successResponse(ctx, http.StatusOK, rsp)
}

// listPaymentsRequest represents a request body for listing payments
type listPaymentsRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0"`
	Limit uint64 `form:"limit" binding:"required,min=5"`
}

// ListPayments lists all payments with pagination
func (ph *PaymentHandler) ListPayments(ctx *gin.Context) {
	var req listPaymentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	payments, err := ph.svc.ListPayments(ctx, req.Skip, req.Limit)
	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	paymentsList := make([]paymentResponse, 0)
	for _, payment := range payments {
		paymentsList = append(paymentsList, newPaymentResponse(payment))
	}

	total := uint64(len(paymentsList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, paymentsList, "payments")

	successResponse(ctx, http.StatusOK, rsp)
}

// updatePaymentRequest represents a request body for updating a payment
type updatePaymentRequest struct {
	Name string `json:"name" binding:"omitempty,required"`
	Type string `json:"type" binding:"omitempty,required"`
	Logo string `json:"logo" binding:"omitempty,required"`
}

// UpdatePayment updates a payment
func (ph *PaymentHandler) UpdatePayment(ctx *gin.Context) {
	var req updatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := convertStringToUint64(idStr)
	if err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	payment := domain.Payment{
		ID:   id,
		Name: req.Name,
		Type: req.Type,
		Logo: req.Logo,
	}

	_, err = ph.svc.UpdatePayment(ctx, &payment)
	if err != nil {
		if err.Error() == "payment not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	rsp := newPaymentResponse(&payment)

	successResponse(ctx, http.StatusOK, rsp)
}

// deletePaymentRequest represents a request body for deleting a payment
type deletePaymentRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// DeletePayment deletes a payment
func (ph *PaymentHandler) DeletePayment(ctx *gin.Context) {
	var req deletePaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := ph.svc.DeletePayment(ctx, req.ID)
	if err != nil {
		if err.Error() == "payment not found" {
			errorResponse(ctx, http.StatusNotFound, err)
			return
		}

		errorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	successResponse(ctx, http.StatusOK, nil)
}
