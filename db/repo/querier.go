// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package repo

import (
	"context"
)

type Querier interface {
	CreateMatch(ctx context.Context, arg CreateMatchParams) (Match, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteExpiredTokens(ctx context.Context) error
	DeleteResetToken(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, arg ForgotPasswordParams) (PasswordReset, error)
	GetResetToken(ctx context.Context, token string) ([]PasswordReset, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id string) ([]User, error)
	UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error)
	UpdateUserProfile(ctx context.Context, arg UpdateUserProfileParams) (User, error)
}

var _ Querier = (*Queries)(nil)
