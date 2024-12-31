package image_product

type InputImageProduct struct {
	IsPrimary int `json:"image_status" form:"image_status"`
	Status    string
	IDProduct uint `json:"id_product" form:"id_product" validate:"required"`
	NameFile  string
	URL       string
}
