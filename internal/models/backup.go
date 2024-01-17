package models

// UserBackup holds components related to a user backup.
type UserBackup struct {
	DeleteSQL  string
	ImagesPath string
	InsertSQL  string
	Recipes    Recipes
	UserID     int64
}
