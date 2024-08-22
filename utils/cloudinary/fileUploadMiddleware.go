package cloudinary

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func FileUploadMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        file, err := c.FormFile("file")
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{
                "error": "Bad request",
            })
        }

        src, err := file.Open()
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "error": "Unable to open file",
            })
        }
        defer src.Close() // close file properly

        c.Set("filePath", file.Filename)
        c.Set("file", src)

        return next(c)
    }
}