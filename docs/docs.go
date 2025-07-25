// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/attendance-periods": {
            "post": {
                "security": [
                    {
                        "JWTBearer": []
                    }
                ],
                "description": "Create a new attendance period by providing start and end dates.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Create Attendance Period",
                "operationId": "create-attendance-period",
                "parameters": [
                    {
                        "description": "Attendance Period Payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AttendancePeriodRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.BodySuccess"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/entity.AttendancePeriod"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    }
                }
            }
        },
        "/admin/payrolls": {
            "post": {
                "security": [
                    {
                        "JWTBearer": []
                    }
                ],
                "description": "Generate payroll for a specific attendance period. Requires authentication.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Process Payroll",
                "operationId": "process-payroll",
                "parameters": [
                    {
                        "description": "Payroll request payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.PayrollRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.BodySuccess"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/entity.Payroll"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    }
                }
            }
        },
        "/admin/payrolls/summary": {
            "get": {
                "security": [
                    {
                        "JWTBearer": []
                    }
                ],
                "description": "Get payroll summary for a given period (via query param)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Get Payroll Summary",
                "operationId": "get-payroll-summary",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Attendance Period ID (UUID)",
                        "name": "period_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.BodySuccess"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.PayrollSummaryResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Authenticate User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User Login",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "Login Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.BodySuccess"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.LoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    }
                }
            }
        },
        "/employee/attendances": {
            "post": {
                "security": [
                    {
                        "JWTBearer": []
                    }
                ],
                "description": "Submit attendance (check-in and check-out) for current active period.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Employee"
                ],
                "summary": "Submit Attendance",
                "operationId": "submit-attendance",
                "parameters": [
                    {
                        "description": "Attendance payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.AttendanceRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/utils.BodySuccess"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    }
                }
            }
        },
        "/employee/overtimes": {
            "post": {
                "security": [
                    {
                        "JWTBearer": []
                    }
                ],
                "description": "Submit overtime request (between 0.5 and 3 hours) for current active period.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Employee"
                ],
                "summary": "Submit Overtime",
                "operationId": "submit-overtime",
                "parameters": [
                    {
                        "description": "Overtime payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.OvertimeRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/utils.BodySuccess"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    }
                }
            }
        },
        "/employee/payslips": {
            "get": {
                "security": [
                    {
                        "JWTBearer": []
                    }
                ],
                "description": "Retrieve payslip for the authenticated employee.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Employee"
                ],
                "summary": "Get Payslip",
                "operationId": "get-payslip",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Attendance Period ID (UUID)",
                        "name": "period_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "payslip retrieved successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.BodySuccess"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.PayslipResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "invalid input or bad request",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "404": {
                        "description": "payslip not found",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "500": {
                        "description": "internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    }
                }
            }
        },
        "/employee/reimbursements": {
            "post": {
                "security": [
                    {
                        "JWTBearer": []
                    }
                ],
                "description": "Submit reimbursement request for the current active attendance period.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Employee"
                ],
                "summary": "Submit Reimbursement",
                "operationId": "submit-reimbursement",
                "parameters": [
                    {
                        "description": "Reimbursement payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ReimbursementRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/utils.BodySuccess"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.BodyFailure"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.Attendance": {
            "type": "object",
            "properties": {
                "check_in": {
                    "type": "string"
                },
                "check_out": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "created_by": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "ip_address": {
                    "type": "string"
                },
                "period_id": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "entity.AttendancePeriod": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "created_by": {
                    "type": "string"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_active": {
                    "type": "boolean"
                },
                "start_date": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by": {
                    "type": "string"
                }
            }
        },
        "entity.Overtime": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "created_by": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "hours": {
                    "type": "number"
                },
                "id": {
                    "type": "string"
                },
                "ip_address": {
                    "type": "string"
                },
                "period_id": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "entity.Payroll": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "created_by": {
                    "type": "string"
                },
                "employee_count": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "period_id": {
                    "type": "string"
                },
                "processed_at": {
                    "type": "string"
                },
                "total_amount": {
                    "type": "number"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by": {
                    "type": "string"
                }
            }
        },
        "entity.Payslip": {
            "type": "object",
            "properties": {
                "attendance_days": {
                    "type": "integer"
                },
                "attendance_salary": {
                    "type": "number"
                },
                "base_salary": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "created_by": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "overtime_hours": {
                    "type": "number"
                },
                "overtime_salary": {
                    "type": "number"
                },
                "payroll_id": {
                    "type": "string"
                },
                "period_id": {
                    "type": "string"
                },
                "reimbursement_total": {
                    "type": "number"
                },
                "total_pay": {
                    "type": "number"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                },
                "working_days": {
                    "type": "integer"
                }
            }
        },
        "entity.Reimbursement": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "created_at": {
                    "type": "string"
                },
                "created_by": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "ip_address": {
                    "type": "string"
                },
                "period_id": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "entity.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "created_by": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_admin": {
                    "type": "boolean"
                },
                "salary": {
                    "type": "number"
                },
                "updated_at": {
                    "type": "string"
                },
                "updated_by": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "request.AttendancePeriodRequest": {
            "type": "object",
            "required": [
                "end_date",
                "start_date"
            ],
            "properties": {
                "end_date": {
                    "type": "string",
                    "example": "2025-01-31"
                },
                "start_date": {
                    "type": "string",
                    "example": "2025-01-01"
                }
            }
        },
        "request.AttendanceRequest": {
            "type": "object",
            "required": [
                "date",
                "type"
            ],
            "properties": {
                "check_in": {
                    "type": "string",
                    "example": "08:00:00"
                },
                "check_out": {
                    "type": "string",
                    "example": "17:00:00"
                },
                "date": {
                    "type": "string",
                    "example": "2025-06-20"
                },
                "type": {
                    "type": "string",
                    "enum": [
                        "check_in",
                        "check_out"
                    ],
                    "example": "check_in"
                }
            }
        },
        "request.LoginRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "example": "admin123"
                },
                "username": {
                    "type": "string",
                    "example": "admin"
                }
            }
        },
        "request.OvertimeRequest": {
            "type": "object",
            "required": [
                "date",
                "hours"
            ],
            "properties": {
                "date": {
                    "type": "string",
                    "example": "2025-06-20"
                },
                "hours": {
                    "type": "number",
                    "maximum": 3,
                    "minimum": 0.5,
                    "example": 1
                }
            }
        },
        "request.PayrollRequest": {
            "type": "object",
            "required": [
                "period_id"
            ],
            "properties": {
                "period_id": {
                    "type": "string",
                    "example": "1e6c4313-3148-4fc8-a6cd-b6cd4674e041"
                }
            }
        },
        "request.ReimbursementRequest": {
            "type": "object",
            "required": [
                "amount",
                "description"
            ],
            "properties": {
                "amount": {
                    "type": "number",
                    "minimum": 0
                },
                "description": {
                    "type": "string"
                }
            }
        },
        "response.EmployeePay": {
            "type": "object",
            "properties": {
                "total_pay": {
                    "type": "number"
                },
                "user_id": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "response.LoginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/entity.User"
                }
            }
        },
        "response.PayrollSummaryResponse": {
            "type": "object",
            "properties": {
                "employee_pays": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.EmployeePay"
                    }
                },
                "payroll": {
                    "$ref": "#/definitions/entity.Payroll"
                },
                "total_amount": {
                    "type": "number"
                }
            }
        },
        "response.PayslipResponse": {
            "type": "object",
            "properties": {
                "attendances": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Attendance"
                    }
                },
                "overtimes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Overtime"
                    }
                },
                "payslip": {
                    "$ref": "#/definitions/entity.Payslip"
                },
                "reimbursements": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Reimbursement"
                    }
                }
            }
        },
        "utils.BodyFailure": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "utils.BodySuccess": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        }
    },
    "securityDefinitions": {
        "JWTBearer": {
            "description": "Provide your JWT token prefixed with \"Bearer\", e.g., \"Bearer \u003cyour_token\u003e\"",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Dealls Payslip System API",
	Description:      "API for the Dealls Payslip System, featuring attendance tracking, overtime, reimbursement, and payroll summary.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
