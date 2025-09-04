// pkg/api/routes.go
package api

import (
	"github.com/dhruvpatel-10/signee/ca-api/db"
	"github.com/dhruvpatel-10/signee/ca-api/internal/service/auth"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, q *db.Queries) {

	v1 := r.Group("/api/v1")
	authService := &auth.AuthService{DB: q}
	// Public endpoints
	public := v1.Group("/")
	{
		public.POST("/auth/login", authService.Login)
		// public.POST("/auth/refresh", handlers.RefreshToken)
		public.POST("/auth/signup", authService.Signup)
		// public.GET("/healthz", auth.AuthService.HealthCheck)
	}

	// Protected endpoints
	// protected := v1.Group("/")
	// protected.Use(middleware.AuthRequired())
	// {
	// 	// User management
	// 	protected.GET("/users/me", handlers.GetCurrentUser)
	// 	protected.PUT("/users/me", handlers.UpdateCurrentUser)
	// 	protected.POST("/users/me/mfa/enable", handlers.EnableMFA)

	// 	// Certificate Authorities
	// 	protected.GET("/cas", handlers.ListCAs)
	// 	protected.POST("/cas", middleware.RequirePermission("ca:create"), handlers.CreateCA)
	// 	protected.GET("/cas/:id", handlers.GetCA)
	// 	protected.PUT("/cas/:id", middleware.RequirePermission("ca:update"), handlers.UpdateCA)
	// 	protected.POST("/cas/:id/rotate", middleware.RequirePermission("ca:rotate"), handlers.RotateCA)

	// 	// Certificate Templates
	// 	protected.GET("/templates", handlers.ListTemplates)
	// 	protected.POST("/templates", middleware.RequirePermission("template:create"), handlers.CreateTemplate)

	// 	// Certificates
	// 	protected.GET("/certificates", handlers.ListCertificates)
	// 	protected.POST("/certificates/request", handlers.RequestCertificate)
	// 	protected.GET("/certificates/:id", handlers.GetCertificate)
	// 	protected.POST("/certificates/:id/revoke", middleware.RequirePermission("cert:revoke"), handlers.RevokeCertificate)

	// 	// Certificate Requests (approval workflow)
	// 	protected.GET("/requests", handlers.ListRequests)
	// 	protected.POST("/requests/:id/approve", middleware.RequirePermission("request:approve"), handlers.ApproveRequest)
	// 	protected.POST("/requests/:id/reject", middleware.RequirePermission("request:approve"), handlers.RejectRequest)

	// 	// Audit logs
	// 	protected.GET("/audit", middleware.RequirePermission("audit:read"), handlers.GetAuditLogs)

	// 	// Admin endpoints
	// 	admin := protected.Group("/admin")
	// 	admin.Use(middleware.RequireRole("admin"))
	// 	{
	// 		admin.GET("/users", handlers.ListUsers)
	// 		admin.POST("/users", handlers.CreateUser)
	// 		admin.PUT("/users/:id/roles", handlers.UpdateUserRoles)
	// 	}
	// }
}
