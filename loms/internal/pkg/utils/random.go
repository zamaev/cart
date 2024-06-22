package utils

import "math/rand"

func GetRandomBool() bool {
	return rand.Intn(2) == 1
}
