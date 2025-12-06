package rho

import (
	"net/http"

	"github.com/go-chi/render"
)

func init() {
	render.Decode = decode
	render.Respond = respond
}

func decode(r *http.Request, v any) error {
	if v, ok := v.(render.Binder); ok {
		return v.Bind(r)
	}

	return render.DefaultDecoder(r, v)
}

func respond(w http.ResponseWriter, r *http.Request, v any) {
	if v, ok := v.(render.Renderer); ok {
		if err := v.Render(w, r); err != nil {
			panic(err)
		}

		return
	}

	if err, ok := v.(error); ok {
		HandleError(w, r, err)
		return
	}

	render.DefaultResponder(w, r, v)
}
