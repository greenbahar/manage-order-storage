package order

type Order struct {
	OrderId int    `json:"order_id"`
	Price   int    `json:"price"`
	Title   string `json:"title"`
}
