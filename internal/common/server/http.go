package server

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func RunHTTPServer(createHandler func(router chi.Router) http.Handler) {
	apiRouter := chi.NewRouter()
	rootRouter := chi.NewRouter()
	rootRouter.Mount("/api", createHandler(apiRouter))
	_ = http.ListenAndServe(":"+os.Getenv("PORT"), rootRouter)
}
