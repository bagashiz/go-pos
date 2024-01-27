package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port/mock"
	"github.com/bagashiz/go-pos/internal/core/service"
	"github.com/bagashiz/go-pos/internal/core/util"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type createPaymentTestedInput struct {
	payment *domain.Payment
}

type createPaymentExpectedOutput struct {
	payment *domain.Payment
	err     error
}

func TestPaymentService_CreatePayment(t *testing.T) {
	ctx := context.Background()
	paymentID := gofakeit.Uint64()
	paymentName := gofakeit.CreditCardType()
	paymentType := domain.EWallet
	paymentLogo := gofakeit.ImageURL(320, 320)

	paymentInput := &domain.Payment{
		Name: paymentName,
		Type: paymentType,
		Logo: paymentLogo,
	}
	paymentOutput := &domain.Payment{
		ID:        paymentID,
		Name:      paymentName,
		Type:      paymentType,
		Logo:      paymentLogo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cacheKey := util.GenerateCacheKey("payment", paymentOutput.ID)
	paymentSerialized, _ := util.Serialize(paymentOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			paymentRepo *mock.MockPaymentRepository,
			cache *mock.MockCacheRepository,
		)
		input    createPaymentTestedInput
		expected createPaymentExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					CreatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(paymentOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("payments:*")).
					Return(nil)
			},
			input: createPaymentTestedInput{
				payment: paymentInput,
			},
			expected: createPaymentExpectedOutput{
				payment: paymentOutput,
				err:     nil,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					CreatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(nil, domain.ErrConflictingData)
			},
			input: createPaymentTestedInput{
				payment: paymentInput,
			},
			expected: createPaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrConflictingData,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					CreatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(nil, domain.ErrInternal)
			},
			input: createPaymentTestedInput{
				payment: paymentInput,
			},
			expected: createPaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					CreatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(paymentOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentSerialized), gomock.Eq(ttl)).
					Return(domain.ErrInternal)
			},
			input: createPaymentTestedInput{
				payment: paymentInput,
			},
			expected: createPaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					CreatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(paymentOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("payments:*")).
					Return(domain.ErrInternal)
			},
			input: createPaymentTestedInput{
				payment: paymentInput,
			},
			expected: createPaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			paymentRepo := mock.NewMockPaymentRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(paymentRepo, cache)

			paymentService := service.NewPaymentService(paymentRepo, cache)

			payment, err := paymentService.CreatePayment(ctx, tc.input.payment)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.payment, payment, "Payment mismatch")
		})
	}
}

type getPaymentTestedInput struct {
	id uint64
}

type getPaymentExpectedOutput struct {
	payment *domain.Payment
	err     error
}

func TestPaymentService_GetPayment(t *testing.T) {
	ctx := context.Background()
	paymentID := gofakeit.Uint64()
	paymentName := gofakeit.CreditCardType()
	paymentType := domain.EWallet
	paymentLogo := gofakeit.ImageURL(320, 320)

	paymentOutput := &domain.Payment{
		ID:   paymentID,
		Name: paymentName,
		Type: paymentType,
		Logo: paymentLogo,
	}

	cacheKey := util.GenerateCacheKey("payment", paymentOutput.ID)
	paymentSerialized, _ := util.Serialize(paymentOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			paymentRepo *mock.MockPaymentRepository,
			cache *mock.MockCacheRepository,
		)
		input    getPaymentTestedInput
		expected getPaymentExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(paymentSerialized, nil)
				err := util.Deserialize(paymentSerialized, &paymentOutput)
				if err != nil {
					return
				}
			},
			input: getPaymentTestedInput{
				id: paymentID,
			},
			expected: getPaymentExpectedOutput{
				payment: paymentOutput,
				err:     nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(paymentOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentSerialized), gomock.Eq(ttl)).
					Return(nil)
			},
			input: getPaymentTestedInput{
				id: paymentID,
			},
			expected: getPaymentExpectedOutput{
				payment: paymentOutput,
				err:     nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(nil, domain.ErrDataNotFound)
			},
			input: getPaymentTestedInput{
				id: paymentID,
			},
			expected: getPaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(nil, domain.ErrInternal)
			},
			input: getPaymentTestedInput{
				id: paymentID,
			},
			expected: getPaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return([]byte("invalid"), nil)
			},
			input: getPaymentTestedInput{
				id: paymentID,
			},
			expected: getPaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(paymentOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentSerialized), gomock.Eq(ttl)).
					Return(domain.ErrInternal)
			},
			input: getPaymentTestedInput{
				id: paymentID,
			},
			expected: getPaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		// tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			// TODO: fix race condition to enable parallel testing
			// t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			paymentRepo := mock.NewMockPaymentRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(paymentRepo, cache)

			paymentService := service.NewPaymentService(paymentRepo, cache)

			payment, err := paymentService.GetPayment(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.payment, payment, "Payment mismatch")
		})
	}
}

type listPaymentsTestedInput struct {
	skip  uint64
	limit uint64
}

type listPaymentsExpectedOutput struct {
	payments []domain.Payment
	err      error
}

func TestPaymentService_ListPayments(t *testing.T) {
	var payments []domain.Payment

	for i := 0; i < 10; i++ {
		payments = append(payments, domain.Payment{
			ID:   gofakeit.Uint64(),
			Name: gofakeit.CreditCardType(),
			Type: domain.EWallet,
			Logo: gofakeit.ImageURL(320, 320),
		})
	}

	ctx := context.Background()
	skip := gofakeit.Uint64()
	limit := gofakeit.Uint64()

	params := util.GenerateCacheKeyParams(skip, limit)
	cacheKey := util.GenerateCacheKey("payments", params)
	paymentsSerialized, _ := util.Serialize(payments)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			paymentRepo *mock.MockPaymentRepository,
			cache *mock.MockCacheRepository,
		)
		input    listPaymentsTestedInput
		expected listPaymentsExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(paymentsSerialized, nil)
			},
			input: listPaymentsTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listPaymentsExpectedOutput{
				payments: payments,
				err:      nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				paymentRepo.EXPECT().
					ListPayments(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Return(payments, nil)
				paymentsSerialized, _ := util.Serialize(payments)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentsSerialized), gomock.Eq(ttl)).
					Return(nil)
			},
			input: listPaymentsTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listPaymentsExpectedOutput{
				payments: payments,
				err:      nil,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return([]byte("invalid"), nil)
			},
			input: listPaymentsTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listPaymentsExpectedOutput{
				payments: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_internalError",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				paymentRepo.EXPECT().
					ListPayments(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Return(nil, domain.ErrInternal)
			},
			input: listPaymentsTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listPaymentsExpectedOutput{
				payments: nil,
				err:      domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				paymentRepo.EXPECT().
					ListPayments(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Return(payments, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentsSerialized), gomock.Eq(ttl)).
					Return(domain.ErrInternal)
			},
			input: listPaymentsTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listPaymentsExpectedOutput{
				payments: nil,
				err:      domain.ErrInternal,
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			paymentRepo := mock.NewMockPaymentRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(paymentRepo, cache)

			paymentService := service.NewPaymentService(paymentRepo, cache)

			payments, err := paymentService.ListPayments(ctx, tc.input.skip, tc.input.limit)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.payments, payments, "Payments mismatch")
		})
	}
}

type updatePaymentTestedInput struct {
	payment *domain.Payment
}

type updatePaymentExpectedOutput struct {
	payment *domain.Payment
	err     error
}

func TestPaymentService_UpdatePayment(t *testing.T) {
	ctx := context.Background()
	paymentID := gofakeit.Uint64()

	paymentInput := &domain.Payment{
		ID:   paymentID,
		Name: gofakeit.CreditCardType(),
		Type: domain.EWallet,
		Logo: gofakeit.ImageURL(320, 320),
	}
	paymentOutput := &domain.Payment{
		ID:   paymentID,
		Name: paymentInput.Name,
		Type: paymentInput.Type,
		Logo: paymentInput.Logo,
	}
	existingPayment := &domain.Payment{
		ID:   paymentID,
		Name: gofakeit.CreditCardType(),
		Type: domain.Cash,
		Logo: gofakeit.ImageURL(320, 320),
	}

	cacheKey := util.GenerateCacheKey("payment", paymentOutput.ID)
	paymentSerialized, _ := util.Serialize(paymentOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			paymentRepo *mock.MockPaymentRepository,
			cache *mock.MockCacheRepository,
		)
		input    updatePaymentTestedInput
		expected updatePaymentExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(existingPayment, nil)
				paymentRepo.EXPECT().
					UpdatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(paymentOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("payments:*")).
					Return(nil)
			},
			input: updatePaymentTestedInput{
				payment: paymentInput,
			},
			expected: updatePaymentExpectedOutput{
				payment: paymentOutput,
				err:     nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(nil, domain.ErrDataNotFound)
			},
			input: updatePaymentTestedInput{
				payment: paymentInput,
			},
			expected: updatePaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(nil, domain.ErrInternal)
			},
			input: updatePaymentTestedInput{
				payment: paymentInput,
			},
			expected: updatePaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_EmptyData",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(existingPayment, nil)
			},
			input: updatePaymentTestedInput{
				payment: existingPayment,
			},
			expected: updatePaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_SameData",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(existingPayment, nil)
			},
			input: updatePaymentTestedInput{
				payment: existingPayment,
			},
			expected: updatePaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(existingPayment, nil)
				paymentRepo.EXPECT().
					UpdatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(nil, domain.ErrConflictingData)
			},
			input: updatePaymentTestedInput{
				payment: paymentInput,
			},
			expected: updatePaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrConflictingData,
			},
		},
		{
			desc: "Fail_InternalErrorUpdate",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(existingPayment, nil)
				paymentRepo.EXPECT().
					UpdatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(nil, domain.ErrInternal)
			},
			input: updatePaymentTestedInput{
				payment: paymentInput,
			},
			expected: updatePaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(existingPayment, nil)
				paymentRepo.EXPECT().
					UpdatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(paymentOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(domain.ErrInternal)
			},
			input: updatePaymentTestedInput{
				payment: paymentInput,
			},
			expected: updatePaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(existingPayment, nil)
				paymentRepo.EXPECT().
					UpdatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(paymentOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentSerialized), gomock.Eq(ttl)).
					Return(domain.ErrInternal)
			},
			input: updatePaymentTestedInput{
				payment: paymentInput,
			},
			expected: updatePaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(existingPayment, nil)
				paymentRepo.EXPECT().
					UpdatePayment(gomock.Any(), gomock.Eq(paymentInput)).
					Return(paymentOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(paymentSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("payments:*")).
					Return(domain.ErrInternal)
			},
			input: updatePaymentTestedInput{
				payment: paymentInput,
			},
			expected: updatePaymentExpectedOutput{
				payment: nil,
				err:     domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			paymentRepo := mock.NewMockPaymentRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(paymentRepo, cache)

			paymentService := service.NewPaymentService(paymentRepo, cache)

			payment, err := paymentService.UpdatePayment(ctx, tc.input.payment)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.payment, payment, "Payment mismatch")
		})
	}
}

type deletePaymentTestedInput struct {
	id uint64
}

type deletePaymentExpectedOutput struct {
	err error
}

func TestPaymentService_DeletePayment(t *testing.T) {
	ctx := context.Background()
	paymentID := gofakeit.Uint64()

	cacheKey := util.GenerateCacheKey("payment", paymentID)

	testCases := []struct {
		desc  string
		mocks func(
			paymentRepo *mock.MockPaymentRepository,
			cache *mock.MockCacheRepository,
		)
		input    deletePaymentTestedInput
		expected deletePaymentExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(&domain.Payment{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("payments:*")).
					Return(nil)
				paymentRepo.EXPECT().
					DeletePayment(gomock.Any(), gomock.Eq(paymentID)).
					Return(nil)
			},
			input: deletePaymentTestedInput{
				id: paymentID,
			},
			expected: deletePaymentExpectedOutput{
				err: nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(nil, domain.ErrDataNotFound)
			},
			input: deletePaymentTestedInput{
				id: paymentID,
			},
			expected: deletePaymentExpectedOutput{
				err: domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(nil, domain.ErrInternal)
			},
			input: deletePaymentTestedInput{
				id: paymentID,
			},
			expected: deletePaymentExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(&domain.Payment{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(domain.ErrInternal)
			},
			input: deletePaymentTestedInput{
				id: paymentID,
			},
			expected: deletePaymentExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(&domain.Payment{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("payments:*")).
					Return(domain.ErrInternal)
			},
			input: deletePaymentTestedInput{
				id: paymentID,
			},
			expected: deletePaymentExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalErrorDelete",
			mocks: func(
				paymentRepo *mock.MockPaymentRepository,
				cache *mock.MockCacheRepository,
			) {
				paymentRepo.EXPECT().
					GetPaymentByID(gomock.Any(), gomock.Eq(paymentID)).
					Return(&domain.Payment{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("payments:*")).
					Return(nil)
				paymentRepo.EXPECT().
					DeletePayment(gomock.Any(), gomock.Eq(paymentID)).
					Return(domain.ErrInternal)
			},
			input: deletePaymentTestedInput{
				id: paymentID,
			},
			expected: deletePaymentExpectedOutput{
				err: domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			paymentRepo := mock.NewMockPaymentRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(paymentRepo, cache)

			paymentService := service.NewPaymentService(paymentRepo, cache)

			err := paymentService.DeletePayment(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
		})
	}
}
