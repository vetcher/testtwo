package resolvers

type OnlyID struct {
	ID uint `validate:"required,min=1"`
}

type Comment struct {
	Text  string `validate:"required,omitempty"`
	Login string `validate:"required"`
}

type User struct {
	Login    string `validate:"required"`
	Password string `validate:"required"`
	Banned   bool
}
