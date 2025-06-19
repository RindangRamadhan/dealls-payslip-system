package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/request"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary     Get Payroll Summary
// @Description Get payroll summary for a given period (via query param)
// @ID          get-payroll-summary
// @Tags  	    Admin
// @Accept      json
// @Produce     json
// @Param       period_id query string true "Attendance Period ID (UUID)"
// @Success     200 {object} utils.BodySuccess{data=response.PayrollSummaryResponse}
// @Failure     400 {object} utils.BodyFailure
// @Failure     500 {object} utils.BodyFailure
// @Security    JWTBearer
// @Router      /admin/payrolls/summary [get]
func (r *V1) getPayrollSummary(ctx *fiber.Ctx) error {
	periodIDStr := ctx.Query("period_id")

	// Validasi UUID
	periodID, err := uuid.Parse(periodIDStr)
	if err != nil {
		r.logger.Error(err, "http - v1 - getPayrollSummary - uuid.Parse")
		return utils.WriteError(ctx, http.StatusBadRequest, "invalid uuid format", nil)
	}

	summary, err := r.payroll.GetPayrollSummary(ctx.UserContext(), periodID)
	if err != nil {
		r.logger.Error(err, "http - v1 - getPayrollSummary - getPayrollSummary - r.u.GetPayrollSummary")

		if strings.Contains(err.Error(), "no rows in result set") {
			return utils.WriteError(ctx, http.StatusNotFound, "payroll summary not found", nil)
		}

		return utils.WriteError(ctx, http.StatusBadRequest, "failed to get payroll summary", nil)
	}

	return utils.WriteSuccess(ctx, "Payroll summary retrieved successfully", summary)
}

// @Summary     Process Payroll
// @Description Generate payroll for a specific attendance period. Requires authentication.
// @ID          process-payroll
// @Tags        Admin
// @Accept      json
// @Produce     json
// @Param       payload body request.PayrollRequest true "Payroll request payload"
// @Success     200 {object} utils.BodySuccess{data=entity.Payroll}
// @Failure     400 {object} utils.BodyFailure
// @Failure     500 {object} utils.BodyFailure
// @Security    JWTBearer
// @Router      /admin/payrolls [post]
func (r *V1) processPayroll(ctx *fiber.Ctx) error {
	var req request.PayrollRequest
	if err := ctx.BodyParser(&req); err != nil {
		r.logger.Error(err, "http - v1 - processPayroll - BodyParser")
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid request body", nil)
	}

	// Validation
	if err := utils.ValidateStruct(req); err != nil {
		r.logger.Error(err, "http - v1 - processPayroll - ValidateStruct")
		return utils.WriteError(ctx, http.StatusBadRequest, "Bad request", err)
	}

	periodID, err := uuid.Parse(req.PeriodID)
	if err != nil {
		r.logger.Error(err, "http - v1 - uuid.Parse")
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid UUID format for period_id", nil)
	}

	userIDRaw := ctx.Locals("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		r.logger.Error(err, "http - v1 - uuid.Assertion")
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid user ID in context", nil)
	}

	payroll, err := r.payroll.ProcessPayroll(ctx.UserContext(), periodID, userID)
	if err != nil {
		r.logger.Error(err, "http - v1 - processPayroll - r.payroll.ProcessPayroll")

		if strings.Contains(err.Error(), "already processed") {
			return utils.WriteError(ctx, http.StatusConflict, err.Error(), nil)
		}

		return utils.WriteError(ctx, http.StatusInternalServerError, "Failed to process payroll", err)
	}

	return utils.WriteSuccess(ctx, "Payroll processed successfully", payroll)
}

// @Summary     Create Attendance Period
// @Description Create a new attendance period by providing start and end dates.
// @ID          create-attendance-period
// @Tags  	    Admin
// @Accept      json
// @Produce     json
// @Param       request body request.AttendancePeriodRequest true "Attendance Period Payload"
// @Success     201 {object} utils.BodySuccess{data=entity.AttendancePeriod}
// @Failure     400 {object} utils.BodyFailure
// @Failure     500 {object} utils.BodyFailure
// @Security    JWTBearer
// @Router      /admin/attendance-periods [post]
func (r *V1) createAttendancePeriod(ctx *fiber.Ctx) error {
	var req request.AttendancePeriodRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid request body", err)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.WriteError(ctx, http.StatusBadRequest, "Validation failed", err)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid start date format", err)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid end date format", err)
	}

	if endDate.Before(startDate) {
		return utils.WriteError(ctx, http.StatusBadRequest, "End date must be after start date", nil)
	}

	userIDRaw := ctx.Locals("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid user ID", nil)
	}

	period := &entity.AttendancePeriod{
		ID:        uuid.New(),
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	if err := r.attendancePeriod.CreateAttendancePeriod(ctx.UserContext(), period); err != nil {
		return utils.WriteError(ctx, http.StatusInternalServerError, "Failed to create attendance period", err)
	}

	// store record_id for audit log
	ctx.Locals("record_id", period.ID)

	return utils.WriteSuccess(ctx, "Attendance period created successfully", period)
}
