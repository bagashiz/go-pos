package handler

import (
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// PaymentHandler represents the HTTP handler for payment-related requests
type PaymentHandler struct {
	svc port.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler instance
func NewPaymentHandler(svc port.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		svc,
	}
}

// createPaymentRequest represents a request body for creating a new payment
type createPaymentRequest struct {
	Name string             `json:"name" binding:"required" example:"Tunai"`
	Type domain.PaymentType `json:"type" binding:"required" example:"CASH"`
	Logo string             `json:"logo" binding:"omitempty,required" example:"https://example.com/cash.png"`
}

// CreatePayment godoc
//
//	@Summary		Create a new payment
//	@Description	create a new payment with name, type, and logo
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			createPaymentRequest	body		createPaymentRequest	true	"Create payment request"
//	@Success		200						{object}	paymentResponse			"Payment created"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/payments [post]
//	@Security		BearerAuth
func (ph *PaymentHandler) CreatePayment(ctx *gin.Context) {
	var req createPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	payment := domain.Payment{
		Name: req.Name,
		Type: req.Type,
		Logo: req.Logo,
	}

	_, err := ph.svc.CreatePayment(ctx, &payment)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newPaymentResponse(&payment)

	handleSuccess(ctx, rsp)
}

// getPaymentRequest represents a request body for retrieving a payment
type getPaymentRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetPayment godoc
//
//	@Summary		Get a payment
//	@Description	get a payment by id
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"Payment ID"
//	@Success		200	{object}	paymentResponse	"Payment retrieved"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/payments/{id} [get]
//	@Security		BearerAuth
func (ph *PaymentHandler) GetPayment(ctx *gin.Context) {
	var req getPaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	payment, err := ph.svc.GetPayment(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newPaymentResponse(payment)

	handleSuccess(ctx, rsp)
}

// listPaymentsRequest represents a request body for listing payments
type listPaymentsRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListPayments godoc
//
//	@Summary		List payments
//	@Description	List payments with pagination
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			true	"Skip"
//	@Param			limit	query		uint64			true	"Limit"
//	@Success		200		{object}	meta			"Payments displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/payments [get]
//	@Security		BearerAuth
func (ph *PaymentHandler) ListPayments(ctx *gin.Context) {
	var req listPaymentsRequest
	var paymentsList []paymentResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	payments, err := ph.svc.ListPayments(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, payment := range payments {
		paymentsList = append(paymentsList, newPaymentResponse(&payment))
	}

	total := uint64(len(paymentsList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, paymentsList, "payments")

	handleSuccess(ctx, rsp)
}

// updatePaymentRequest represents a request body for updating a payment
type updatePaymentRequest struct {
	Name string             `json:"name" binding:"omitempty,required" example:"Gopay"`
	Type domain.PaymentType `json:"type" binding:"omitempty,required,payment_type" example:"E-WALLET"`
	Logo string             `json:"logo" binding:"omitempty,required" example:"https://example.com/gopay.png"`
}

// UpdatePayment godoc
//
//	@Summary		Update a payment
//	@Description	update a payment's name, type, or logo by id
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			id						path		int						true	"Payment ID"
//	@Param			updatePaymentRequest	body		updatePaymentRequest	true	"Update payment request"
//	@Success		200						{object}	paymentResponse			"Payment updated"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/payments/{id} [put]
//	@Security		BearerAuth
func (ph *PaymentHandler) UpdatePayment(ctx *gin.Context) {
	var req updatePaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := stringToUint64(idStr)
	if err != nil {
		validationError(ctx, err)
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
		handleError(ctx, err)
		return
	}

	rsp := newPaymentResponse(&payment)

	handleSuccess(ctx, rsp)
}

// deletePaymentRequest represents a request body for deleting a payment
type deletePaymentRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeletePayment godoc
//
//	@Summary		Delete a payment
//	@Description	Delete a payment by id
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Payment ID"
//	@Success		200	{object}	response		"Payment deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/payments/{id} [delete]
//	@Security		BearerAuth
func (ph *PaymentHandler) DeletePayment(ctx *gin.Context) {
	var req deletePaymentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := ph.svc.DeletePayment(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
