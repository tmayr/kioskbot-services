package lib

import (
	"fmt"
	KioskTypes "kioskbot-services/types"
	"log"
	"os"

	"github.com/nlopes/slack"
)

func Notify(wire KioskTypes.BankWire) {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text: fmt.Sprintf(
			"Hey, parece que recibi un deposito de *%v*, desde el *%v*, por *$%v*?",
			wire.Name,
			wire.Bank,
			wire.Amount,
		),
		MarkdownIn: []string{"text", "pretext"},
		Color:      "#3AA3E3",
		CallbackID: "wire_user_selection?amount=" + wire.Amount + "&apiKey=" + os.Getenv("KB_SERVICES_API_KEY"),
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
		"*¡Nuevo Deposito!*",
		params,
	)

	if err != nil {

		log.Fatal(err)
		return
	}
}
