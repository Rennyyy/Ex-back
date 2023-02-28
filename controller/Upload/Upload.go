package Upload

import (
	"context"
	"crypto/rand"
	"fmt"
	"mime/multipart"

	"example.com/m/v2/keys"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2/bson"
)

func UploadPv(rn *echo.Group) {

	rn.POST("/upload", UploadImage)
	rn.POST("/destroy", DeleteImage)

}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func UploadPhoto(x *multipart.FileHeader) (string, error) {
	// fmt.Println("Upload X ==> ", x)
	ctx := context.Background()
	image := x
	cld, err3 := cloudinary.NewFromURL(keys.CLOUDINARY_URL)
	if err3 != nil {
		return "", err3
	}

	ImageName, err := GenerateRandomString(32)
	if err != nil {
		return "", err
	}

	src, err := image.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	resp, _ := cld.Upload.Upload(ctx, src, uploader.UploadParams{PublicID: ImageName, Folder: "products"})

	// fmt.Println("resp ==> ", resp.SecureURL)
	ss := resp.SecureURL

	return ss, nil

}

func UploadImage(_pr echo.Context) error {

	ctx := context.Background()
	cld, err3 := cloudinary.NewFromURL(keys.CLOUDINARY_URL)
	if err3 != nil {
		return _pr.JSON(500, echo.Map{"error": err3.Error(), "message": "Error_cloudinary_server"})
	}

	image, err := _pr.FormFile("image")
	if err != nil {
		return _pr.JSON(400, echo.Map{"error": err.Error(), "message": "Error_Image"})
	}

	ImageName, err := GenerateRandomString(32)
	if err != nil {
		return _pr.JSON(500, echo.Map{"error": err.Error(), "message": "Error_GenerateRandomString"})
	}

	src, err := image.Open()
	if err != nil {
		return _pr.JSON(500, echo.Map{"error": err.Error(), "message": "Error_Open_Image"})
	}
	defer src.Close()

	resp, _ := cld.Upload.Upload(ctx, src, uploader.UploadParams{PublicID: ImageName})

	// fmt.Println("resp ==> ", resp.SecureURL)
	ss := bson.M{
		"PublicID":  ImageName,
		"SecureURL": resp.SecureURL,
	}

	return _pr.JSON(200, echo.Map{
		"status":  true,
		"message": "Add Product Completed",
		"result":  ss,
	})
}

func DeleteImage(_pr echo.Context) error {

	ctx := context.Background()
	cld, err3 := cloudinary.NewFromURL(keys.CLOUDINARY_URL)
	if err3 != nil {
		return _pr.JSON(500, echo.Map{"error": err3.Error(), "message": "Error_cloudinary_server"})
	}

	_PublicID := _pr.FormValue("public_id")
	fmt.Println("_PublicID ==> ", _PublicID)

	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: _PublicID})
	if err != nil {
		return _pr.JSON(500, echo.Map{"error": err.Error(), "message": "Error_Destroy"})
	}

	return _pr.JSON(200, echo.Map{
		"status":  true,
		"message": "Deleted Image",
		// "result":  ss,
	})

}
