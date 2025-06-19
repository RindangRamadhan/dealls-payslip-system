package v1

import (
	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/middleware"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/usecase"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(
	api fiber.Router,
	logger logger.Interface,
	user usecase.User,
	payroll usecase.Payroll,
	payslip usecase.Payslip,
	overtime usecase.Overtime,
	attendance usecase.Attendance,
	reimbursement usecase.Reimbursement,
	attendancePeriod usecase.AttendancePeriod,
) {
	r := &V1{
		logger: logger,
		validator: validator.New(
			validator.WithRequiredStructEnabled(),
		),

		// Initialize the usecases
		user:             user,
		payroll:          payroll,
		payslip:          payslip,
		overtime:         overtime,
		attendance:       attendance,
		reimbursement:    reimbursement,
		attendancePeriod: attendancePeriod,
	}

	// Auth Routes
	api.Post("/auth/login", r.login)

	// Admin Routes
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.AdminOnly())
	{
		adminGroup.Post("/attendance-periods", r.createAttendancePeriod)
		adminGroup.Post("/payrolls", r.processPayroll)
		adminGroup.Get("/payrolls/summary", r.getPayrollSummary)
	}

	// Employee Routes
	employeeGroup := api.Group("/employee")
	{
		employeeGroup.Get("/payslips", r.getPayslip)
		employeeGroup.Post("/attendances", r.submitAttendance)
		employeeGroup.Post("/overtimes", r.submitOvertime)
		employeeGroup.Post("/reimbursements", r.submitReimbursement)
	}
}
