package response

import "github.com/RindangRamadhan/dealls-payslip-system/internal/entity"

type LoginResponse struct {
	Token string      `json:"token"`
	User  entity.User `json:"user"`
}
