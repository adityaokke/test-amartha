package rest

import (
	"fmt"
	"net/http"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/service"
	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(
	fileService service.FileService,
) FileHandler {
	return FileHandler{
		fileService: fileService,
	}
}

func (d FileHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.String(http.StatusBadRequest, "missing file")
	}

	fileName, err := d.fileService.UploadFile(c.Request().Context(), entity.UploadFileInput{
		File: file,
	})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	hostName := c.Request().Host
	url := fmt.Sprintf("http://%s/%s/%s", hostName, entity.PublicUploadPath, fileName)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"file": map[string]interface{}{
				"url": url,
			},
		},
	})
}
