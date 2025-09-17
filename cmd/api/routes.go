package main

func (app *Application) routes() {
	apiGroup := app.Server.Group("/api/v1")
	publicAuthRoutes := apiGroup.Group("/auth")
	{
		publicAuthRoutes.POST("/user", app.handler.RegisterHandler)
		publicAuthRoutes.POST("/login", app.handler.LoginHandler)
		publicAuthRoutes.POST("/forgot/password", app.handler.ForgotPassword)
		publicAuthRoutes.POST("/reset/password", app.handler.ResetPassword)
	}

	profileRoutes := apiGroup.Group("/", app.Middleware.AuthenticationMiddleware)
	{
		profileRoutes.GET("/authenticated/user", app.handler.GetAuthenticatedUser)
		profileRoutes.PATCH("/change/password", app.handler.ChangePassword)
	}

	categoryRoutes := apiGroup.Group("/categories", app.Middleware.AuthenticationMiddleware)
	{
		categoryRoutes.GET("/all", app.handler.ListCategories)
		categoryRoutes.POST("/", app.handler.CreateCategory)
		categoryRoutes.DELETE("/:id", app.handler.DeleteCategory)
	}

	app.Server.GET("/", app.handler.HealthCheck)

}
