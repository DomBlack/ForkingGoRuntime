# Forking the Go Runtime

This repo will hold the code for the talk.


## Updating the go-src root

```bash
git subtree pull --prefix go-src https://go.googlesource.com/go release-branch.go1.20 --squash
```

## Example app

The example app is a simple todo app. It contains three services:

- `todo-svc` - The todo service is responsible for managing the todos for users
- `user-svc` - The user service is responsible for authenticating users. In this example it's hardcoded to only allow
  one user with bearer token `secret`.
- `api-svc` - The API service acts the the gateway for the user's requests. It is responsible for authenticating the user
  against the `user-svc` and then forwarding the request to the `todo-svc`.

### Running the example app

To run the apps you will need to brew install `make`, `postgres` and `overmind`.

```bash
make initdb        # Create the database
make postgres &    # Start the database
make microservices # Start the microservices
```

### Example API calls

```bash
# List todos for user 1
curl -H "Authorization: Bearer secret" http://localhost:8080/todos

# Create todo for user 1
curl -H "Authorization: Bearer secret" http://localhost:8080/todos -X "POST" -d `{"title":"My Todo"}`

# Read the first todo
curl -H "Authorization: Bearer secret" http://localhost:8080/todos/1

# Update title for the first todo
curl -H "Authorization: Bearer secret" http://localhost:8080/todos/1 -X "PATCH" -d `{"title":"New title"}`

# Update completed for the first todo
curl -H "Authorization: Bearer secret" http://localhost:8080/todos/1 -X "PATCH" -d `{"completed":true}`

# Delete the first todo
curl -H "Authorization: Bearer secret" http://localhost:8080/todos/1 -X "DELETE"

```
