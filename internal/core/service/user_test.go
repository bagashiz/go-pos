package service

import (
	"context"
	"errors"
	"testing"

	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/bagashiz/go-pos/mocks"
	"github.com/stretchr/testify/mock"
)

func TestUserService_DeleteUser(t *testing.T) {
	type fields struct {
		repo  port.UserRepository
		cache port.CacheService
	}
	type args struct {
		ctx context.Context
		id  uint64
	}
	type testStruct struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	tests := []testStruct{}
	DeleteUser_Error_GetUserByID := func() testStruct {
		repo := mocks.NewUserRepository(t)

		repo.
			On("GetUserByID", mock.Anything, mock.Anything).
			Return(nil, errors.New("dummy error"))

		return testStruct{
			name:    "if repo.GetUserByID error, should return error",
			fields:  fields{repo: repo},
			args:    args{},
			wantErr: true,
		}
	}

	DeleteUser_Error_DeleteByPrefix := func() testStruct {
		repo := mocks.NewUserRepository(t)
		cache := mocks.NewCacheService(t)

		repo.
			On("GetUserByID", mock.Anything, mock.Anything).
			Return(nil, nil)

		cache.
			On("Delete", mock.Anything, mock.Anything).
			Return(nil)

		cache.
			On("DeleteByPrefix", mock.Anything, "users:*").
			Return(errors.New("dummy error"))

		return testStruct{
			name:    "if cache.DeleteByPrefix error, should return error",
			fields:  fields{repo: repo, cache: cache},
			args:    args{},
			wantErr: true,
		}
	}

	DeleteUser_Error_DeleteUser := func() testStruct {
		repo := mocks.NewUserRepository(t)
		cache := mocks.NewCacheService(t)

		repo.
			On("GetUserByID", mock.Anything, mock.Anything).
			Return(nil, nil)

		cache.
			On("Delete", mock.Anything, mock.Anything).
			Return(nil)

		cache.
			On("DeleteByPrefix", mock.Anything, "users:*").
			Return(nil)

		repo.
			On("DeleteUser", mock.Anything, mock.Anything).
			Return(errors.New("dummy error"))

		return testStruct{
			name:    "if repo.DeleteUser error, should return error",
			fields:  fields{repo: repo, cache: cache},
			args:    args{},
			wantErr: true,
		}
	}

	DeleteUser_Success := func() testStruct {
		repo := mocks.NewUserRepository(t)
		cache := mocks.NewCacheService(t)

		repo.
			On("GetUserByID", mock.Anything, mock.Anything).
			Return(nil, nil)

		cache.
			On("Delete", mock.Anything, mock.Anything).
			Return(nil)

		cache.
			On("DeleteByPrefix", mock.Anything, "users:*").
			Return(nil)

		repo.
			On("DeleteUser", mock.Anything, mock.Anything).
			Return(nil)

		return testStruct{
			name:    "if no error, should return success",
			fields:  fields{repo: repo, cache: cache},
			args:    args{},
			wantErr: false,
		}
	}

	tests = append(tests, DeleteUser_Error_GetUserByID())
	tests = append(tests, DeleteUser_Error_DeleteByPrefix())
	tests = append(tests, DeleteUser_Error_DeleteUser())
	tests = append(tests, DeleteUser_Success())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				repo:  tt.fields.repo,
				cache: tt.fields.cache,
			}
			if err := us.DeleteUser(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserService.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
