package fcm

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"regexp"

	"io/ioutil"
	"log"
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
	response, err := Client.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
		Topic: topic,
		Data:  data,
	})

	if err != nil {
		log.Println(err.Error())
	}
	log.Println(response)
}

func RegisterTopic(email string, deviceId string) {
	response, err := Client.SubscribeToTopic(context.Background(), []string{deviceId}, cleanDeviceId(email))
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(response)
}

func cleanDeviceId(emailAddress string) string {
	re := regexp.MustCompile(`\W`)
	return re.ReplaceAllString(emailAddress, "")
}

func SendMultiNotification(title string, message string, devices []string) {
	response, err := Client.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
		Tokens: devices,
	})

	if err != nil {
		log.Println(err.Error())
	}
	log.Println(response)
}
