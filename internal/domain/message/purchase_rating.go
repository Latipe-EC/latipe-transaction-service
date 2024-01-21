package message

type RatingMessage struct {
	OrderItemId string `json:"orderItemId"`
	ProductId   string `json:"productId"`
	RatingId    string `json:"ratingId"`
	Rating      int    `json:"rating"`
	Op          string `json:"op"`
}
