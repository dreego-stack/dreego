package catalog

type Product struct {
	Name    string
	Price   string
	InStock bool
}

func Products() []Product {
	return []Product{
		{Name: "Dreego Mug", Price: "$12", InStock: true},
		{Name: "Dreego Tee", Price: "$24", InStock: false},
	}
}
