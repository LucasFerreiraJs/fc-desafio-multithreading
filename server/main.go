package main

import (
	handlers "fc-desafio-multi/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Route("/consultaCep", func(chi.Router) {
		r.Get("/{cep}", handlers.GetCepValue)
	})

	http.ListenAndServe("0.0.0.0:8000", r)
}
