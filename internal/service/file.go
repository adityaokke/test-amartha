package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/google/uuid"
)

type FileService interface {
	UploadFile(ctx context.Context, input entity.UploadFileInput) (fileURL string, err error)
}

func (s *fileService) UploadFile(ctx context.Context, input entity.UploadFileInput) (result string, err error) {
	if input.File == nil {
		err = errors.New("file is required")
		return
	}

	ext := strings.ToLower(filepath.Ext(input.File.Filename))
	if ext != ".pdf" && ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		err = errors.New("only pdf/png/jpg/jpeg allowed")
		return
	}

	// open upload
	src, err := input.File.Open()
	if err != nil {
		return "", fmt.Errorf("open: %w", err)
	}
	defer src.Close()

	dir := entity.LocalUploadPath

	// generate a server-side name (or keep original if you prefer)
	filename := uuid.New().String() + ext
	dstPath := filepath.Join(dir, filename)

	// write directly to final path (NO atomic temp)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("create: %w", err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		_ = dst.Close()
		return "", fmt.Errorf("write: %w", err)
	}
	if err := dst.Close(); err != nil {
		return "", fmt.Errorf("close: %w", err)
	}

	u, _ := url.Parse(os.Getenv("APP_HOST"))
	fileURL, err := url.JoinPath(u.String(), entity.PublicUploadPath, filename)
	if err != nil {
		return
	}
	return fileURL, nil
}

type fileService struct {
}

type InitiatorFile func(s *fileService) *fileService

func NewFileService() InitiatorFile {
	return func(s *fileService) *fileService {
		return s
	}
}

func (i InitiatorFile) Build() FileService {
	return i(&fileService{})
}
