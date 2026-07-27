package transport

import (
	"net/http"

	"github.com/Thomika1/KarHub/pkg/domains/business/products"
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(svc products.ServiceI) http.Handler {
	r := chi.NewRouter()

	create := kithttp.NewServer(
		products.Create(svc),
		DecodeCreateRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(ErrorEncoder),
	)

	get := kithttp.NewServer(
		products.Get(svc),
		DecodeIDRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(ErrorEncoder),
	)

	list := kithttp.NewServer(
		products.List(svc),
		DecodeListRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(ErrorEncoder),
	)

	update := kithttp.NewServer(
		products.Update(svc),
		DecodeUpdateRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(ErrorEncoder),
	)

	delete := kithttp.NewServer(
		products.Delete(svc),
		DecodeIDRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(ErrorEncoder),
	)

	r.Post("/", create.ServeHTTP)
	r.Post("/list", list.ServeHTTP)
	r.Get("/{id}", get.ServeHTTP)
	r.Put("/{id}", update.ServeHTTP)
	r.Delete("/{id}", delete.ServeHTTP)

	return r
}
