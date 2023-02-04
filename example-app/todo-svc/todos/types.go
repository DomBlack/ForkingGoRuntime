package todos

import (
	"context"
	"strconv"
	"time"

	"github.com/DomBlack/ForkingGoRuntime/example-app/pkg/rest"
)

const Port = 8082

type Todo struct {
	ID   int `json:"id,omitempty"`
	User int `json:"user,omitempty"`

	Title string `json:"title"`

	Created   time.Time `json:"created"`
	Completed bool      `json:"completed"`
}

func List(ctx context.Context, userID int) ([]*Todo, error) {
	return rest.DoGet[[]*Todo](ctx, Port, "by-user", strconv.Itoa(userID))
}

func Create(ctx context.Context, userID int, params *CreateParams) (*Todo, error) {
	return rest.DoPost[*Todo](ctx, Port, params, "by-user", strconv.Itoa(userID))
}

type CreateParams struct {
	Title string `json:"title"`
}

func Read(ctx context.Context, userID, todoID int) (*Todo, error) {
	return rest.DoGet[*Todo](ctx, Port, "by-user", strconv.Itoa(userID), strconv.Itoa(todoID))
}

type UpdateParams struct {
	Title     *string `json:"title,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

func Update(ctx context.Context, userID, todoID int, params *UpdateParams) (*Todo, error) {
	return rest.DoPatch[*Todo](ctx, Port, params, "by-user", strconv.Itoa(userID), strconv.Itoa(todoID))
}

func Delete(ctx context.Context, userID, todoID int) (*Todo, error) {
	return rest.DoDelete[*Todo](ctx, Port, "by-user", strconv.Itoa(userID), strconv.Itoa(todoID))
}
