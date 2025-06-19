package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/request"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// @Summary     Get Payslip
// @Description Retrieve payslip for the authenticated employee.
// @ID          get-payslip
// @Tags  	    Employee
// @Accept      json
// @Produce     json
// @Param       period_id query string true "Attendance Period ID (UUID)"
// @Success     200 {object} utils.BodySuccess{data=response.PayslipResponse} "payslip retrieved successfully"
// @Failure     400 {object} utils.BodyFailure "invalid input or bad request"
// @Failure     404 {object} utils.BodyFailure "payslip not found"
// @Failure     500 {object} utils.BodyFailure "internal server error"
// @Security    JWTBearer
// @Router      /employee/payslips [get]
func (r *V1) getPayslip(ctx *fiber.Ctx) error {
	periodIDStr := ctx.Query("period_id")

	// Validasi UUID
	periodID, err := uuid.Parse(periodIDStr)
	if err != nil {
		r.logger.Error(err, "http - v1 - getPayslip - uuid.Parse")
		return utils.WriteError(ctx, http.StatusBadRequest, "invalid uuid format", nil)
	}

	userIDRaw := ctx.Locals("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		r.logger.Error(err, "http - v1 - getPayslip - uuid.Assertion")
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid user ID in context", nil)
	}

	payslip, err := r.payslip.GetPayslip(ctx.UserContext(), userID, periodID)
	if err != nil {
		r.logger.Error(err, "http - v1 - getPayslip - r.u.GetPayslip")

		if strings.Contains(err.Error(), "no rows in result set") {
			return utils.WriteError(ctx, http.StatusNotFound, "payslip not found", nil)
		}

		return utils.WriteError(ctx, http.StatusBadRequest, "failed to get payslip", nil)
	}

	return utils.WriteSuccess(ctx, "payslip retrieved successfully", payslip)
}

// @Summary     Submit Attendance
// @Description Submit attendance (check-in and check-out) for current active period.
// @ID          submit-attendance
// @Tags        Employee
// @Accept      json
// @Produce     json
// @Param       payload body request.AttendanceRequest true "Attendance payload"
// @Success     201 {object} utils.BodySuccess
// @Failure     400 {object} utils.BodyFailure
// @Failure     409 {object} utils.BodyFailure
// @Failure     422 {object} utils.BodyFailure
// @Failure     500 {object} utils.BodyFailure
// @Security    JWTBearer
// @Router      /employee/attendances [post]
func (r *V1) submitAttendance(ctx *fiber.Ctx) error {
	var req request.AttendanceRequest
	if err := ctx.BodyParser(&req); err != nil {
		r.logger.Error(err, "http - v1 - submitAttendance - BodyParser")
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := utils.ValidateStruct(req); err != nil {
		r.logger.Error(err, "http - v1 - submitAttendance - ValidateStruct")
		return utils.WriteError(ctx, http.StatusBadRequest, "Validation failed", err)
	}

	// Parse date
	date, err := utils.ParseDateWIB(req.Date)
	if err != nil {
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid date format (expected YYYY-MM-DD)", nil)
	}

	// Check weekend
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return utils.WriteError(ctx, http.StatusBadRequest, "Cannot submit attendance on weekends", nil)
	}

	var (
		checkIn  *time.Time
		checkOut *time.Time
	)

	// Parse check-in time
	if req.Type == "check_in" {
		if req.CheckIn == "" {
			return utils.WriteError(ctx, http.StatusBadRequest, "check_in field is required", nil)
		}

		checkInTime, err := utils.ParseDateTimeWIB(req.Date + "T" + req.CheckIn)
		if err != nil {
			return utils.WriteError(ctx, http.StatusBadRequest, "Invalid check-in time format (expected HH:mm:ss)", nil)
		}

		checkIn = &checkInTime
	}

	// Parse check-out time
	if req.Type == "check_out" {
		if req.CheckOut == "" {
			return utils.WriteError(ctx, http.StatusBadRequest, "check_out field is required", nil)
		}

		checkOutTime, err := utils.ParseDateTimeWIB(req.Date + "T" + req.CheckOut)
		if err != nil {
			return utils.WriteError(ctx, http.StatusBadRequest, "Invalid check-in time format (expected HH:mm:ss)", nil)
		}

		checkOut = &checkOutTime
	}

	// Get user ID from context
	userIDRaw := ctx.Locals("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid user ID in context", nil)
	}

	req.UserID = userID
	req.ClientIP = ctx.IP()
	req.DateTime = date
	req.CheckInTime = checkIn
	req.CheckOutTime = checkOut

	// Call usecase
	err = r.attendance.SubmitAttendance(ctx.UserContext(), req)
	if err != nil {
		r.logger.Error(err, "http - v1 - submitAttendance - r.attendance.SubmitAttendance")

		if strings.Contains(err.Error(), "422") {
			return utils.WriteError(ctx, http.StatusUnprocessableEntity, "No active attendance period", nil)
		}
		if strings.Contains(err.Error(), "409") {
			return utils.WriteError(ctx, http.StatusConflict, "Attendance already submitted for this date", nil)
		}
		return utils.WriteError(ctx, http.StatusInternalServerError, "Failed to submit attendance", err)
	}

	return utils.WriteSuccess(ctx, "Attendance submitted successfully", nil)
}

// @Summary     Submit Overtime
// @Description Submit overtime request (between 0.5 and 3 hours) for current active period.
// @ID          submit-overtime
// @Tags        Employee
// @Accept      json
// @Produce     json
// @Param       payload body request.OvertimeRequest true "Overtime payload"
// @Success     201 {object} utils.BodySuccess
// @Failure     400 {object} utils.BodyFailure
// @Failure     409 {object} utils.BodyFailure
// @Failure     422 {object} utils.BodyFailure
// @Failure     500 {object} utils.BodyFailure
// @Security    JWTBearer
// @Router      /employee/overtimes [post]
func (r *V1) submitOvertime(ctx *fiber.Ctx) error {
	var req request.OvertimeRequest
	if err := ctx.BodyParser(&req); err != nil {
		r.logger.Error(err, "http - v1 - submitOvertime - BodyParser")
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := utils.ValidateStruct(req); err != nil {
		r.logger.Error(err, "http - v1 - submitOvertime - ValidateStruct")
		return utils.WriteError(ctx, http.StatusBadRequest, "Validation failed", err)
	}

	// Parse & validate date
	date, err := utils.ParseDateWIB(req.Date)
	if err != nil {
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD", nil)
	}

	// Validate hours
	if req.Hours < 0.5 || req.Hours > 3 {
		return utils.WriteError(ctx, http.StatusBadRequest, "Overtime hours must be between 0.5 and 3 hours", nil)
	}

	// Get user ID from context
	userIDRaw := ctx.Locals("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid user ID in context", nil)
	}

	// Inject tambahan data ke request
	req.UserID = userID
	req.ClientIP = ctx.IP()
	req.DateTime = date

	// Call usecase
	err = r.overtime.SubmitOvertime(ctx.UserContext(), req)
	if err != nil {
		r.logger.Error(err, "http - v1 - submitOvertime - r.overtime.SubmitOvertime")

		if strings.Contains(err.Error(), "422") {
			return utils.WriteError(ctx, http.StatusUnprocessableEntity, "No active attendance period", nil)
		}

		if strings.Contains(err.Error(), "409") {
			return utils.WriteError(ctx, http.StatusConflict, "Total overtime cannot exceed 3 hours per day", nil)
		}

		return utils.WriteError(ctx, http.StatusInternalServerError, "Failed to submit overtime", err)
	}

	return utils.WriteSuccess(ctx, "Overtime submitted successfully", nil)
}

// @Summary     Submit Reimbursement
// @Description Submit reimbursement request for the current active attendance period.
// @ID          submit-reimbursement
// @Tags        Employee
// @Accept      json
// @Produce     json
// @Param       payload body request.ReimbursementRequest true "Reimbursement payload"
// @Success     201 {object} utils.BodySuccess
// @Failure     400 {object} utils.BodyFailure
// @Failure     409 {object} utils.BodyFailure
// @Failure     422 {object} utils.BodyFailure
// @Failure     500 {object} utils.BodyFailure
// @Security    JWTBearer
// @Router      /employee/reimbursements [post]
func (r *V1) submitReimbursement(ctx *fiber.Ctx) error {
	var req request.ReimbursementRequest
	if err := ctx.BodyParser(&req); err != nil {
		r.logger.Error(err, "http - v1 - submitReimbursement - BodyParser")
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid request body", nil)
	}

	if err := utils.ValidateStruct(req); err != nil {
		r.logger.Error(err, "http - v1 - submitReimbursement - ValidateStruct")
		return utils.WriteError(ctx, http.StatusBadRequest, "Validation failed", err)
	}

	// Get user ID from context
	userIDRaw := ctx.Locals("user_id")
	userID, ok := userIDRaw.(uuid.UUID)
	if !ok {
		return utils.WriteError(ctx, http.StatusBadRequest, "Invalid user ID in context", nil)
	}

	// Inject tambahan data ke request
	req.UserID = userID
	req.ClientIP = ctx.IP()

	// Call usecase
	err := r.reimbursement.SubmitReimbursement(ctx.UserContext(), req)
	if err != nil {
		r.logger.Error(err, "http - v1 - submitReimbursement - r.reimbursement.SubmitReimbursement")

		if strings.Contains(err.Error(), "422") {
			return utils.WriteError(ctx, http.StatusUnprocessableEntity, "No active attendance period", nil)
		}
		return utils.WriteError(ctx, http.StatusInternalServerError, "Failed to submit reimbursement", err)
	}

	return utils.WriteSuccess(ctx, "Reimbursement submitted successfully", nil)
}
