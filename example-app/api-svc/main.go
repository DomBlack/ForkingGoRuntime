package main

import (
	"strconv"
	"strings"

	"github.com/DomBlack/ForkingGoRuntime/example-app/pkg/rest"
	"github.com/DomBlack/ForkingGoRuntime/example-app/todo-svc/todos"
	"github.com/DomBlack/ForkingGoRuntime/example-app/user-svc/users"
)

func main() {
	srv := rest.NewServer("api", 8080)

	rest.Get(srv, "/todos", ListTodos)
	rest.Post(srv, "/todos", CreateTodo)
	rest.Get(srv, "/todos/:todoID", ReadTodo)
	rest.Patch(srv, "/todos/:todoID", UpdateTodo)
	rest.Delete(srv, "/todos/:todoID", DeleteTodo)

	srv.Start()
}

func ListTodos(ctx *rest.Context) ([]*todos.Todo, error) {
	userID, err := authedUser(ctx)
	if err != nil {
		return nil, err
	}

	return todos.List(ctx, userID)
}

func CreateTodo(ctx *rest.Context, req *todos.CreateParams) (*todos.Todo, error) {
	userID, err := authedUser(ctx)
	if err != nil {
		return nil, err
	}

	return todos.Create(ctx, userID, req)
}

func ReadTodo(ctx *rest.Context) (*todos.Todo, error) {
	userID, err := authedUser(ctx)
	if err != nil {
		return nil, err
	}

	todoID, err := todoIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	return todos.Read(ctx, userID, todoID)
}

func UpdateTodo(ctx *rest.Context, p *todos.UpdateParams) (*todos.Todo, error) {
	userID, err := authedUser(ctx)
	if err != nil {
		return nil, err
	}

	todoID, err := todoIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	return todos.Update(ctx, userID, todoID, p)
}

func DeleteTodo(ctx *rest.Context) (*todos.Todo, error) {
	userID, err := authedUser(ctx)
	if err != nil {
		return nil, err
	}

	todoID, err := todoIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	return todos.Delete(ctx, userID, todoID)
}

// authedUser returns the user ID of the user who is authenticated
func authedUser(ctx *rest.Context) (int, error) {
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		return 0, rest.Unauthorized("no auth header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, rest.Unauthorized("only bearer auth supported")
	}

	user, err := users.CheckAuth(ctx, authHeader[len("Bearer "):])
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func todoIDFromReq(ctx *rest.Context) (int, error) {
	userID := ctx.Param("todoID")
	if userID == "" {
		return 0, rest.BadRequest("missing todoID")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, rest.BadRequest("non-numeric todoID")
	}

	return id, nil
}
