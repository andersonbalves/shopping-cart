package model

type CartItem struct {
	UserId      string `json:"userId"`
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	Quantity    int    `json:"quantity"`
}
