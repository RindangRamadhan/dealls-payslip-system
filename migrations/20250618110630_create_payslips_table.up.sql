CREATE TABLE payslips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    payroll_id UUID NOT NULL REFERENCES payrolls(id),
    period_id UUID NOT NULL REFERENCES attendance_periods(id),
    base_salary DECIMAL(12,2) NOT NULL,
    attendance_days INTEGER NOT NULL DEFAULT 0,
    working_days INTEGER NOT NULL DEFAULT 0,
    attendance_salary DECIMAL(12,2) NOT NULL DEFAULT 0,
    overtime_hours DECIMAL(6,2) NOT NULL DEFAULT 0,
    overtime_salary DECIMAL(12,2) NOT NULL DEFAULT 0,
    reimbursement_total DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_pay DECIMAL(12,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    UNIQUE(user_id, payroll_id)
);

CREATE INDEX idx_payslips_user ON payslips(user_id);
CREATE INDEX idx_payslips_payroll ON payslips(payroll_id);
CREATE INDEX idx_payslips_period ON payslips(period_id);