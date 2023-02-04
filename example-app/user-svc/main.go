package main

import (
	"github.com/DomBlack/ForkingGoRuntime/example-app/pkg/rest"
	"github.com/DomBlack/ForkingGoRuntime/example-app/user-svc/users"
)

func main() {
	srv := rest.NewServer("user", users.Port)

	rest.Post(srv, "/check-auth", CheckAuth)

	srv.Start()
}

func CheckAuth(_ *rest.Context, p *users.CheckAuthParams) (*users.User, error) {
	switch p.Token {
	case "secret":
		return &users.User{ID: 1}, nil

	default:
		return nil, rest.Unauthorized("invalid token")
	}
}
