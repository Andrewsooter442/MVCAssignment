package config

type Client struct {
	Name    string
	IsAdmin bool
	IsChef  bool
	TableNo int
}

type HomePageData struct {
	Client        Client
	Menu          Menu
	PendingOrders []Order
	StatusMessage string
}
