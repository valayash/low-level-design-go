package onlineshopping

type User struct {
	Id int
	Name string
	Email string
	Orders []*Order
}

func NewUser(id int, name , email string) {

	u := &User{
		Id : id,
		Name : name,
		Email : email,
		Orders : []*Order{}
	}

	return u
}