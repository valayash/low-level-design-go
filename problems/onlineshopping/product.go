package onlineshopping

type Product struct {
	ProductId int
	ProductName string
	Quantity int
	Price int
}

func NewProduct(id int, name string, quantity int, price int) *Product{
	p := &Product{
		ProductId : id,
		ProductName : name,
		Quantity : quantity,
		Price : price
	}

	return p
}

func (p *Product) UpdateQuantity(quantity int){
	p.Quantity += quantity
}