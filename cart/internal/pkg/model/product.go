package model

import (
	"encoding/json"
)

type ProductSku int64

type Product struct {
	Sku   ProductSku
	Name  string
	Price uint32
}

func (p Product) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Product) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
