package transport

import (
	"net/http"

	productDecoders "github.com/Thomika1/KarHub/pkg/domains/business/products/transport"
	"github.com/Thomika1/KarHub/pkg/domains/business/restock"
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(svc restock.ServiceI) http.Handler {
	r := chi.NewRouter()

	priorities := kithttp.NewServer(
		restock.Priorities(svc),
		productDecoders.DecodeListRequest,
		kithttp.EncodeJSONResponse,
		kithttp.ServerErrorEncoder(productDecoders.ErrorEncoder),
	)

	r.Post("/priorities", priorities.ServeHTTP)

	return r
}
