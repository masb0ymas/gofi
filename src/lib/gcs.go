package lib

import (
	"context"
	"encoding/pem"
	"errors"
	"fmt"
	"gofi/src/config"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/masb0ymas/go-utils/pkg"
	"google.golang.org/api/option"
)

var client *storage.Client

type GoogleCloudStorage struct {
	Bucket         string
	ProjectID      string
	Region         string
	GoogleAccessID string
	ExpiresIn      string
	FilePath       string
}

func NewGoogleCloudStorageConfig() *GoogleCloudStorage {
	return &GoogleCloudStorage{
		Bucket:         config.Env("STORAGE_BUCKET", "gofi"),
		ProjectID:      config.Env("STORAGE_PROJECT_ID", "your-project-id"),
		Region:         config.Env("STORAGE_REGION", "asia-southeast2"),
		GoogleAccessID: config.Env("STORAGE_ACCESS_ID", "your-access-id"),
		ExpiresIn:      config.Env("STORAGE_EXPIRES_IN", "7"),
		FilePath:       config.Env("STORAGE_FILE_PATH", "./secret/service-account-gcs.json"),
	}
}

// Initial Google Cloud Storage
func InitGCS() {
	var err error
	ctx := context.Background()
	config := NewGoogleCloudStorageConfig()

	// Create new client
	client, err = storage.NewClient(ctx, option.WithCredentialsFile(config.FilePath))
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()

	// client to bucket
	bucket := client.Bucket(config.Bucket)

	// Check if bucket exists
	_, err = bucket.Attrs(ctx)
	if err != nil {
		// Create new bucket if it doesn't exist
		if err == storage.ErrBucketNotExist {
			// Create bucket
			if err := bucket.Create(ctx, config.ProjectID, &storage.BucketAttrs{Location: config.Region}); err != nil {
				panic(fmt.Sprintf("Failed to create bucket: %v", err))
			}
			// Print success message
			msg := pkg.Println("Google Cloud Storage", fmt.Sprintf("Bucket %s created successfully", config.Bucket))
			log.Println(msg)
		} else {
			// Check if bucket exists
			panic(fmt.Sprintf("Failed to check bucket: %v", err))
		}
	} else {
		// Print success message
		msg := pkg.Println("Google Cloud Storage", fmt.Sprintf("Connected to existing bucket: %s", config.Bucket))
		log.Println(msg)
	}
}

// Upload File to Cloud Storage
func UploadFile(filename string, file_data *multipart.File) (string, time.Time, error) {
	config := NewGoogleCloudStorageConfig()
	ctx := context.Background()

	// Read file into byte slice
	fileBytes, err := io.ReadAll(*file_data)
	if err != nil {
		return "", time.Time{}, err
	}

	object := client.Bucket(config.Bucket).Object(filename)
	wc := object.NewWriter(ctx)

	// Copy the file data to GCS
	if _, err := wc.Write(fileBytes); err != nil {
		return "", time.Time{}, err
	}

	// Close the writer
	if err := wc.Close(); err != nil {
		return "", time.Time{}, err
	}

	// Generate public URL for the uploaded file
	signed_url, expires_at, err := GenerateSignedURL(filename)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate public URL: %v", err)
	}

	return signed_url, expires_at, nil
}

// Upload File to Cloud Storage
func UploadFilePublicView(filename string, file_data *multipart.File) (string, time.Time, error) {
	config := NewGoogleCloudStorageConfig()
	ctx := context.Background()

	// Read file into byte slice
	fileBytes, err := io.ReadAll(*file_data)
	if err != nil {
		return "", time.Time{}, err
	}

	object := client.Bucket(config.Bucket).Object(filename)
	wc := object.NewWriter(ctx)

	// Copy the file data to GCS
	if _, err := wc.Write(fileBytes); err != nil {
		return "", time.Time{}, err
	}

	// Close the writer
	if err := wc.Close(); err != nil {
		return "", time.Time{}, err
	}

	err = object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		return "", time.Time{}, err
	}

	// Generate public URL for the uploaded file
	public_url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", config.Bucket, filename)

	return public_url, time.Time{}, nil
}

// Upload File From Local to Cloud Storage
func UploadFileFromLocal(filename string, file_data *os.File) (string, time.Time, error) {
	config := NewGoogleCloudStorageConfig()
	ctx := context.Background()

	// Read file into byte slice
	fileBytes, err := io.ReadAll(file_data)
	if err != nil {
		return "", time.Time{}, err
	}

	object := client.Bucket(config.Bucket).Object(filename)
	wc := object.NewWriter(ctx)

	// Copy the file data to GCS
	if _, err := wc.Write(fileBytes); err != nil {
		return "", time.Time{}, err
	}

	// Close the writer
	if err := wc.Close(); err != nil {
		return "", time.Time{}, err
	}

	// Generate public URL for the uploaded file
	signed_url, expires_at, err := GenerateSignedURL(filename)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate public URL: %v", err)
	}

	return signed_url, expires_at, nil
}

// Generate Signed URL from Cloud Storage
func GenerateSignedURL(keyFile string) (string, time.Time, error) {
	config := NewGoogleCloudStorageConfig()

	// Set the expiration time for the signed URL (e.g., 7 days from now)
	expires_at := time.Now().Add(time.Duration(pkg.StringToInt32(config.ExpiresIn)) * 24 * time.Hour)
	privateKey, err := getPrivateKeyFromFile("secret/pk-gcs.pem")
	if err != nil {
		return "", time.Time{}, err
	}

	opts := &storage.SignedURLOptions{
		Method:         http.MethodGet,
		Expires:        expires_at,
		GoogleAccessID: config.GoogleAccessID,
		PrivateKey:     privateKey,
	}

	// Generate a signed URL for the object
	signed_url, err := client.Bucket(config.Bucket).SignedURL(keyFile, opts)
	if err != nil {
		return "", time.Time{}, err
	}

	// Print the signed URL
	fmt.Println("Signed URL:", signed_url)

	return signed_url, expires_at, nil
}

// Get Google Cloud Storage Client
func GetGCSClient() *storage.Client {
	return client
}

// getPrivateKeyFromFile reads a private key from a PEM file
func getPrivateKeyFromFile(keyPath string) ([]byte, error) {
	// Read the private key file
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	// Decode the PEM block
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	// Return the original PEM bytes
	return keyBytes, nil
}
