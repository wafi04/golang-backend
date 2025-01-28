package producthandler


type ProductRequest struct {
	Id            string     `json:"id"`
	Name          string     `json:"name"`
	SubTitle      *string    `json:"sub_title"`
	Description   string     `json:"description"`
	Sku           string     `json:"sku"` 
	Price         float64    `json:"price"`
	CategoryId    string     `json:"category_id"`
}
