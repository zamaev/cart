package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetIntPahtValue(r *http.Request, field string) (int64, error) {
	valueRaw := r.PathValue(field)
	value, err := strconv.ParseInt(valueRaw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt %s: %w", field, err)
	}
	return value, nil
}
