package onlineshopping

type Cart struct {
	Cart map[int]OrderItems
}

func NewCart() *Cart{
	c := &Cart{
		Cart: make(map[int]OrderItems)
	}
	return c
}

func (c *Cart) AddInCart(product *Product, quantity int) {

}

func (c *Cart) RemoveFromCart(product *Product, quantity int) {

}

func (c *Cart) ShowCartItems() {

}

func (c *Cart) ClearCart(){

}