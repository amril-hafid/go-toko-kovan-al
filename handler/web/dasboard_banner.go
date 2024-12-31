package web

import (
	"context"
	"errors"
	"fmt"
	"go-toko-kovan-al/app/banner"
	"go-toko-kovan-al/app/users"
	"go-toko-kovan-al/config"
	"go-toko-kovan-al/helper"
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type adminWebBannerHandler struct {
	service      banner.Service
	userService  users.Service
	sessionStore *session.Store
	config       config.Config
}

func NewAdminWebBannerHandler(service banner.Service, userService users.Service, sessionStore *session.Store, config config.Config) *adminWebBannerHandler {
	return &adminWebBannerHandler{service, userService, sessionStore, config}
}

func (h *adminWebBannerHandler) ShowAllBanner(ctx *fiber.Ctx) error {

	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
	} else {

		sessionResult = session.Get("msg-alert-new-banner")
		staussession = session.Get("msg-alert-new-banner-status")
	}

	alert := helper.AlertString(sessionResult, staussession)

	cookisUserId := ctx.Cookies("sessionLog")
	idUser, err := helper.GetSessionID(cookisUserId, h.config)
	if err != nil {
		return ctx.Redirect("/login")
	}

	// var idUser uint = 14

	userMain, err := h.userService.GetUserByID(c, idUser)
	if err != nil {
		return ctx.Redirect("/login")
	}

	banner, _ := h.service.GetAllBannder(c)
	// if err != nil {
	// 	return ctx.Render("error400", fiber.Map{})
	// }

	return ctx.Render("dasboard/admin/banner/dasboard_banner_list", fiber.Map{
		"header": userMain,
		"data":   banner,
		"layout": "table",
		"alert":  template.HTML(alert),
	})
}

func (h *adminWebBannerHandler) NewBanner(ctx *fiber.Ctx) error {

	c, cencel := context.WithTimeout(ctx.Context(), 60*time.Second)
	defer cencel()
	root := ctx.Get("Referer")

	session, err := h.sessionStore.Get(ctx)
	if err != nil {
		helper.AlertMassage("msg-alert-new-banner", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	inputNew := new(banner.InputBanner)
	if err := ctx.BodyParser(inputNew); err != nil {
		helper.AlertMassage("msg-alert-new-banner", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	file, err := ctx.FormFile("file-image")
	if err != nil {
		helper.AlertMassage("msg-alert-new-banner", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if fileCheckSize := helper.ChekFileSize(file, h.config); fileCheckSize == false {
		helper.AlertMassage("msg-alert-new-banner", "Ukuran File Melebih yang sistem tentukan maksimal 15MB", "error", session)
		return ctx.Redirect(root)
	}

	if file == nil {
		helper.AlertMassage("msg-alert-new-banner", errors.New("File gagal di upload!").Error(), "error", session)
		return ctx.Redirect(root)
	}

	a, _ := helper.AllowedFileType(file, h.config)
	if a != true {
		helper.AlertMassage("msg-alert-new-banner", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	fileName := helper.RenameFile(file)

	pathFile := fmt.Sprintf("assets/%s%s", h.config.FileConf.FileDerektory, fileName)
	err = ctx.SaveFile(file, pathFile)
	if err != nil {
		helper.AlertMassage("msg-alert-new-banner", err.Error(), "error", session)
		return ctx.Redirect(root)
	}
	path := fmt.Sprintf("%s%s", h.config.FileConf.FileDerektory, fileName)

	inputNew.NameImage = fileName
	inputNew.ImageURL = path

	_, err = h.service.CreateBanner(c, *inputNew)
	if err != nil {
		errFile := helper.DeleteFile(path)
		if errFile != nil {
			helper.AlertMassage("msg-alert-new-banner", errFile.Error(), "error", session)
			return ctx.Redirect(root)
		}
		helper.AlertMassage("msg-alert-new-banner", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-banner", "Banner baru berhasil di unggah.", "success", session)
	pathRoot := fmt.Sprintf("/dasboard/banner")
	return ctx.Redirect(pathRoot)
}

func (h *adminWebBannerHandler) UpdateBannerView(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
	} else {

		sessionResult = session.Get("msg-alert-update-banner")
		staussession = session.Get("msg-alert-update-banner-status")
	}

	alert := helper.AlertString(sessionResult, staussession)

	// cookisUserId := ctx.Cookies("sessionLog")
	// idUser, err := helper.GetSessionID(cookisUserId, h.config)
	// if err != nil {
	// 	return ctx.Redirect("/login")
	// }

	var idUser uint = 14

	userMain, err := h.userService.GetUserByID(c, idUser)
	if err != nil {
		return ctx.Redirect("/login")
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Redirect("/banner")
	}

	banner, _ := h.service.GetBannerByID(c, uint(id))
	// if err != nil {
	// 	return ctx.Render("error400", fiber.Map{})
	// }

	return ctx.Render("dasboard/admin/banner/dasboard_banner_update", fiber.Map{
		"header": userMain,
		"data":   banner,
		"layout": "table",
		"alert":  template.HTML(alert),
	})
}

// func (h *adminWebBannerHandler) UpdateBanner(ctx *fiber.Ctx) error {

// }

func (h *adminWebBannerHandler) DeleteBanner(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-banner", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-banner", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.service.DeleteBanner(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-banner", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-banner", "Banner berhasil di hapus!", "success", session)
	return ctx.Redirect(root)
}
