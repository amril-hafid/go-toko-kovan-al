package web

import (
	"context"
	"fmt"
	"go-toko-kovan-al/app/category"
	"go-toko-kovan-al/app/users"
	"go-toko-kovan-al/config"
	"go-toko-kovan-al/helper"
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type adminWebCategoryHandler struct {
	categoryService category.Service
	userService     users.Service
	sessionStore    *session.Store
	config          config.Config
}

func NewAdminWebCategoryHandler(categoryService category.Service, userService users.Service, sessionStore *session.Store, config config.Config) *adminWebCategoryHandler {
	return &adminWebCategoryHandler{categoryService, userService, sessionStore, config}
}

func (h *adminWebCategoryHandler) ShowAllCategory(ctx *fiber.Ctx) error {

	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
	} else {

		sessionResult = session.Get("msg-alert-new-category")
		staussession = session.Get("msg-alert-new-category-status")
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

	category, err := h.categoryService.GetAllCategory(c)
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	return ctx.Render("dasboard/admin/category/dasboard_category_list", fiber.Map{
		"header": userMain,
		"data":   category,
		"layout": "table",
		"alert":  template.HTML(alert),
	})
}
func (h *adminWebCategoryHandler) ShowDetailCategory(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

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

	category, err := h.categoryService.GetCategoryByID(c, uint(id))
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	return ctx.Render("dasboard/admin/category/dasboard_category_detail", fiber.Map{
		"header": userMain,
		"data":   category,
		"layout": "table",
	})
}
func (h *adminWebCategoryHandler) NewCategory(ctx *fiber.Ctx) error {

	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()
	root := ctx.Get("Referer")

	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	inputNew := new(category.InputCategory)

	if err := ctx.BodyParser(inputNew); err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}
	fmt.Printf("data input handler 1 :", inputNew)
	a, err := h.categoryService.CreateCategory(c, *inputNew)
	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	fmt.Printf("data input handler 2 :", a)

	helper.AlertMassage("msg-alert-new-category", "Category baru berhasil di buat.", "success", session)
	pathRoot := fmt.Sprintf("/dasboard/category")
	return ctx.Redirect(pathRoot)
}
func (h *adminWebCategoryHandler) UpdateCategoryView(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()

	} else {
		sessionResult = session.Get("msg-alert-update-category")
		staussession = session.Get("msg-alert-update-category-status")
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
	categoryRow, err := h.categoryService.GetCategoryByID(c, uint(id))
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	return ctx.Render("dasboard/admin/category/dasboard_category_update", fiber.Map{
		"header": userMain,
		"data":   categoryRow,
		"layout": "form",
		"alert":  template.HTML(alert),
	})
}
func (h *adminWebCategoryHandler) UpdateCategory(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()
	root := ctx.Get("Referer")
	input := new(category.UpdateCategory)
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-update-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if err = ctx.BodyParser(input); err != nil {
		helper.AlertMassage("msg-alert-update-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if _, err = h.categoryService.UpdateCategory(c, *input); err != nil {
		helper.AlertMassage("msg-alert-update-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-update-category", "Update category berhasil.", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebCategoryHandler) DeleteCategorySoft(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.categoryService.DeleteCategorySoft(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-category", "Category berhasil di hapus!", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebCategoryHandler) DeleteCategoryRecycle(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.categoryService.DeleteCategory(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-category", "Category berhasil di hapus!", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebCategoryHandler) ShowAllCategoryRecycle(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
	} else {

		sessionResult = session.Get("msg-alert-new-category")
		staussession = session.Get("msg-alert-new-user-category")
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

	categorys, _ := h.categoryService.GetAllCategoryDeleted(c)
	// if err != nil {
	// 	return ctx.Render("error400", fiber.Map{})
	// }

	return ctx.Render("dasboard/admin/category/dasboard_category_list_recycle", fiber.Map{
		"header": userMain,
		"data":   categorys,
		"layout": "table",
		"alert":  template.HTML(alert),
	})
}
func (h *adminWebCategoryHandler) RestoreCategory(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.categoryService.RestoreCategory(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-category", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-category", "category berhasil di kembalikan!", "success", session)
	return ctx.Redirect(root)
}
