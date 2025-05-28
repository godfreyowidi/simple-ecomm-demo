package pkg

import (
	"fmt"
	"log"
	"os"

	africastalking "github.com/tech-kenya/africastalkingsms"
)

type SMSService struct {
	client *africastalking.SMSClient
}

func NewSMSService() (*SMSService, error) {
	apiKey := os.Getenv("AT_API_KEY")
	username := os.Getenv("AT_USERNAME")
	shortCode := os.Getenv("AT_SHORTCODE")
	sandbox := os.Getenv("AT_SANDBOX")

	client, err := africastalking.NewSMSClient(apiKey, username, shortCode, sandbox)
	if err != nil {
		return nil, err
	}

	return &SMSService{client: client}, nil
}

func (s *SMSService) SendOrderConfirmationSMS(toPhone, customerName string) error {
	message := fmt.Sprintf("Hi %s, your order has been received and is being processed. Thank you!", customerName)

	resp, err := s.client.SendSMS(toPhone, message)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	for _, recipient := range resp.SMSMessageData.Recipients {
		if recipient.StatusCode != 101 && recipient.StatusCode != 200 {
			return fmt.Errorf("SMS not delivered to %s: status=%s", recipient.Number, recipient.Status)
		}
	}

	log.Printf("SMS sent successfully to %s", toPhone)
	return nil
}
