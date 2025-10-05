package entity

import "mime/multipart"

const (
	LocalUploadPath  = "storage/uploads"
	PublicUploadPath = "public/uploads"
)

type UploadFileInput struct {
	File *multipart.FileHeader
}
