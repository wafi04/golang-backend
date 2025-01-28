package internal

type CustomStock struct {
	Quantity  int64  `json:"quantity"`
	VariantId string `json:"variant_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
