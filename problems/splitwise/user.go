package splitwise

type User struct {
	Id int
	Name string
	Email string
	Balances map[int]float64
}

func NewUser(id int, name, email string){

	return &User{
		Id : id,
		Name : name,
		Email : email
		Balances : make(map[int]float64)
	}
}
