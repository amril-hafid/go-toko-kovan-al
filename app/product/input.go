package product

type Paging struct {
	Page  int
	Limit int
	// Offset int
	Category string
	Total    int64
	Search   string
}

type InputProduct struct {
	SKU              string `form:"sku" validate:"required"`
	Name             string `form:"name" validate:"required"`
	Price            int    `form:"price" validate:"required"`
	Stock            int    `form:"stock" validate:"required"`
	ShortDescription string `form:"short_description"`
	LongDescription  string `form:"long_description"`
	SizeType         string `form:"size_type"`
	Long             int    `form:"long"`
	Wide             int    `form:"wide"`
	Tall             int    `form:"tall"`
	Diameter         int    `form:"diameter"`
	IDKategori       int    `form:"id_kategory" validate:"required"`
	Status           int    `form:"status" validate:"required"`
}

type UpdateProduct struct {
	ID               uint   `form:"id" validate:"required"`
	SKU              string `form:"sku" validate:"required"`
	Name             string `form:"name" validate:"required"`
	Price            int    `form:"price" validate:"required"`
	Stock            int    `form:"stock" validate:"required"`
	ShortDescription string `form:"short_description"`
	LongDescription  string `form:"long_description"`
	SizeType         string `form:"size_type"`
	Long             int    `form:"long"`
	Wide             int    `form:"wide"`
	Tall             int    `form:"tall"`
	Diameter         int    `form:"diameter"`
	IDKategori       uint   `form:"id_kategory" validate:"required"`
	Status           int    `form:"status" validate:"required"`
}
