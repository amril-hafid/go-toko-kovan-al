package web

import (
	"context"
	"go-toko-kovan-al/app/auth"
	"go-toko-kovan-al/app/users"
	"go-toko-kovan-al/config"
	"go-toko-kovan-al/helper"
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type sessionHandler struct {
	userService  users.Service
	authService  auth.Service
	sessionStore *session.Store
	config       config.Config
}

func NewSessionHandler(userService users.Service, authService auth.Service, sessionStore *session.Store, config config.Config) *sessionHandler {
	return &sessionHandler{userService, authService, sessionStore, config}
}

func (h *sessionHandler) LoginView(ctx *fiber.Ctx) error {
	_, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var sessionResult, staussession interface{}
	session, err := h.sessionStore.Get(ctx)
	defer session.Save()

	if err != nil {
		sessionResult = err.Error()
		session.Destroy()
	} else {
		sessionResult = session.Get("msg-alert-login-user")
		staussession = session.Get("msg-alert-login-user-status")
	}

	alert := helper.AlertString(sessionResult, staussession)
	return ctx.Render("dasboard/session/login", fiber.Map{
		"alert": template.HTML(alert),
	})
}
func (h *sessionHandler) Login(ctx *fiber.Ctx) error {
	c, cencel := context.WithTimeout(ctx.Context(), 2*time.Second)
	defer cencel()

	var input users.LoginInput
	root := ctx.Get("Referer")
	session, err := h.sessionStore.Get(ctx)
	if err != nil {
		helper.AlertMassage("msg-alert-login-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	if err := ctx.BodyParser(&input); err != nil {
		helper.AlertMassage("msg-alert-login-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	rowUser, err := h.userService.Login(c, input)
	if err != nil {
		helper.AlertMassage("msg-alert-login-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	idJWT, err := h.authService.GenerateToken(rowUser.ID, h.config)
	if err != nil {
		helper.AlertMassage("msg-alert-login-user", err.Error(), "error", session)
		return ctx.Redirect(root)
	}

	helper.FormDeleteError("msg-alert-login-user", session)

	cookie := new(fiber.Cookie)
	cookie.Name = "sessionLog"
	cookie.Value = idJWT
	cookie.Expires = time.Now().Add(24 * time.Hour)
	ctx.Cookie(cookie)

	return ctx.Redirect("/dasboard/user")
}
func (h *sessionHandler) Destroy(c *fiber.Ctx) error {

	cookisUserId := c.Cookies("sessionLog")
	idUser, err := helper.GetSessionID(cookisUserId, h.config)
	if err != nil {
		return c.Redirect("/login")
	}

	if idUser == 0 {
		return c.Redirect("/login")
	}

	c.ClearCookie("sessionLog")
	c.ClearCookie()

	return c.Redirect("/login")
}
