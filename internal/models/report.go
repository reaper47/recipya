package models

import "time"

// ReportType represents
type ReportType int64

// ImportReportType is the ReportType for importing recipes, either from files or the web.
const ImportReportType ReportType = 1

// NewReport creates a new, initialized and empty Report of the given ReportType.
func NewReport(reportType ReportType) Report {
	return Report{
		CreatedAt: time.Now(),
		Logs:      make([]ReportLog, 0),
		Type:      reportType,
	}
}

// NewReportLog creates a new ReportLog from the title and error.
func NewReportLog(title string, isSuccess bool, err error, action string) ReportLog {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	return ReportLog{
		Error:     errStr,
		IsError:   err != nil,
		IsSuccess: isSuccess && err == nil,
		IsWarning: !isSuccess && err == nil,
		Title:     title,
		Action:    action,
	}
}

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
	ID        int64
	Title     string
	IsError   bool
	IsSuccess bool
	IsWarning bool
	Error     string
	Action    string
}
