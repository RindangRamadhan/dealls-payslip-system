CREATE TABLE attendance_periods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    CONSTRAINT check_dates CHECK (end_date >= start_date)
);

CREATE INDEX idx_attendance_periods_dates ON attendance_periods(start_date, end_date);
CREATE INDEX idx_attendance_periods_active ON attendance_periods(is_active);