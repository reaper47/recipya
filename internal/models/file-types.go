package models

import "strings"

// These constants enumerate all possible file types used by the software.
const (
	JSON FileType = iota
	PDF
	MXP
	TXT
	InvalidFileType
)

// FileType is an alias for a file type, e.g. JSON and PDF.
type FileType int64

// NewFileType creates a FileType from the file type name.
func NewFileType(fileType string) FileType {
	switch strings.ToLower(fileType) {
	case "json":
		return JSON
	case "pdf":
		return PDF
	case "mxp":
		return MXP
	case "txt":
		return TXT
	default:
		return InvalidFileType
	}
}

// Ext returns the FileType's extension.
func (f FileType) Ext() string {
	switch f {
	case JSON:
		return ".json"
	case PDF:
		return ".pdf"
	case MXP:
		return ".mxp"
	case TXT:
		return ".txt"
	default:
		return ""
	}
}
