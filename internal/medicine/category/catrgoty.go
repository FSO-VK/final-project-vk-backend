package category

type Category struct {
	ID   uint
	Name string
}

func NewCategory(name string) *Category {
	return &Category{
		Name: name,
	}
}
