package main

import (
	handlers "fc-desafio-multi/handler"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/consulta-cep", func(r chi.Router) {
		r.Get("/{cep}", handlers.GetCepValue)
	})

	http.ListenAndServe("127.0.0.1:8000", r)
}
