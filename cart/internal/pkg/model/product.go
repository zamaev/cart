package model

type ProductSku int64

type Product struct {
	Sku   ProductSku
	Name  string
	Price uint32
}
