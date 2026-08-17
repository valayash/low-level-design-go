package splitwise

type Expense struct {
	Id int
	Amount float64
	PaidBy *User
	Splits 
}