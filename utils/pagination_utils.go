package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

// GetPaginationLimit и GetPaginationOffset возвращают параметры пагинации
func GetPaginationLimit(r *http.Request) (int, error) {
    limitStr := r.URL.Query().Get("limit")
    if limitStr == "" {
        return 5, nil // дефолтное значение
    }
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        return 0, err
    }
    return limit, nil
}

func GetPaginationOffset(r *http.Request) (int, error) {
    offsetStr := r.URL.Query().Get("offset")
    if offsetStr == "" {
        return 0, nil // дефолтное значение
    }
    offset, err := strconv.Atoi(offsetStr)
    if err != nil {
        return 0, err
    }
    return offset, nil
}

func GetPagination(r *http.Request) (int, int, int, error) {
	var limit, offset int
	var err error
	limit, err = GetPaginationLimit(r)
	if err != nil {
		return limit, offset, http.StatusBadRequest, fmt.Errorf("неверный лимит возвращаемых объектов (limit)")
	}

	offset, err = GetPaginationOffset(r)
	if err != nil {
		return limit, offset, http.StatusBadRequest, fmt.Errorf("неверное смещение пагинации (offset)")
	}
	return limit, offset, http.StatusOK, nil
}