package request

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

const (
	ParamNotFound = -1
)

// Read path parameter `id`
func IDOf(r *http.Request) (int64, error) {
	ids := chi.URLParam(r, "id")

	if ids == "" {
		return -1, errors.Errorf("id was not presence")
	}

	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		return -1, errors.Errorf("id was not number")
	}

	return id, nil
}

func DecodeParam2Int64(r *http.Request, key string) (int64, error) {
	q := r.URL.Query()
	str := q.Get(key)
	if str == "" {
		return -1, nil
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return -1, errors.Errorf("query value (?%s=%v) was not number", key, str)
	}
	return n, nil
}
