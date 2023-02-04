package users

import (
	"context"

	"github.com/DomBlack/ForkingGoRuntime/example-app/pkg/rest"
)

const Port = 8081

type User struct {
	ID int `json:"id"`
}

type CheckAuthParams struct {
	Token string `json:"token"`
}

func CheckAuth(ctx context.Context, token string) (*User, error) {
	return rest.DoPost[*User](ctx, Port, &CheckAuthParams{Token: token}, "check-auth")
}
