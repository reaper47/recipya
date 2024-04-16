package models

import "strings"

// These constants enumerate all possible file types used by the software.
const (
	Crumb FileType = iota
	JSON
	PDF
	MXP
	Paprika
	TXT
	InvalidFileType
)

// FileType is an alias for a file type, e.g. JSON and PDF.
type FileType int64

// NewFileType creates a FileType from the file type name.
func NewFileType(fileType string) FileType {
	switch strings.ToLower(fileType) {
	case "crumb":
		return Crumb
	case "json":
		return JSON
	case "mxp":
		return MXP
	case "paprikarecipes":
		return Paprika
	case "pdf":
		return PDF
	case "txt":
		return TXT
	default:
		return InvalidFileType
	}
}

// Ext returns the FileType's extension.
func (f FileType) Ext() string {
	switch f {
	case Crumb:
		return ".crumb"
	case JSON:
		return ".json"
	case MXP:
		return ".mxp"
	case Paprika:
		return ".paprikarecipes"
	case PDF:
		return ".pdf"
	case TXT:
		return ".txt"
	default:
		return ""
	}
}
