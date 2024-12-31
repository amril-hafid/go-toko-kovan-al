package web

import (
	"context"
	"errors"
	"fmt"
	"go-toko-kovan-al/app/image_product"
	"go-toko-kovan-al/app/product"
	"go-toko-kovan-al/app/users"
	"go-toko-kovan-al/config"
	"go-toko-kovan-al/helper"
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type adminWebImageProductHandler struct {
	imageProductService image_product.Service
	productService      product.Service
	userService         users.Service
	sessionStore        *session.Store
	config              config.Config
}

func NewAdminWebImageProductHandler(imageProductService image_product.Service, productService product.Service, userService users.Service, sessionStore *session.Store, config config.Config) *adminWebImageProductHandler {
	return &adminWebImageProductHandler{imageProductService, productService, userService, sessionStore, config}
}

func (h *adminWebImageProductHandler) NewImageProductView(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 60*time.Second)
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

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	product, err := h.productService.GetProductByID(c, uint(id))
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	return ctx.Render("dasboard/admin/product/dasboard_product_image_new", fiber.Map{
		"header": userMain,
		"role":   userMain.Role,
		"data":   product,
		"layout": "form",
		"alert":  template.HTML(alert),
	})
}

func (h *adminWebImageProductHandler) NewImageProduct(ctx *fiber.Ctx) error {

	c, cencel := context.WithTimeout(ctx.Context(), 60*time.Second)
	defer cencel()
	root := ctx.Get("Referer")

	session, err := h.sessionStore.Get(ctx)
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	inputNew := new(image_product.InputImageProduct)
	if err := ctx.BodyParser(inputNew); err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	product, err := h.productService.GetProductByID(c, inputNew.IDProduct)
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if product.Product.ID != inputNew.IDProduct {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	file, err := ctx.FormFile("file-image")
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if fileCheckSize := helper.ChekFileSize(file, h.config); fileCheckSize == false {
		helper.AlertMassage("msg-alert-new-product", "Ukuran File Melebih yang sistem tentukan maksimal 15MB", "error", session)
		return ctx.Redirect(root)
	}

	if file == nil {
		helper.AlertMassage("msg-alert-new-product", errors.New("File gagal di upload!").Error(), "error", session)
		return ctx.Redirect(root)
	}

	a, _ := helper.AllowedFileType(file, h.config)
	if a != true {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	fileName := helper.RenameFile(file)

	pathFile := fmt.Sprintf("assets/%s%s", h.config.FileConf.FileDerektory, fileName)
	err = ctx.SaveFile(file, pathFile)
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}
	path := fmt.Sprintf("%s%s", h.config.FileConf.FileDerektory, fileName)

	inputNew.NameFile = fileName
	inputNew.URL = path

	_, err = h.imageProductService.CreateImageProduc(c, *inputNew)
	if err != nil {
		errFile := helper.DeleteFile(path)
		if errFile != nil {
			helper.AlertMassage("msg-alert-new-product", errFile.Error(), "error", session)
			return ctx.Redirect(root)
		}
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-product", "Image baru berhasil di unggah.", "success", session)
	pathRoot := fmt.Sprintf("/dasboard/product/detail/%d", product.Product.ID)
	return ctx.Redirect(pathRoot)
}

func (h *adminWebImageProductHandler) DeleteImageProduct(ctx *fiber.Ctx) error {
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

	_, err = h.imageProductService.DeleteImageProduct(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-product", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-product", "Product berhasil di hapus!", "success", session)
	return ctx.Redirect(root)
}
