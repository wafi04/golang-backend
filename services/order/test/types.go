package testing

type CreateOrderRequest struct {
	VariantsID string  `json:"variants_id"`
	Quantity   int64   `json:"quantity"`
	Total      float64 `json:"total"`
	UserID     string  `json:"user_id"`
}
