package splitwise

type Group struct {
	Id int
	Name string
	Members []*User
	Expenses []*Expense
}

func NewGroup(id int, name string) *Group {

	u := &User{
		Id : id,
		Name : name,
		Members : []*User{},
		Expenses : []*Expense{}
	}

	return u
}

func (g *Group) AddMember(user User) {

	g.Members = append(g.Members, user)
}

func (g *Group) AddExpense(expense Expense) {

	g.Expenses = append(g.Expenses, expense)
}
