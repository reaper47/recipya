package models

import "time"

// ReportType represents
type ReportType int64

// ImportReportType is the ReportType for an import report.
const ImportReportType ReportType = 1

// Report holds information on a report.
type Report struct {
	CreatedAt time.Time
	ExecTime  time.Duration
	ID        int64
	Logs      []ReportLog
	Type      ReportType
}

// ReportLog holds information on a report's log.
type ReportLog struct {
	Error     string
	ID        int64
	IsSuccess bool
	Title     string
}
