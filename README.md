# rho

Generic HTTP request handler abstraction using [github.com/go-chi/chi/v5](https://github.com/go-chi/chi)
and [github.com/go-chi/render](https://github.com/go-chi/render).

## Usage

```go
package main

import (
    "context"
    "io"
    "net"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/bdreece/rho"
)

func main() {
    mux := chi.NewMux()

    mux.Post("/echo", rho.HandleFunc(echo))

    // ...
}

func echo(ctx context.Context, req *echoRequest) (*echoResponse, error) {
    return &echoResponse(*echoRequest), nil
}

type echoRequest string

func (req *echoRequest) Bind(r *http.Request) error {
    defer r.Body.Close()

    content, err := io.ReadAll(r.Body)
    if err != nil {
        return err
    }

    *req = echoRequest(content)
    return nil
}

type echoResponse string

func (res *echoResponse) Render(w http.ResponseWriter, _ *http.Request) error {
    if _, err := io.WriteString(w, string(*res)); err != nil {
        return err
    }

    return nil
}
```
