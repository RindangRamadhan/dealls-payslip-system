package v1

import (
	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/request"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/v1/response"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// @Summary     User Login
// @Description Authenticate User
// @ID          login
// @Tags  	    Auth
// @Accept      json
// @Produce     json
// @Param       request body request.LoginRequest true "Login Request"
// @Success     200 {object} utils.BodySuccess{data=response.LoginResponse}
// @Failure     400 {object} utils.BodyFailure
// @Failure     401 {object} utils.BodyFailure
// @Failure     500 {object} utils.BodyFailure
// @Router      /auth/login [post]
func (r *V1) login(ctx *fiber.Ctx) error {
	var req request.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := utils.ValidateStruct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := r.user.GetUserByUsername(ctx.UserContext(), req.Username)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := utils.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	resp := response.LoginResponse{
		Token: token,
		User:  *user,
	}

	return utils.WriteSuccess(ctx, "successfully login", resp)
}
