package banner

type InputBanner struct {
	Name        string ` form:"name"`
	NameImage   string
	ImageURL    string
	URL         string ` form:"url"`
	Description string ` form:"description"`
	IsPrimary   int    ` form:"is_promary"`
	Status      string ` form:"status"`
}

type UpdateBanner struct {
	ID          uint   ` form:"id"`
	Name        string ` form:"name"`
	URL         string `form:"url"`
	Description string ` form:"description"`
	IsPrimary   int    ` form:"is_promary"`
	Status      string ` form:"status"`
}
