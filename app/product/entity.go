package product

import (
	"go-toko-kovan-al/app/image_product"
	"html/template"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Index            int
	ID               uint `gorm:"primaryKey"`
	SKU              string
	Name             string
	Price            int
	Stock            int
	ShortDescription string
	LongDescription  string
	SizeType         string
	Long             int
	Wide             int
	Tall             int
	Diameter         int
	IDKategori       uint
	Status           int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type ProductView struct {
	Index            int
	ID               uint
	SKU              string
	Name             string
	Price            int
	Stock            int
	ShortDescription string
	LongDescription  string
	SizeType         string
	Long             int
	Wide             int
	Tall             int
	Diameter         int
	IDKategori       uint
	CategoryName     string
	Status           int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}

type ProductListWebViews struct {
	Index    int
	Product  []ProductWebList
	Total    int
	Page     int
	LastPage int
	PrevPage int
	NextPage int
	Search   string
	Category string
}

type ProductEndImageDetail struct {
	Product Product
	Image   image_product.ImageFormatter
	Search  string
}

type ProductDetail struct {
	Product ProductWebDetail
	Image   []image_product.ImageProduct
}

type ProductWebList struct {
	Index        int
	ID           uint
	SKU          string
	Name         string
	Price        string
	Stock        int
	ImagePrimary template.HTML
	IDKategori   uint
	Status       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type ProductWebDetail struct {
	Index            int
	ID               uint `gorm:"primaryKey"`
	SKU              string
	Name             string
	Price            string
	Stock            int
	ShortDescription string
	LongDescription  string
	SizeType         string
	Long             int
	Wide             int
	Tall             int
	Diameter         int
	IDKategori       uint
	Status           int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

// {{range .data.Product}}
//                 <div class="col-6 col-lg-3 col-md-6 col-sm-6 mix hot-sales">
//                     <a href="">
//                     <div class="product__item sale">
//                         <div class="product__item__pic bag-cont" >
//                             <!-- <span class="label">Sale</span> -->
//                             <img src="{{.ImagePrimary}}"   alt="">

//                         </div>

//                         <div class="product__item__text">

//                             <h6>{{.Name}}</h6>
//                             <!-- <a href="#" class="add-cart"></a> -->
//                             <div class="rating">
//                                 <i class="fa fa-star"></i>
//                                 <i class="fa fa-star"></i>
//                                 <i class="fa fa-star"></i>
//                                 <i class="fa fa-star"></i>
//                                 <i class="fa fa-star-o"></i>
//                             </div>
//                             <h5>$98.49</h5>
//                             <div class="product__color__select">
//                                 <label for="pc-16">
//                                     <input type="radio" id="pc-16">
//                                 </label>
//                                 <label class="active black" for="pc-17">
//                                     <input type="radio" id="pc-17">
//                                 </label>
//                                 <label class="grey" for="pc-18">
//                                     <input type="radio" id="pc-18">
//                                 </label>
//                             </div>
//                         </div>
//                     </div>
//                 </a>
//                 </div>
//                 {{end}}

// <li><a href="#">Pages</a>
// <ul class="dropdown">
// 	<li><a href="./about.html">About Us</a></li>
// 	<!-- <li><a href="./shop-details.html">Shop Details</a></li>
// 	<li><a href="./shopping-cart.html">Shopping Cart</a></li>
// 	<li><a href="./checkout.html">Check Out</a></li> -->
// 	<li {{if eq .page "blog-detail"}} class="active" {{else}} {{end}}><a href="/blog-detail">Blog Details</a></li>
// </ul>
// </li>
// <li {{if eq .page "blog"}} class="active" {{else}} {{end}}><a href="/blog">Blog</a></li>
// <li {{if eq .page "contact"}} class="active" {{else}} {{end}}><a href="/contact">Contacts</a></li>
