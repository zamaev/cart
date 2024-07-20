package utils

import (
	"strings"
)

func IsSqlWrite(sql string) bool {
	return strings.Contains(sql, "INSERT") || strings.Contains(sql, "UPDATE") || strings.Contains(sql, "DELETE")
}

func GetSqlType(sql string) string {
	switch {
	case strings.Contains(sql, "SELECT"):
		return "SELECT"
	case strings.Contains(sql, "INSERT"):
		return "INSERT"
	case strings.Contains(sql, "UPDATE"):
		return "UPDATE"
	case strings.Contains(sql, "DELETE"):
		return "DELETE"
	}
	return "UNKNOWN"
}
