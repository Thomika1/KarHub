package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Thomika1/KarHub/pkg/core/crud"
	"github.com/Thomika1/KarHub/pkg/domains/business/products"
	domainErrors "github.com/Thomika1/KarHub/pkg/domains/shared/errors"
	"github.com/Thomika1/KarHub/pkg/domains/shared/models/business"
	"github.com/go-chi/chi"
)

type createRequest struct {
	business.ProductData
}

type listRequest struct {
	Query crud.Query
}

func DecodeCreateRequest(_ context.Context, r *http.Request) (any, error) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req.ProductData.ToModel(), nil
}

func DecodeListRequest(_ context.Context, r *http.Request) (any, error) {
	if r.ContentLength == 0 {
		return crud.Query{}, nil
	}
	var req listRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req.Query, nil
}

func DecodeIDRequest(_ context.Context, r *http.Request) (any, error) {
	id := chi.URLParam(r, "id")
	return id, nil
}

func DecodeUpdateRequest(_ context.Context, r *http.Request) (any, error) {
	id := chi.URLParam(r, "id")
	var data business.ProductData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return products.UpdateRequest{ID: id, Data: data}, nil
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	var conflict domainErrors.ConflictError
	var validation domainErrors.ValidationError

	switch {
	case errors.As(err, &conflict):
		w.WriteHeader(http.StatusConflict)
	case errors.As(err, &validation):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
