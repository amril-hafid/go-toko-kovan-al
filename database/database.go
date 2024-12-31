package database

import (
	"fmt"
	"go-toko-kovan-al/app/banner"
	"go-toko-kovan-al/app/category"
	"go-toko-kovan-al/app/image_product"
	"go-toko-kovan-al/app/product"
	"go-toko-kovan-al/app/users"
	"go-toko-kovan-al/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(conf config.Config) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.DB.User, conf.DB.Pass, conf.DB.Host, conf.DB.Port, conf.DB.Name)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	DB = db

	db.AutoMigrate(&users.User{})
	db.AutoMigrate(&product.Product{})
	db.AutoMigrate(&category.Category{})
	db.AutoMigrate(&image_product.ImageProduct{})
	db.AutoMigrate(&banner.Banner{})

}
