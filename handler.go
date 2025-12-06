package rho

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
)

type Handler[Req, Res any] interface {
	Handle(ctx context.Context, req *Req) (*Res, error)
}

type HandlerFunc[Req, Res any] func(ctx context.Context, req *Req) (*Res, error)

func (handle HandlerFunc[Req, Res]) Handle(ctx context.Context, req *Req) (*Res, error) {
	return handle(ctx, req)
}

func Handle[Req, Res any](handler Handler[Req, Res]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(Req)
		if err := render.Decode(r, req); err != nil {
			render.Respond(w, r, err)
			return
		}

		if err := Validate(r.Context(), req); err != nil {
			render.Respond(w, r, err)
			return
		}

		res, err := handler.Handle(r.Context(), req)
		if err != nil {
			render.Respond(w, r, err)
			return
		}

		render.Respond(w, r, res)
	}
}

func HandleFunc[Req, Res any](handler HandlerFunc[Req, Res]) http.HandlerFunc {
	return Handle(handler)
}
