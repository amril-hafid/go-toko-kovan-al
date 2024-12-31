package image_product

import (
	"time"
)

type ImageFormatter struct {
	ID             uint
	Index          int
	Name           string
	IDProduct      uint
	URL            string
	IsPrimary      int
	CreatedAt      time.Time
	ImageNoPrimary []ImageAllNoPrimary
}

type ImageAllNoPrimary struct {
	ID        uint
	Index     int
	Name      string
	IDProduct uint
	URL       string
	IsPrimary int
	CreatedAt time.Time
}

func ImageProductFormatter(imagePrimary ImageProduct, imageNoPrimary []ImageProduct) ImageFormatter {

	var imageFormatter ImageFormatter

	imageFormatter.ID = imagePrimary.ID
	imageFormatter.Index = 1
	imageFormatter.Name = imagePrimary.Name
	imageFormatter.IDProduct = imagePrimary.IDProduct
	imageFormatter.IsPrimary = imagePrimary.IsPrimary
	imageFormatter.URL = imagePrimary.URL
	imageFormatter.CreatedAt = imagePrimary.CreatedAt

	images := []ImageAllNoPrimary{}
	for i, image := range imageNoPrimary {
		var imageAllNoPrimary ImageAllNoPrimary

		imageAllNoPrimary.ID = image.ID
		imageAllNoPrimary.Index = i + 2
		imageAllNoPrimary.Name = image.Name
		imageAllNoPrimary.URL = image.URL
		imageAllNoPrimary.IsPrimary = image.IsPrimary
		imageAllNoPrimary.IDProduct = image.IDProduct
		imageAllNoPrimary.CreatedAt = image.CreatedAt

		images = append(images, imageAllNoPrimary)
	}

	imageFormatter.ImageNoPrimary = images

	return imageFormatter
}
