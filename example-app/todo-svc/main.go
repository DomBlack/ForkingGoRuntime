package main

import (
	"database/sql"
	"strconv"

	"github.com/DomBlack/ForkingGoRuntime/example-app/pkg/rest"
	"github.com/DomBlack/ForkingGoRuntime/example-app/todo-svc/todos"
	"github.com/rs/zerolog/log"
)

func main() {
	srv := rest.NewServer("todo", todos.Port)

	// Setup the handlers
	rest.Get(srv, "/by-user/:userID", ListTodos)
	rest.Post(srv, "/by-user/:userID", CreateTodo)

	rest.Get(srv, "/by-user/:userID/:todoID", ReadTodo)
	rest.Patch(srv, "/by-user/:userID/:todoID", UpdateTodo)
	rest.Delete(srv, "/by-user/:userID/:todoID", DeleteTodo)

	// Then connect to the database
	if err := connectToDB(); err != nil {
		srv.Log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Start listening
	srv.Start()
}

func ListTodos(ctx *rest.Context) ([]*todos.Todo, error) {
	// Get the user ID
	userID, err := userIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	// validate the request
	switch {
	case userID <= 0:
		return nil, rest.BadRequest("invalid userID")
	}

	// Query the database
	rows, err := database.Query("SELECT id, user_id, title, created, completed FROM todos WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Scan each row
	rtn := make([]*todos.Todo, 0)
	for rows.Next() {
		var todo todos.Todo
		if err := rows.Scan(&todo.ID, &todo.User, &todo.Title, &todo.Created, &todo.Completed); err != nil {
			log.Error().Err(err).Msg("failed to scan row")
			return nil, err
		}
		rtn = append(rtn, &todo)
	}

	// Check for errors
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("error during listing")
		return nil, err
	}

	return rtn, nil
}

func CreateTodo(ctx *rest.Context, p *todos.CreateParams) (*todos.Todo, error) {
	userID, err := userIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	// validate the request
	switch {
	case userID <= 0:
		return nil, rest.BadRequest("invalid userID")
	case p.Title == "":
		return nil, rest.BadRequest("missing title")
	}

	// Insert into the database
	var todo todos.Todo
	err = database.QueryRow(""+
		"INSERT INTO todos (user_id, title) VALUES ($1, $2) RETURNING id, user_id, title, created, completed",
		userID, p.Title,
	).Scan(&todo.ID, &todo.User, &todo.Title, &todo.Created, &todo.Completed)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert todo")
		return nil, err
	}

	return &todo, nil
}

func ReadTodo(ctx *rest.Context) (*todos.Todo, error) {
	userID, err := userIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	todoID, err := todoIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	// validate the request
	switch {
	case userID <= 0:
		return nil, rest.BadRequest("invalid userID")
	case todoID <= 0:
		return nil, rest.BadRequest("invalid todoID")
	}

	// Query the database
	var todo todos.Todo
	err = database.QueryRow(""+
		"SELECT id, user_id, title, created, completed FROM todos WHERE id = $1 AND user_id = $2",
		todoID, userID,
	).Scan(&todo.ID, &todo.User, &todo.Title, &todo.Created, &todo.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, rest.NotFound("todo not found")
		}

		log.Error().Err(err).Msg("failed to read todo")
		return nil, err
	}

	return &todo, nil
}

func UpdateTodo(ctx *rest.Context, p *todos.UpdateParams) (*todos.Todo, error) {
	userID, err := userIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	todoID, err := todoIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	// validate the request
	switch {
	case userID <= 0:
		return nil, rest.BadRequest("invalid userID")
	case todoID <= 0:
		return nil, rest.BadRequest("invalid todoID")
	case p.Title != nil && *p.Title == "":
		return nil, rest.BadRequest("empty title")
	}

	// Read the todo first
	todo, err := ReadTodo(ctx)
	if err != nil {
		return nil, err
	}

	// Patch the in memory todo with the changes
	if p.Title != nil {
		todo.Title = *p.Title
	}
	if p.Completed != nil {
		todo.Completed = *p.Completed
	}

	// Save the todo
	_, err = database.Exec(
		"UPDATE todos SET title = $1, completed = $2 WHERE id = $3 AND user_id = $4",
		todo.Title, todo.Completed, todoID, userID,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to update todo")
		return nil, err
	}

	return todo, nil
}

func DeleteTodo(ctx *rest.Context) (*todos.Todo, error) {
	userID, err := userIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	todoID, err := todoIDFromReq(ctx)
	if err != nil {
		return nil, err
	}

	// validate the request
	switch {
	case userID <= 0:
		return nil, rest.BadRequest("invalid userID")
	case todoID <= 0:
		return nil, rest.BadRequest("invalid todoID")
	}

	// Delete from the database
	var todo todos.Todo
	err = database.QueryRow(""+
		"DELETE FROM todos WHERE id = $1 AND user_id = $2 RETURNING id, user_id, title, created, completed",
		todoID, userID,
	).Scan(&todo.ID, &todo.User, &todo.Title, &todo.Created, &todo.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			// No row  means it's already deleted, so that's not an error
			return nil, nil
		}

		log.Error().Err(err).Msg("failed to delete todo")
		return nil, err
	}

	return &todo, nil
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

func userIDFromReq(ctx *rest.Context) (int, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return 0, rest.BadRequest("missing userID")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, rest.BadRequest("non-numeric userID")
	}

	return id, nil
}
