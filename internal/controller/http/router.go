// Package v1 implements routing paths. Each services in own file.
package http

import (
	"github.com/RindangRamadhan/dealls-payslip-system/config"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/usecase"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/RindangRamadhan/dealls-payslip-system/docs" // Swagger docs.
	v1 "github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1"
)

// NewRouter initializes the API routes.

// @title           Dealls Payslip System API
// @version         1.0
// @description     API for the Dealls Payslip System, featuring attendance tracking, overtime, reimbursement, and payroll summary.
// @host            localhost:8080
// @BasePath        /v1

// @securityDefinitions.apikey JWTBearer
// @in                         header
// @name                       Authorization
// @description                Provide your JWT token prefixed with "Bearer", e.g., "Bearer <your_token>"
func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	lgr logger.Interface,
	user usecase.User,
	payroll usecase.Payroll,
	payslip usecase.Payslip,
	overtime usecase.Overtime,
	attendance usecase.Attendance,
	reimbursement usecase.Reimbursement,
	attendancePeriod usecase.AttendancePeriod,
) {
	// Prometheus Metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("dealls-payslip-system")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	app.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendStatus(200) })

	// V1 Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewTranslationRoutes(
			apiV1Group,
			lgr,
			user,
			payroll,
			payslip,
			overtime,
			attendance,
			reimbursement,
			attendancePeriod,
		)
	}
}
