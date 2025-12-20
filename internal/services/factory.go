package services

type Services struct {
	Email  EmailService
	Google GoogleService
	S3     S3Service
}
