package object

type Author struct {
	Name    string
	Surname string
	Email   string
}

func (a *Author) Type() string {
	return "author"
}
