package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	UserID    uuid.UUID       `json:"user_id" db:"user_id"`
	Action    string          `json:"action" db:"action"`
	TableName string          `json:"table_name" db:"table_name"`
	RecordID  uuid.UUID       `json:"record_id" db:"record_id"`
	OldValues json.RawMessage `json:"old_values" db:"old_values"`
	NewValues json.RawMessage `json:"new_values" db:"new_values"`
	IPAddress string          `json:"ip_address" db:"ip_address"`
	RequestID string          `json:"request_id" db:"request_id"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
