package fcm

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"regexp"
)

var Client *messaging.Client

func InitializeFirebase() {
	credentials, err := ioutil.ReadFile("./credentials.json")

	if err != nil {
		log.Println(err.Error())
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(credentials)}

	app, err := firebase.NewApp(context.Background(), nil, opts...)

	if err != nil {
		log.Println(err.Error())
	}

	Client, err = app.Messaging(context.Background())

	if err != nil {
		panic(err)
	}
}

func SendNotification(title string, message string, topic string, data map[string]string) {
	if _, err := Client.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
		Topic: cleanDeviceId(topic),
		Data:  data,
	}); err != nil {
		log.Println(err.Error())
	}
}

func RegisterTopic(email string, deviceId string) {
	if _, err := Client.SubscribeToTopic(context.Background(), []string{deviceId}, cleanDeviceId(email)); err != nil {
		log.Println(err.Error())
	}
}

func cleanDeviceId(emailAddress string) string {
	re := regexp.MustCompile(`\W`)
	return re.ReplaceAllString(emailAddress, "")
}
