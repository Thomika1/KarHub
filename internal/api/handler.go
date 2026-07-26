package api

import (
	"context"
	"net/http"

	"github.com/Thomika1/KarHub/internal/container"
	productTransport "github.com/Thomika1/KarHub/pkg/domains/business/products/transport"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

var r *chi.Mux

func Handler(ctx context.Context, dep *container.Dependency) http.Handler {

	r = chi.NewMux()

	r.Use(cors.Handler(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	// Routes

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/product", productTransport.NewHTTPHandler(dep.Services.ProductService))
	})

	return r
}
