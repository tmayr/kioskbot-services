package lib

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

func Notify() {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text:       "Hey, recibi un deposito de XXXX, Banco XXXX, de Wire McWireface \n ¿A quien se lo abonamos?",
		Color:      "#3AA3E3",
		CallbackID: "wire_user_selection?amount=xxxx",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:       "Users",
				Type:       "select",
				DataSource: "users",
			},
		},
	}

	params.Attachments = []slack.Attachment{attachment}
	_, _, err := api.PostMessage(
		os.Getenv("SLACK_PAYMENTS_NOTIFICATION_CHANNEL"),
		"¡Nuevo Deposito!",
		params,
	)

	if err != nil {
		log.Fatal(err)
		return
	}
}
