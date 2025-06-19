package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/RindangRamadhan/dealls-payslip-system/internal/entity"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/repo"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AuditLog(r repo.AuditLogRepo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()
		c.Locals("request_id", requestID)

		// Capture request body (for POST/PUT)
		newValues := make(map[string]interface{})
		if c.Method() == fiber.MethodPost || c.Method() == fiber.MethodPut {
			if err := json.Unmarshal(c.Body(), &newValues); err != nil {
				newValues = map[string]interface{}{
					"unparsable_body": string(c.Body()),
				}
			}
		}

		newValuesJSON, _ := json.Marshal(newValues)

		method := c.Method()
		path := c.Path()
		ip := c.IP()

		err := c.Next()

		userIDRaw := c.Locals("user_id")
		userID, ok := userIDRaw.(uuid.UUID)
		if !ok {
			// skip if user not available (unauthenticated request)
			return err
		}

		recordID := uuid.Nil
		if rid, ok := c.Locals("record_id").(uuid.UUID); ok {
			recordID = rid
		}

		go func() {
			ctx := context.Background()

			auditLog := &entity.AuditLog{
				ID:        uuid.New(),
				UserID:    userID,
				Action:    method,
				TableName: extractTableFromPath(path),
				RecordID:  recordID,
				OldValues: json.RawMessage(`{}`),
				NewValues: newValuesJSON,
				IPAddress: ip,
				RequestID: requestID,
				CreatedAt: time.Now(),
			}

			if err := r.CreateAuditLog(ctx, auditLog); err != nil {
				fmt.Println("Error creating audit log:", err)
			}
		}()

		return err
	}
}

// extractTableFromPath tries to guess the table/entity from the request path.
func extractTableFromPath(path string) string {
	path = strings.ToLower(path)
	switch {
	case strings.Contains(path, "attendance"):
		return "attendances"
	case strings.Contains(path, "overtime"):
		return "overtimes"
	case strings.Contains(path, "reimbursement"):
		return "reimbursements"
	case strings.Contains(path, "payroll"):
		return "payrolls"
	default:
		return "unknown"
	}
}
