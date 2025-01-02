package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ValidateSignedURLMiddleware is Gin middleware for validating signed URLs
func ValidateSignedURLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the full query string from the URL
		rawURL := c.Request.URL.String()

		// Validate the signed URL
		filePath, err := validateSignedURL(rawURL)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Store the validated file path in the context
		c.Set("filePath", filePath)

		// Continue to the next handler
		c.Next()
	}
}

// Validation logic (unchanged)
func validateSignedURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("Invalid URL")
	}

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	// Extract query parameters
	values := parsedURL.Query()
	filePath := values.Get("file")
	expires := values.Get("expires")
	signature := values.Get("signature")

	// Check if parameters are present
	if filePath == "" || expires == "" || signature == "" {
		return "", errors.New("Missing required query parameters")
	}

	// Check expiration time
	expirationTime, err := strconv.ParseInt(expires, 10, 64)
	if err != nil || time.Now().Unix() > expirationTime {
		return "", errors.New("URL has expired")
	}

	// Recreate the string to sign
	stringToSign := fmt.Sprintf("file=%s&expires=%s", filePath, expires)

	// Validate the signature
	mac := hmac.New(sha256.New, secretKey)
	_, err = mac.Write([]byte(stringToSign))
	if err != nil {
		return "", err
	}
	expectedSignature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	if signature != expectedSignature {
		return "", errors.New("Invalid signature")
	}

	return filePath, nil
}
