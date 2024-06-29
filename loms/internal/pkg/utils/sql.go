package utils

import (
	"strings"
)

func IsSqlWrite(sql string) bool {
	return strings.Contains(sql, "INSERT") || strings.Contains(sql, "UPDATE") || strings.Contains(sql, "DELETE")
}
