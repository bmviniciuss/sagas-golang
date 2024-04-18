package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

func RenderError(w http.ResponseWriter, r *http.Request, errorResponse Error) {
	w.WriteHeader(errorResponse.Status)
	render.JSON(w, r, errorResponse)
}
