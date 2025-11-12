package assets_firebase

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FirebaseMessaging struct {
	client *messaging.Client
}

func NewFirebase(ctx context.Context, credentialsFile string) (*FirebaseMessaging, error) {
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseMessaging{client: client}, nil

}

func (f *FirebaseMessaging) SendToToken(ctx context.Context, token string, notification *messaging.Notification) error {
	message := &messaging.Message{
		Token:        token,
		Notification: notification,
	}
	_, err := f.client.Send(ctx, message)
	if err != nil {
		return err
	}
	// log.Printf("Successfully sent message to token: %s\n", response)
	return nil
}

func (f *FirebaseMessaging) SendToTopic(ctx context.Context, topic string, notification *messaging.Notification) error {
	message := &messaging.Message{
		Topic:        topic,
		Notification: notification,
	}
	_, err := f.client.Send(ctx, message)
	if err != nil {
		return err
	}
	// log.Printf("Successfully sent message to topic: %s\n", response)
	return nil
}
