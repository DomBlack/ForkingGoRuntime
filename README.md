# The adventurer's guide to forking the go runtime

This repo contains the code and [sides](./slides.pdf) for a talk I gave at GopherCon UK 2023 on how by creating
a rolling fork of the language we use every day can massively improve our experience as developers.

This talk is loosely based on how [Encore](https://encore.dev) uses a rolling fork of the Go runtime to add automatic tracing and unit test isolation to applications built using Encore without the developers of those applications having to add anything to their code bases. You can checkout [Encore Go Rolling Fork](https://github.com/encoredev/go) and [Encore's runtime library](https://github.com/encoredev/encore/tree/main/runtime) to see the results of this talk being used in practice.

## Talk Videos

- [A short version given at a London Gophers meetup](https://www.youtube.com/watch?v=CymVdee2Q8Y)
- [The full version given at GopherCon UK 2023](https://www.youtube.com/watch?v=MRZU5J29Rys)

## Example app

The example app is a simple todo app. It contains three services:

- `todo-svc` - The todo service is responsible for managing the todos for users
- `user-svc` - The user service is responsible for authenticating users. In this example it's hardcoded to only allow
  one user with bearer token `secret`.
- `api-svc` - The API service acts the the gateway for the user's requests. It is responsible for authenticating the user
  against the `user-svc` and then forwarding the request to the `todo-svc`.

### Branches

There are several branches with various different stages of tracing enabled:
- `before-tracing` contains all a clean version of the example application with no modifications to Go
- `initial-tracing-code` tracks trace context against Go routines, adding hooks into the standard library to track HTTP servers handling requests, and HTTP clients making calls.
- `main` contains a final version of the code, in which we pass a Trace Context between services to maintain context, track database calls being made and emit traces to Jaeger
- `with-goroutine-tracing` adds spans for every Go routine which is spawned during the trace.

### Running the example app

To run the apps you will need to brew install `make`, `postgres` and `overmind`.

```bash
make initdb        # Create the database
make postgres &    # Start the database
make jaeger &      # Start Jaeger via a Docker image
make microservices # Start the microservices (you only need to run this when changing branches)
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

## Differences between the talk and this repo

In the talk, I said we'd embed Go as a submodule, however to keep this code easier to switch between states, this repo actually uses subtrees which allows us to track changes in the Go runtime per branch without having to reapply patches or push changes into an upstream submodule.

For a example of how I talked about managing the fork, check out the [Encore Go Rolling Fork](https://github.com/encoredev/go).

The initial subtree was added to this repo using this command:

```bash
git subtree pull --prefix go-src https://go.googlesource.com/go release-branch.go1.20 --squash
```
