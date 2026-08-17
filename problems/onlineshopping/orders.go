package onlineshopping

type Order struct {

	OrderId int
	User *User
	Items []OrderItems
	Amount int
}

func NewOrder() *Order {

} 

func (o *Order) calculateAmount() int {

}

func (o *Order) updateStatus(){

}

