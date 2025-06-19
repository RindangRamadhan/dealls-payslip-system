package middleware

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func buildRequestMessage(ctx *fiber.Ctx) string {
	var result strings.Builder

	traceID, ok := ctx.Locals("request_id").(string)
	if !ok || traceID == "" {
		traceID = uuid.New().String()
	}

	result.WriteString(fmt.Sprintf("Trace-ID=%s", traceID))
	result.WriteString(" - ")
	result.WriteString(ctx.IP())
	result.WriteString(" - ")
	result.WriteString(ctx.Method())
	result.WriteString(" ")
	result.WriteString(ctx.OriginalURL())
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(ctx.Response().StatusCode()))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(len(ctx.Response().Body())))

	return result.String()
}

func Logger(l logger.Interface) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()

		l.Info(buildRequestMessage(ctx))

		return err
	}
}
