package config

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Item struct {
	ID          int     `json:"id"`
	CategoryID  int     `json:"categoryId"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}
type Menu struct {
	Categories []Category
	Items      []Item
}
