CREATE TABLE payrolls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    period_id UUID NOT NULL REFERENCES attendance_periods(id),
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    total_amount DECIMAL(15,2) NOT NULL DEFAULT 0,
    employee_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    UNIQUE(period_id)
);

CREATE INDEX idx_payrolls_period ON payrolls(period_id);
CREATE INDEX idx_payrolls_processed_at ON payrolls(processed_at);