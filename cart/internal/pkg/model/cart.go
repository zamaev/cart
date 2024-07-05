package model

import "sync"

type Cart map[ProductSku]uint16

type CartFull map[Product]uint16

type CartFullMx struct {
	mx   sync.Mutex
	cart CartFull
}

func NewCartFullMx(length int) CartFullMx {
	return CartFullMx{
		cart: make(CartFull, length),
	}
}

func (cf *CartFullMx) Add(product Product, count uint16) {
	cf.mx.Lock()
	defer cf.mx.Unlock()
	cf.cart[product] += count
}

func (cf *CartFullMx) GetCartFull() CartFull {
	return cf.cart
}
