[![GitHub Actions Workflow Status](https://github.com/bdreece/chi-action/actions/workflows/build.yml/badge.svg)](https://github.com/bdreece/chi-action/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/bdreece/chi-action.svg)](https://pkg.go.dev/github.com/bdreece/chi-action)

# chi-action

Generic HTTP request handler abstraction using [github.com/go-chi/chi/v5](https://pkg.go.dev/github.com/go-chi/chi/v5)
and [github.com/go-chi/render](https://pkg.go.dev/github.com/go-chi/render).

## Usage

```go
package main

import (
    "context"
    "io"
    "net"
    "net/http"

    "github.com/bdreece/chi-action"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
)

func main() {
    mux := chi.NewMux()

    mux.Post("/echo", action.HandleFunc(echo))

    // ...
}

func echo(ctx context.Context, req *echoRequest) (*echoResponse, error) {
    return &echoResponse(*echoRequest), nil
}

type (
    echoRequest string
    echoResponse string
)

func (req *echoRequest) Bind(r *http.Request) error {
    return render.DecodeJSON(r.Body, req)
}


func (res *echoResponse) Render(w http.ResponseWriter, r *http.Request) error {
    render.JSON(w, r, res)
    return nil
}
```
