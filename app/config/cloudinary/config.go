package cloudinary

import (
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	godotenv.Load()
	cldSecret := os.Getenv("API_SECRET")
	cldName := os.Getenv("API_NAME")
	cldKey := os.Getenv("API_KEY")

	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
	if err != nil {
		return nil, err
	}

	return cld, nil
}
