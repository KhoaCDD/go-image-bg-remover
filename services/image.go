package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func ProcessedImage(fileName string, clientID string) error {
	// Construct the command to remove the background

	UploadPath := os.Getenv("UPLOAD_DIR")
	ProcessPath := os.Getenv("PROCESSED_DIR")
	uploadDst := filepath.Join(UploadPath, fileName)
	processedDst := filepath.Join(ProcessPath, fileName)

	cmd := exec.Command("rembg", "i", uploadDst, processedDst)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("failed to process image: %v, output: %s", err, string(output))
		return err
	}

	// Log the sessionId for tracking
	fmt.Printf("Image processed for clientID: %s\n", clientID)

	return nil
}

// generateRandomString generates a random string of the specified length
func GenerateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	for i := range bytes {
		bytes[i] = letters[bytes[i]%byte(len(letters))]
	}
	return string(bytes)
}

// GenerateSignedURL generates a signed URL with an expiration timestamp
func GenerateSignedURL(filePath string, expiration time.Time) (string, error) {
	baseURL := os.Getenv("BASE_URL")
	secretKey := []byte(os.Getenv("SECRET_KEY"))

	// Create a query string with the file path and expiration time
	values := url.Values{}
	values.Set("file", filePath)
	values.Set("expires", fmt.Sprintf("%d", expiration.Unix()))

	// Create a string to sign
	stringToSign := values.Encode()

	// Generate the HMAC signature
	mac := hmac.New(sha256.New, secretKey)
	_, err := mac.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	// Append the signature to the query parameters
	values.Set("signature", signature)

	// Construct the signed URL
	signedURL := fmt.Sprintf("%s?%s", baseURL, values.Encode())
	return signedURL, nil
}
