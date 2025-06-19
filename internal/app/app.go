// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RindangRamadhan/dealls-payslip-system/config"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/controller/http/middleware"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/repo/persistent"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/attendance"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/overtime"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/payroll"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/payslip"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/reimbursement"
	"github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/user"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/httpserver"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/logger"
	"github.com/RindangRamadhan/dealls-payslip-system/pkg/postgres"

	_ "github.com/RindangRamadhan/dealls-payslip-system/docs" // Swagger docs.
	attendanceperiod "github.com/RindangRamadhan/dealls-payslip-system/internal/usecase/attendance_period"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	lgr := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, lgr, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		lgr.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Initialize Usecases
	userUsecase := user.NewUserUsecase(
		persistent.NewUserRepo(pg),
	)

	payrollUsecase := payroll.NewPayrollUsecase(
		*lgr,
		persistent.NewPayrollRepo(pg),
		persistent.NewUserRepo(pg),
		persistent.NewPayslipRepo(pg),
		persistent.NewOvertimeRepo(pg),
		persistent.NewAttendanceRepo(pg),
		persistent.NewReimbursementRepo(pg),
		persistent.NewAttendancePeriodRepo(pg),
	)

	payslipUsecase := payslip.NewPayslipUsecase(
		*lgr,
		persistent.NewPayslipRepo(pg),
		persistent.NewPayrollRepo(pg),
		persistent.NewOvertimeRepo(pg),
		persistent.NewAttendanceRepo(pg),
		persistent.NewReimbursementRepo(pg),
	)

	overtimeUsecase := overtime.NewOvertimeUsecase(
		*lgr,
		persistent.NewOvertimeRepo(pg),
		persistent.NewAttendancePeriodRepo(pg),
	)

	attendanceUsecase := attendance.NewAttendanceUsecase(
		*lgr,
		persistent.NewAttendanceRepo(pg),
		persistent.NewAttendancePeriodRepo(pg),
	)

	reimbursementUsecase := reimbursement.NewReimbursementUsecase(
		*lgr,
		persistent.NewReimbursementRepo(pg),
		persistent.NewAttendancePeriodRepo(pg),
	)

	attendancePeriodUsecase := attendanceperiod.NewAttendancePeriodUsecase(
		*lgr,
		persistent.NewAttendancePeriodRepo(pg),
	)

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))

	// Middleware
	httpServer.App.Use(middleware.CORS())
	httpServer.App.Use(middleware.Recovery(lgr))
	httpServer.App.Use(middleware.AuthMiddleware())
	httpServer.App.Use(middleware.AuditLog(persistent.NewAuditLogRepo(pg)))
	httpServer.App.Use(middleware.Logger(lgr))

	// Initialize HTTP Router
	http.NewRouter(
		httpServer.App, cfg, lgr,
		userUsecase,
		payrollUsecase,
		payslipUsecase,
		overtimeUsecase,
		attendanceUsecase,
		reimbursementUsecase,
		attendancePeriodUsecase,
	)

	// Start servers
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		lgr.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		lgr.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		lgr.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
