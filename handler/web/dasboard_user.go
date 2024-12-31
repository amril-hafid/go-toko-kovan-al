package web

import (
	"context"
	"fmt"
	"go-toko-kovan-al/app/users"
	"go-toko-kovan-al/config"
	"go-toko-kovan-al/helper"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type adminWebHandler struct {
	userService  users.Service
	sessionStore *session.Store
	config       config.Config
}

func NewAdminWebHandler(userService users.Service, sessionStore *session.Store, config config.Config) *adminWebHandler {
	return &adminWebHandler{userService, sessionStore, config}
}

func (h *adminWebHandler) ShowAllUser(ctx *fiber.Ctx) error {

	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
	} else {

		sessionResult = session.Get("msg-alert-new-user")
		staussession = session.Get("msg-alert-new-user-status")
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

	users, err := h.userService.GetAllUsers(c)
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	return ctx.Render("dasboard/admin/user/dasboard_user_list", fiber.Map{
		"header": userMain,
		"role":   userMain.Role,
		"data":   users,
		"layout": "table",
		"alert":  template.HTML(alert),
	})
}
func (h *adminWebHandler) ShowUserProfile(ctx *fiber.Ctx) error {
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

	return ctx.Render("dasboard/admin/user/dasboard_user_profile", fiber.Map{
		"header": userMain,
		"role":   userMain.Role,
		"data":   userMain,
		"layout": "table",
	})
}
func (h *adminWebHandler) ShowDetailUser(ctx *fiber.Ctx) error {
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
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect("/dasboard/user")
	}

	user, err := h.userService.GetUserByID(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect("/dasboard/user")
	}

	return ctx.Render("dasboard/admin/user/dasboard_user_profile", fiber.Map{
		"header": userMain,
		"role":   userMain.Role,
		"data":   user,
		"layout": "table",
	})
}
func (h *adminWebHandler) NewUser(ctx *fiber.Ctx) error {

	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()
	root := ctx.Get("Referer")

	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	inputNew := new(users.RegisterUserInput)

	if err := ctx.BodyParser(inputNew); err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.userService.RegisterUser(c, *inputNew)
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-user", "Resgistrasi Pengguna baru berhasil.", "success", session)
	pathRoot := fmt.Sprintf("/dasboard/user")
	return ctx.Redirect(pathRoot)
}
func (h *adminWebHandler) UploadImageProfileView(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
	} else {
		sessionResult = session.Get("msg-alert-new-user")
		staussession = session.Get("msg-alert-new-user-status")
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
		return ctx.Redirect("/login")
	}

	userRow, err := h.userService.GetUserByID(c, uint(id))
	if err != nil {
		return ctx.Redirect("/dasboard/user")
	}
	return ctx.Render("dasboard/admin/user/dasboard_user_image_profile", fiber.Map{
		"header": userMain,
		"data":   userRow,
		"layout": "form",
		"alert":  template.HTML(alert),
	})
}
func (h *adminWebHandler) UploadImageProfile(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, _ := h.sessionStore.Get(ctx)
	id := ctx.FormValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	file, err := ctx.FormFile("image_avatar")
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	a, _ := helper.AllowedFileType(file, h.config)
	if a != true {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	fileName := helper.StringWithoutSpaces(file.Filename)

	path := fmt.Sprintf("%s%s", h.config.FileConf.ImageDerektory, fileName)

	userRow, err := h.userService.GetUserByID(c, uint(idInt))
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	// update image profile
	if userRow.ProfileFile != "" {
		err = ctx.SaveFile(file, path)
		if err != nil {
			helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
			return ctx.Redirect(root)
		}

		_, err = h.userService.CreateImageProfile(c, fileName, userRow.ID)
		if err != nil {
			helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
			return ctx.Redirect(root)
		}

		pathRemoveFile := fmt.Sprintf("%s%s", h.config.FileConf.ImageDerektory, userRow.ProfileFile)
		err = os.Remove(pathRemoveFile)
		if err != nil {
			helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
			return ctx.Redirect(root)
		}

		helper.AlertMassage("msg-alert-new-user", "Image profile berhasil di perbarui!", "success", session)
		return ctx.Redirect("/dasboard/admin/user")
	}
	// end update image profile

	err = ctx.SaveFile(file, path)
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.userService.CreateImageProfile(c, fileName, userRow.ID)
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-user", "Image profile berhasil di upload !", "success", session)
	return ctx.Redirect("/dasboard/admin/user")
}
func (h *adminWebHandler) UpdateUserView(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()

	} else {
		sessionResult = session.Get("msg-alert-update-user")
		staussession = session.Get("msg-alert-update-user-status")
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

	fmt.Println("data user :", userMain)

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	userRow, err := h.userService.GetUserByID(c, uint(id))
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	return ctx.Render("dasboard/admin/user/dasboard_user_update", fiber.Map{
		"header": userMain,
		"role":   userMain.Role,
		"data":   userRow,
		"layout": "form",
		"alert":  template.HTML(alert),
	})
}
func (h *adminWebHandler) UpdateUser(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()
	root := ctx.Get("Referer")
	input := new(users.UpdateUserInput)
	session, err := h.sessionStore.Get(ctx)
	if err != nil {
		helper.AlertMassage("msg-alert-update-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	cookisUserId := ctx.Cookies("sessionLog")
	idUser, err := helper.GetSessionID(cookisUserId, h.config)
	if err != nil {
		return ctx.Redirect("/login")
	}

	if err = ctx.BodyParser(input); err != nil {
		helper.AlertMassage("msg-alert-update-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if input.Role == "admin" {
		if _, err = h.userService.UpdateUser(c, *input); err != nil {
			helper.AlertMassage("msg-alert-update-user", err.Error(), "error", session)
			return ctx.Redirect(root)
		}
	} else if input.Role == "admin" {

		if input.ID == idUser {
			input.Role = "user"
			if _, err = h.userService.UpdateUser(c, *input); err != nil {
				helper.AlertMassage("msg-alert-update-user", err.Error(), "error", session)
				return ctx.Redirect(root)
			}
		}
	} else {
		helper.AlertMassage("msg-alert-update-user", "Update data gagal", "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-update-user", "Update User berhasil.", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebHandler) UpdatePasswordView(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()

	} else {
		sessionResult = session.Get("msg-alert-update-password")
		staussession = session.Get("msg-alert-update-password-status")
	}

	// end from error
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
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect("/dasboard/user")
	}

	userRow, err := h.userService.GetUserByID(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect("/dasboard/user")
	}

	return ctx.Render("dasboard/admin/user/dasboard_user_password_reset", fiber.Map{
		"header": userMain,
		"data":   userRow,
		"layout": "form",
		"alert":  template.HTML(alert),
	})
}
func (h *adminWebHandler) ResetPassword(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	input := new(users.UpdatePasswordInput)
	session, err := h.sessionStore.Get(ctx)
	if err != nil {
		helper.AlertMassage("msg-alert-update-password", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if err = ctx.BodyParser(input); err != nil {
		helper.AlertMassage("msg-alert-update-password", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if _, err = h.userService.UpdatePassword(c, *input); err != nil {
		helper.AlertMassage("msg-alert-update-password", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-update-password", "Password User berhasil di perbarui.", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebHandler) DeleteUserSoft(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.userService.DeleteUserSoft(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-user", "User berhasil di hapus!", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebHandler) DeleteUserRecycle(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.userService.DeleteUser(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-user", "User berhasil di hapus!", "success", session)
	return ctx.Redirect(root)
}
func (h *adminWebHandler) ShowAdminAllRecycle(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
	} else {

		sessionResult = session.Get("msg-alert-new-user")
		staussession = session.Get("msg-alert-new-user-status")
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

	users, err := h.userService.GetAllUsersDeleted(c)
	if err != nil {
		return ctx.Render("error400", fiber.Map{})
	}

	return ctx.Render("dasboard/admin/user/dasboard_user_list_recycle", fiber.Map{
		"header": userMain,
		"data":   users,
		"layout": "table",
		"alert":  template.HTML(alert),
	})
}
func (h *adminWebHandler) RestoreUser(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)

	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	id, err := ctx.ParamsInt("id")
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	_, err = h.userService.RestoreUser(c, uint(id))
	if err != nil {
		helper.AlertMassage("msg-alert-new-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.AlertMassage("msg-alert-new-user", "User berhasil di hapus!", "success", session)
	return ctx.Redirect(root)
}
