package web

import (
	"context"
	"fmt"
	"go-toko-kovan-al/app/banner"
	"go-toko-kovan-al/app/product"
	"go-toko-kovan-al/config"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type productHandler struct {
	service       product.Service
	bannerService banner.Service
	sessionStore  *session.Store
	config        config.Config
}

func NewProductHandler(service product.Service, bannerService banner.Service, sessionStore *session.Store, config config.Config) *productHandler {
	return &productHandler{service, bannerService, sessionStore, config}
}

func (h *productHandler) Home(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()
	banner, _ := h.bannerService.GetBannerByLimit(c, int(4))

	product, _ := h.service.GetAllWebProductHOme(c, int(8))

	return ctx.Render("web/web/home", fiber.Map{
		"page":   "home",
		"data":   product,
		"banner": banner,
	})
}

func (h *productHandler) ShopDetail(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 9*time.Second)
	defer cencel()

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Redirect("/shop")
	}

	product, err := h.service.GetProductWebByID(c, uint(id))
	if err != nil {
		return ctx.Redirect("/shop")
	}

	relatedProduct, _ := h.service.GetAllWebProductHOme(c, 4)
	fmt.Println(relatedProduct)
	return ctx.Render("web/web/shop_detail", fiber.Map{
		"page":    "shop-detail",
		"data":    product,
		"related": relatedProduct,
	})
}

func (h *productHandler) Shop(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "12"))
	search := ctx.Query("search", "")
	category := ctx.Query("category", "")

	setIn := product.Paging{
		Page:     page,
		Limit:    limit,
		Search:   search,
		Category: category,
	}

	product, _ := h.service.GetAllWebProduct(c, &setIn)
	// if err != nil {
	// 	return ctx.Redirect("/home")
	// }

	fmt.Println("data handler satu : ", product)
	return ctx.Render("web/web/shop", fiber.Map{
		"page": "shop",
		"data": product,
	})
}

func (h *productHandler) ProductDetail(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Redirect("/shop")
	}

	productDetail, err := h.service.GetbyIDWebProduct(c, uint(id))
	if err != nil {
		return ctx.Redirect("/shop")
	}

	return ctx.Render("web/web/shop_detail", fiber.Map{
		"page": "shop_detail",
		"data": productDetail,
	})
}

func (h *productHandler) Blog(ctx *fiber.Ctx) error {
	return ctx.Render("web/web/blog", fiber.Map{
		"page": "blog",
	})
}

func (h *productHandler) BlogDetail(ctx *fiber.Ctx) error {
	return ctx.Render("web/web/blog_detail", fiber.Map{
		"page": "blog_detail",
	})
}

func (h *productHandler) Contact(ctx *fiber.Ctx) error {
	return ctx.Render("web/web/contact", fiber.Map{
		"page": "contact",
	})
}
