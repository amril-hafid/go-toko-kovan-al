package router

import (
	"context"
	"go-toko-kovan-al/app/auth"
	"go-toko-kovan-al/app/banner"
	"go-toko-kovan-al/app/category"
	"go-toko-kovan-al/app/image_product"
	"go-toko-kovan-al/app/product"
	"go-toko-kovan-al/app/users"
	"time"

	"go-toko-kovan-al/config"
	"go-toko-kovan-al/handler/web"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, validate *validator.Validate, store *session.Store, conf config.Config) {
	// Repository
	userRepository := users.NewRepository(db)
	productRepository := product.NewRepository(db)
	imageProductRepository := image_product.NewRepository(db)
	caregoryProductRepository := category.NewRepository(db)
	bannerRepository := banner.NewRepository(db)
	//End Repository

	// Service
	aurhService := auth.NewService()
	userServiceWeb := users.NewService(userRepository, validate, db)                                                                  // Users Service
	productServiceWeb := product.NewService(productRepository, caregoryProductRepository, imageProductRepository, validate, db, conf) // Product Service
	caregoryProductService := category.NewService(caregoryProductRepository, validate, db)                                            // Category Service
	imageProductService := image_product.NewService(imageProductRepository, validate, db)                                             // Image Product Service
	bannerService := banner.NewService(bannerRepository, validate, db)
	//End Service

	// Handler
	webAdminHandler := web.NewAdminWebHandler(userServiceWeb, store, conf)
	webProduct := web.NewProductHandler(productServiceWeb, bannerService, store, conf)
	sessionHandler := web.NewSessionHandler(userServiceWeb, aurhService, store, conf)
	categoryProdactHandler := web.NewAdminWebCategoryHandler(caregoryProductService, userServiceWeb, store, conf)
	productHandler := web.NewAdminWebProductHandler(productServiceWeb, caregoryProductService, userServiceWeb, store, conf)
	imageHandler := web.NewAdminWebImageProductHandler(imageProductService, productServiceWeb, userServiceWeb, store, conf)
	bannerHandler := web.NewAdminWebBannerHandler(bannerService, userServiceWeb, store, conf)

	//End Handler

	app.Static("/tem", "./public/template/template-admin/AdminLTE/")
	app.Static("/themes", "./public/template/template-web/")
	app.Static("/image", "./assets/")

	// web conten
	app.Get("/", webProduct.Home)
	app.Get("/shop", webProduct.Shop)
	app.Get("/shop-detail/:id", webProduct.ShopDetail)
	app.Get("/contact", webProduct.Contact)

	// Auth Session
	app.Get("/login", LoginMiddleware(aurhService, userServiceWeb, conf), sessionHandler.LoginView)
	app.Post("/session", LoginMiddleware(aurhService, userServiceWeb, conf), sessionHandler.Login)
	app.Get("/logout", sessionHandler.Destroy)
	//End Auth Session

	// dasboard
	dasboard := app.Group("/dasboard")

	// profile
	dasboard.Get("/profile", webAdminHandler.ShowUserProfile)

	// User
	adminUser := dasboard.Group("/user")
	adminUser.Get("", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.ShowAllUser)
	// adminUser.Get("/detail/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.ShowDetailUser)
	adminUser.Post("/new", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.NewUser)
	adminUser.Get("/update/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.UpdateUserView)
	adminUser.Post("/update", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.UpdateUser)
	dasboard.Get("/detail/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.ShowDetailUser)
	// adminUser.Get("/profile/image/:id", webAdminHandler.UploadImageProfileView)
	// adminUser.Post("/profile/image/", webAdminHandler.UploadImageProfile)
	adminUser.Get("/new-password/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.UpdatePasswordView)
	adminUser.Post("/new-password", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.ResetPassword)
	adminUser.Get("/delete/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.DeleteUserSoft)
	adminUser.Get("/recycle", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.ShowAdminAllRecycle)
	adminUser.Get("/recycle/restore/:id", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.RestoreUser)
	adminUser.Get("/recycle/delete/:id", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), webAdminHandler.DeleteUserRecycle)
	// end User

	// Category
	adminCategory := dasboard.Group("/category")
	adminCategory.Get("", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), categoryProdactHandler.ShowAllCategory)
	adminCategory.Get("/detail/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), categoryProdactHandler.ShowDetailCategory)
	adminCategory.Post("/new", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), categoryProdactHandler.NewCategory)
	adminCategory.Get("/update/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), categoryProdactHandler.UpdateCategoryView)
	adminCategory.Post("/update", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), categoryProdactHandler.UpdateCategory)
	adminCategory.Get("/delete/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), categoryProdactHandler.DeleteCategorySoft)
	adminCategory.Get("/recycle", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), categoryProdactHandler.ShowAllCategoryRecycle)
	adminCategory.Get("/recycle/restore/:id", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), categoryProdactHandler.RestoreCategory)
	adminCategory.Get("/recycle/delete/:id", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), categoryProdactHandler.DeleteCategoryRecycle)
	// End Category

	// Product
	adminProduct := dasboard.Group("/product")
	adminProduct.Get("", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), productHandler.ShowAllProduct)
	adminProduct.Get("/detail/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), productHandler.ShowDetailProduct)
	adminProduct.Post("/new", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), productHandler.NewProduct)
	adminProduct.Get("/image/new/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), imageHandler.NewImageProductView)
	adminProduct.Post("/image/new", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), imageHandler.NewImageProduct)
	adminProduct.Get("/update/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), productHandler.UpdateProductView)
	adminProduct.Post("/update", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), productHandler.UpdateProduct)
	adminProduct.Get("/delete/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), productHandler.DeleteProductSoft)
	adminProduct.Get("/recycle", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), productHandler.ShowAllProductRecycle)
	adminProduct.Get("/recycle/restore/:id", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), productHandler.RestoreProduct)
	adminProduct.Get("/recycle/delete/:id", authRoleAdminMiddleware(aurhService, userServiceWeb, conf), productHandler.DeleteProductRecycle)
	// // End Product

	// Image Product
	adminProduct.Get("/image/delete/:id", authRoleUserEndAdminMiddleware(aurhService, userServiceWeb, conf), imageHandler.DeleteImageProduct)

	// banner
	adminBanner := dasboard.Group("/banner")
	adminBanner.Get("", bannerHandler.ShowAllBanner)
	adminBanner.Post("/new", bannerHandler.NewBanner)
	adminBanner.Get("/delete/:id", bannerHandler.DeleteBanner)

}

func LoginMiddleware(authService auth.Service, userSession users.Service, conf config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, cencel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cencel()
		cookisUserId := c.Cookies("sessionLog")

		token, err := authService.ValidateToken(cookisUserId, conf)
		if err != nil {
			return c.Next()
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Next()
		}

		idUser := uint(claim["user_id"].(float64))
		exp := claim["exp"].(float64)
		if time.Now().Unix() > int64(exp) {
			c.ClearCookie("sessionLog")
			c.ClearCookie()
			return c.Next()
		}

		user, err := userSession.GetUserByID(ctx, idUser)
		if err != nil {
			return c.Next()
		}

		if user.Role == "admin" || user.Role == "user" {
			return c.Redirect("/dasboard")
		}

		return c.Next()
	}
}
func authMiddleware(authService auth.Service, userSession users.Service, conf config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, cencel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cencel()
		cookisUserId := c.Cookies("sessionLog")

		token, err := authService.ValidateToken(cookisUserId, conf)
		if err != nil {
			return c.Redirect("/login")
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Redirect("/login")
		}

		idUser := uint(claim["user_id"].(float64))

		exp := claim["exp"].(float64)
		if time.Now().Unix() > int64(exp) {
			c.ClearCookie("sessionLog")
			c.ClearCookie()
			return c.Next()
		}
		_, err = userSession.GetUserByID(ctx, idUser)
		if err != nil {
			return c.Redirect("/login")
		}

		return c.Next()
	}
}
func authRoleAdminMiddleware(authService auth.Service, userSession users.Service, conf config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, cencel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cencel()
		cookisUserId := c.Cookies("sessionLog")

		token, err := authService.ValidateToken(cookisUserId, conf)
		if err != nil {
			return c.Redirect("/login")
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Redirect("/login")
		}

		idUser := uint(claim["user_id"].(float64))
		exp := claim["exp"].(float64)
		if time.Now().Unix() > int64(exp) {
			c.ClearCookie("sessionLog")
			c.ClearCookie()
			return c.Next()
		}

		user, err := userSession.GetUserByID(ctx, idUser)
		if err != nil {
			return c.Redirect("/login")
		}

		if user.Role == "user" {
			return c.Redirect("/dasboard")
		}

		return c.Next()
	}
}
func authRoleUserEndAdminMiddleware(authService auth.Service, userSession users.Service, conf config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, cencel := context.WithTimeout(c.Context(), 2*time.Second)
		defer cencel()
		cookisUserId := c.Cookies("sessionLog")
		token, err := authService.ValidateToken(cookisUserId, conf)
		if err != nil {
			return c.Redirect("/login")
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Redirect("/login")
		}

		idUser := uint(claim["user_id"].(float64))
		exp := claim["exp"].(float64)
		if time.Now().Unix() > int64(exp) {
			c.ClearCookie("sessionLog")
			c.ClearCookie()
			return c.Redirect("/login")

		}

		user, err := userSession.GetUserByID(ctx, idUser)
		if err != nil {
			return c.Redirect("/login")
		}

		if user.Role == "admin" || user.Role == "user" {
			return c.Next()
		}

		return c.Redirect("/login")
	}
}
