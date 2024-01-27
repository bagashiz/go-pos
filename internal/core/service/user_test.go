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

type registerTestedInput struct {
	user *domain.User
}

type registerExpectedOutput struct {
	user *domain.User
	err  error
}

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()
	userName := gofakeit.Name()
	userEmail := gofakeit.Email()
	userPassword := gofakeit.Password(true, true, true, true, false, 8)
	hashedPassword, _ := util.HashPassword(userPassword)

	userInput := &domain.User{
		Name:     userName,
		Email:    userEmail,
		Password: userPassword,
	}
	userOutput := &domain.User{
		ID:        gofakeit.Uint64(),
		Name:      userName,
		Email:     userEmail,
		Password:  hashedPassword,
		Role:      domain.Cashier,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cacheKey := util.GenerateCacheKey("user", userOutput.ID)
	userSerialized, _ := util.Serialize(userOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock.MockUserRepository,
			cache *mock.MockCacheRepository,
		)
		input    registerTestedInput
		expected registerExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(nil)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(nil, domain.ErrInternal)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(nil, domain.ErrConflictingData)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  domain.ErrConflictingData,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(domain.ErrInternal)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(domain.ErrInternal)
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
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

			userRepo := mock.NewMockUserRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := service.NewUserService(userRepo, cache)

			user, err := userService.Register(ctx, tc.input.user)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.user, user, "User mismatch")
		})
	}
}

type getUserTestedInput struct {
	id uint64
}

type getUserExpectedOutput struct {
	user *domain.User
	err  error
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()
	userID := gofakeit.Uint64()
	userOutput := &domain.User{
		ID:       userID,
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, 8),
		Role:     domain.Cashier,
	}

	cacheKey := util.GenerateCacheKey("user", userID)
	userSerialized, _ := util.Serialize(userOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock.MockUserRepository,
			cache *mock.MockCacheRepository,
		)
		input    getUserTestedInput
		expected getUserExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(userSerialized, nil)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, domain.ErrDataNotFound)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: nil,
				err:  domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrInternal)
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, domain.ErrInternal)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(userOutput, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(domain.ErrInternal)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return([]byte("invalid"), nil)
			},
			input: getUserTestedInput{
				id: userID,
			},
			expected: getUserExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock.NewMockUserRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := service.NewUserService(userRepo, cache)

			user, err := userService.GetUser(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.user, user, "User mismatch")
		})
	}
}

type listUsersTestedInput struct {
	skip  uint64
	limit uint64
}

type listUsersExpectedOutput struct {
	users []domain.User
	err   error
}

func TestUserService_ListUsers(t *testing.T) {
	var users []domain.User

	for i := 0; i < 10; i++ {
		userPassword := gofakeit.Password(true, true, true, true, false, 8)
		hashedPassword, _ := util.HashPassword(userPassword)

		users = append(users, domain.User{
			ID:       gofakeit.Uint64(),
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: hashedPassword,
			Role:     domain.Cashier,
		})
	}

	ctx := context.Background()
	skip := gofakeit.Uint64()
	limit := gofakeit.Uint64()

	params := util.GenerateCacheKeyParams(skip, limit)
	cacheKey := util.GenerateCacheKey("users", params)
	usersSerialized, _ := util.Serialize(users)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock.MockUserRepository,
			cache *mock.MockCacheRepository,
		)
		input    listUsersTestedInput
		expected listUsersExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(usersSerialized, nil)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: users,
				err:   nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				userRepo.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Return(users, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(usersSerialized), gomock.Eq(ttl)).
					Return(nil)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: users,
				err:   nil,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return([]byte("invalid"), nil)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: nil,
				err:   domain.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				userRepo.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Return(nil, domain.ErrInternal)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: nil,
				err:   domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				cache.EXPECT().
					Get(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil, domain.ErrDataNotFound)
				userRepo.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(skip), gomock.Eq(limit)).
					Return(users, nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(usersSerialized), gomock.Eq(ttl)).
					Return(domain.ErrInternal)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: nil,
				err:   domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock.NewMockUserRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := service.NewUserService(userRepo, cache)

			users, err := userService.ListUsers(ctx, tc.input.skip, tc.input.limit)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.users, users, "Users mismatch")
		})
	}
}

type updateUserTestedInput struct {
	user *domain.User
}

type updateUserExpectedOutput struct {
	user *domain.User
	err  error
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	userID := gofakeit.Uint64()

	// TODO: test with hashed password

	userInput := &domain.User{
		ID:    userID,
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Role:  domain.Cashier,
	}
	userOutput := &domain.User{
		ID:    userID,
		Name:  userInput.Name,
		Email: userInput.Email,
		Role:  userInput.Role,
	}
	existingUser := &domain.User{
		ID:    userID,
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Role:  domain.Admin,
	}

	cacheKey := util.GenerateCacheKey("user", userID)
	userSerialized, _ := util.Serialize(userOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock.MockUserRepository,
			cache *mock.MockCacheRepository,
		)
		input    updateUserTestedInput
		expected updateUserExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(nil)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, domain.ErrDataNotFound)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, domain.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_EmptyData",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
			},
			input: updateUserTestedInput{
				user: &domain.User{
					ID: userID,
				},
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_SameData",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
			},
			input: updateUserTestedInput{
				user: existingUser,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ErrNoUpdatedData,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(nil, domain.ErrConflictingData)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ErrConflictingData,
			},
		},
		{
			desc: "Fail_InternalErrorUpdate",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(nil, domain.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(domain.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(domain.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(existingUser, nil)
				userRepo.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(userInput)).
					Return(userOutput, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					Set(gomock.Any(), gomock.Eq(cacheKey), gomock.Eq(userSerialized), gomock.Eq(ttl)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(domain.ErrInternal)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
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

			userRepo := mock.NewMockUserRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := service.NewUserService(userRepo, cache)

			user, err := userService.UpdateUser(ctx, tc.input.user)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.user, user, "User mismatch")
		})
	}
}

type deleteUserTestedInput struct {
	id uint64
}

type deleteUserExpectedOutput struct {
	err error
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	userID := gofakeit.Uint64()

	cacheKey := util.GenerateCacheKey("user", userID)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mock.MockUserRepository,
			cache *mock.MockCacheRepository,
		)
		input    deleteUserTestedInput
		expected deleteUserExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(&domain.User{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(nil)
				userRepo.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(userID)).
					Return(nil)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, domain.ErrDataNotFound)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: domain.ErrDataNotFound,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(nil, domain.ErrInternal)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(&domain.User{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(domain.ErrInternal)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(&domain.User{}, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(domain.ErrInternal)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
				err: domain.ErrInternal,
			},
		},
		{
			desc: "Fail_InternalErrorDelete",
			mocks: func(
				userRepo *mock.MockUserRepository,
				cache *mock.MockCacheRepository,
			) {
				user := &domain.User{
					ID: userID,
				}
				userRepo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(userID)).
					Return(user, nil)
				cache.EXPECT().
					Delete(gomock.Any(), gomock.Eq(cacheKey)).
					Return(nil)
				cache.EXPECT().
					DeleteByPrefix(gomock.Any(), gomock.Eq("users:*")).
					Return(nil)
				userRepo.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(userID)).
					Return(domain.ErrInternal)
			},
			input: deleteUserTestedInput{
				id: userID,
			},
			expected: deleteUserExpectedOutput{
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

			userRepo := mock.NewMockUserRepository(ctrl)
			cache := mock.NewMockCacheRepository(ctrl)

			tc.mocks(userRepo, cache)

			userService := service.NewUserService(userRepo, cache)

			err := userService.DeleteUser(ctx, tc.input.id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
		})
	}
}
