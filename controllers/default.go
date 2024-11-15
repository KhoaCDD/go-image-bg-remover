package controllers

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"os"
	"os/exec"

	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"

	"github.com/beego/beego/v2/server/web"
)

const MaxUploadSize = 10 * 1024 * 1024 // 10 MB
const UploadPath = "./uploads"         // Directory where files will be saved

type MainController struct {
	web.Controller
}

func (c *MainController) Get() {
	c.Data["Title"] = "Image Background Remover"
	c.Data["Content"] = "example@example.com"
	c.Data["User"] = "Khoa CDD"
	c.TplName = "index.tpl"
}

func (c *MainController) UploadImage() {
	fmt.Println("UploadImage")
	var errUpload string
	defer func() {
		c.Data["Error"] = errUpload
		c.TplName = "index.tpl"
	}()
	if err := c.Ctx.Request.ParseMultipartForm(MaxUploadSize); err != nil {
		errUpload = ("File too big")
	}
	file, handler, err := c.Ctx.Request.FormFile("file")
	if err != nil {
		errUpload = ("Error get file")
	}
	defer file.Close()
	fmt.Println(handler.Filename)

	// dst, err := os.Create(filepath.Join(UploadPath, handler.Filename))
	// if err != nil {
	// 	errUpload = ("Unable to create the file for saving: " + err.Error())
	// }
	// defer dst.Close()

	// // Copy the uploaded file to the destination file
	// if _, err = io.Copy(dst, file); err != nil {
	// 	errUpload = ("Unable to save the file")
	// }

	// Decode the image

	// img, format, err := image.Decode(file)
	// if err != nil {
	// 	errUpload = ("Unable to Decode the file for saving: " + err.Error())
	// 	return
	// }

	// width := img.Bounds().Dx()
	// height := img.Bounds().Dy()
	// fmt.Println("Image width: ", width)
	// fmt.Println("Image height: ", height)
	// Convert image to bytes using a buffer
	// imageBytes, err := io.ReadAll(file)
	// if err != nil {
	// 	errUpload = ("Unable to convert file to byte: " + err.Error())
	// 	return
	// }
	// r := bytes.NewReader(imageBytes)
	img, err := jpeg.Decode(file)
	if err != nil {
		errUpload = ("Unable to create the file for saving1: " + err.Error())
		return
	}

	rgb24Stream, err := imageToRGB24(img)
	if err != nil {
		errUpload = ("Unable to create the file for saving3: " + err.Error())
		return
	}

	fmt.Println("\n\n\n\nImage:", len(rgb24Stream))

	processedImage, err := ProcessImageWithRembg(rgb24Stream, fmt.Sprintf("%d", "1892"), fmt.Sprintf("%d", "1419"))
	if err != nil {
		errUpload = ("Unable to create the file for saving: " + err.Error())
	}
	fmt.Println("Processed image \n", len(processedImage))

	// err = saveImageToFolder(processedImage, filepath.Join(UploadPath, handler.Filename))
	// if err != nil {
	// 	errUpload = ("Unable to save the file" + err.Error())
	// }

	// File saved successfully
	fmt.Println("File uploaded successfully")

}

func ProcessImageWithRembg(imageBytes []byte, height string, width string) ([]byte, error) {
	cmd := exec.Command("rembg", "b", "1419", "1892")

	// imageFile, errRead := os.Open("./uploads/TUS-175.jpeg")
	// if errRead != nil {
	// 	fmt.Println("Image processing completed successfully0.", errRead.Error())
	// }
	// defer imageFile.Close()
	reader := bytes.NewReader(imageBytes)
	cmd.Stdin = reader

	fs := http.FileServer(http.Dir("/path/to/images"))
	http.Handle("/", fs)

	// Create a pipe to pass image bytes to the Python process
	// stdin, err := cmd.StdinPipe()
	// if err != nil {
	// 	fmt.Println("Image processing completed successfully1.", err.Error())
	// 	return nil, err
	// }

	// Create a buffer to capture the output

	// // Start the process
	// if err := cmd.Start(); err != nil {
	// 	fmt.Println("Image processing completed successfully.2", err.Error())
	// 	return nil, err
	// }

	// Write the image bytes to stdin
	// go func() {
	// 	defer stdin.Close() // Ensure stdin is closed after writing
	// 	if size, err := stdin.Write(imageBytes); err != nil {
	// 		fmt.Println("Image processing completed successfully3.", size, err.Error())
	// 	}
	// }()

	// Wait for the process to finish
	logs, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Image processing completed successfully.4", err.Error())
		fmt.Println(string(logs))
		return nil, err
	}

	fmt.Println("Image processing completed.", err.Error())

	// Return the processed image bytes
	return nil, nil
}

func saveImageToFolder(imgBytes []byte, path string) error {
	// Create or open the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the bytes to the file
	_, err = file.Write(imgBytes)
	if err != nil {
		return err
	}

	return nil
}

func imageToRGB24(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	rgbStream := make([]byte, 0, width*height*3) // RGB24 = 3 bytes per pixel

	// Iterate over each pixel and extract RGB values
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the color of the pixel
			r, g, b, _ := img.At(x, y).RGBA()

			// Convert the pixel values from 16-bit (0-65535) to 8-bit (0-255)
			rgbStream = append(rgbStream, byte(r>>8), byte(g>>8), byte(b>>8))
		}
	}

	return rgbStream, nil
}
