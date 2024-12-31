package web

import (
	"context"
	"fmt"
	"go-toko-kovan-al/app/category"
	"go-toko-kovan-al/app/product"
	"go-toko-kovan-al/app/users"
	"go-toko-kovan-al/config"
	"go-toko-kovan-al/helper"
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type adminWebProductHandler struct {
	productService  product.Service
	categoryService category.Service
	userService     users.Service
	sessionStore    *session.Store
	config          config.Config
}

func NewAdminWebProductHandler(productService product.Service, categoryService category.Service, userService users.Service, sessionStore *session.Store, config config.Config) *adminWebProductHandler {
	return &adminWebProductHandler{productService, categoryService, userService, sessionStore, config}
}

func (h *adminWebProductHandler) ShowAllProduct(ctx *fiber.Ctx) error {

	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
	} else {

		sessionResult = session.Get("msg-alert-new-product")
		staussession = session.Get("msg-alert-new-product-status")
	}

	alert := helper.AlertString(sessionResult, staussession)

	cookisUserId := ctx.Cookies("sessionLog")
	idUser, err := helper.GetSessionID(cookisUserId, h.config)
	if err != nil {
		return ctx.Redirect("/login")
	}

	userMain, err := h.userService.GetUserByID(c, idUser)
	if err != nil {
		return ctx.Redirect("/login")
	}

	product, _ := h.productService.GetAllProduct(c)

	category, _ := h.categoryService.GetAllCategory(c)

	return ctx.Render("dasboard/admin/product/dasboard_product_list", fiber.Map{
		"header":   userMain,
		"role":     userMain.Role,
		"data":     product,
		"category": category,
		"layout":   "table",
		"alert":    template.HTML(alert),
	})
}
func (h *adminWebProductHandler) ShowDetailProduct(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	session, err := h.sessionStore.Get(ctx)

	cookisUserId := ctx.Cookies("sessionLog")
	idUser, err := helper.GetSessionID(cookisUserId, h.config)
	if err != nil {
		return ctx.Redirect("/login")
	}

	userMain, err := h.userService.GetUserByID(c, idUser)
	if err != nil {
		return ctx.Redirect("/login")
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect("/dasboard/product")
	}

	product, err := h.productService.GetProductByID(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect("/dasboard/product")
	}

	return ctx.Render("dasboard/admin/product/dasboard_product_detail", fiber.Map{
		"header": userMain,
		"role":   userMain.Role,
		"data":   product,
		"layout": "table",
	})
}
func (h *adminWebProductHandler) NewProduct(ctx *fiber.Ctx) error {

	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()
	root := ctx.Get("Referer")

	session, err := h.sessionStore.Get(ctx)
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	cookisUserId := ctx.Cookies("sessionLog")
	userID, err := helper.GetSessionID(cookisUserId, h.config)
	if err != nil {
		return ctx.Redirect("/login")
	}

	inputNew := new(product.InputProduct)
	if err := ctx.BodyParser(inputNew); err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	fmt.Println("ini data input handler new product 1 :", inputNew)

	productRow, err := h.productService.CreateProduct(c, *inputNew, userID)
	fmt.Println("ini data input handler new product 2:", productRow)

	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-product", "Product baru berhasil di unggah.", "success", session)
	pathRoot := fmt.Sprintf("/dasboard/product/image/new/%d", productRow.ID)
	return ctx.Redirect(pathRoot)
}
func (h *adminWebProductHandler) UpdateProductView(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()

	} else {
		sessionResult = session.Get("msg-alert-update-product")
		staussession = session.Get("msg-alert-update-product-status")
	}

	alert := helper.AlertString(sessionResult, staussession)

	cookisUserId := ctx.Cookies("sessionLog")
	idUser, err := helper.GetSessionID(cookisUserId, h.config)
	if err != nil {
		return ctx.Redirect("/login")
	}

	userMain, err := h.userService.GetUserByID(c, idUser)
	if err != nil {
		return ctx.Redirect("/login")
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect("/dasboard/product")
	}
	productRow, err := h.productService.GetProductByID(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect("/dasboard/product")
	}

	category, _ := h.categoryService.GetAllCategory(c)

	return ctx.Render("dasboard/admin/product/dasboard_product_update", fiber.Map{
		"header":   userMain,
		"role":     userMain.Role,
		"data":     productRow,
		"category": category,
		"layout":   "form",
		"alert":    template.HTML(alert),
	})
}
func (h *adminWebProductHandler) UpdateProduct(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()
	root := ctx.Get("Referer")
	input := new(product.UpdateProduct)
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-update-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if err = ctx.BodyParser(input); err != nil {
		helper.AlertMassage("msg-alert-update-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if _, err = h.productService.UpdateProduct(c, *input); err != nil {
		helper.AlertMassage("msg-alert-update-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-update-product", "Update product berhasil.", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebProductHandler) DeleteProductSoft(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.productService.DeleteProduct(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-product", "Product berhasil di hapus!", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebProductHandler) DeleteProductRecycle(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.productService.DeleteProduct(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-product", "Product berhasil di hapus!", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebProductHandler) ShowAllProductRecycle(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
	} else {

		sessionResult = session.Get("msg-alert-new-product")
		staussession = session.Get("msg-alert-new-product-status")
	}

	alert := helper.AlertString(sessionResult, staussession)

	cookisUserId := ctx.Cookies("sessionLog")
	idUser, err := helper.GetSessionID(cookisUserId, h.config)
	if err != nil {
		return ctx.Redirect("/login")
	}

	userMain, err := h.userService.GetUserByID(c, idUser)
	if err != nil {
		return ctx.Redirect("/login")
	}
	products, err := h.productService.GetAllProductDeleted(c)
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	return ctx.Render("dasboard/admin/product/dasboard_product_list_recycle", fiber.Map{
		"header": userMain,
		"role":   userMain.Role,
		"data":   products,
		"layout": "table",
		"alert":  template.HTML(alert),
	})
}
func (h *adminWebProductHandler) RestoreProduct(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.productService.RestoreProduct(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-product", "DAta berhasil di pilihkan!", "success", session)
	return ctx.Redirect(root)
}
