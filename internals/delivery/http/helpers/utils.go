package helpers

import (
	"context"
	"encoding/json"
	"gostarter/internals/domain"
	"gostarter/pkg/utils"
	"io"
	"net/http"
	"strconv"
)

func ParseRequest[T any](data io.ReadCloser) (T, error) {
	var obj T
	err := json.NewDecoder(data).Decode(&obj)
	return obj, err
}

func WriteResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func GetAccountFromContext(ctx context.Context) (*domain.Account, error) {
	acc, ok := ctx.Value("account").(*domain.Account)
	if !ok {
		return nil, domain.ErrGettingAccountInfo
	}
	return acc, nil
}

func GetPaginationParams(r *http.Request) utils.PaginationParams {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageParams := utils.PaginationParams{
		Page: 1,
		Size: 10,
	}

	if page != "" {
		pageParams.Page, _ = strconv.Atoi(page)
	}
	if limit != "" {
		pageParams.Size, _ = strconv.Atoi(limit)
	}

	return pageParams
}
