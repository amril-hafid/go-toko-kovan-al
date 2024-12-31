package category

type InputCategory struct {
	Name   string `json:"name" form:"name" validate:"required"`
	Status string `json:"status" form:"status"`
}

type UpdateCategory struct {
	ID     uint   `json:"id" form:"id" validate:"required"`
	Name   string `json:"name" form:"name" validate:"required"`
	Status string `json:"status" form:"status"`
}
