package lib

import (
	KioskbotTypes "kioskbot-services/types"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	imapreader "github.com/erizocosmico/go-imapreader"
)

func getAmount(matchExpression string, body string) string {
	r, _ := regexp.Compile(matchExpression)
	amountMatched := r.FindStringSubmatch(body)
	if amountMatched != nil {
		return strings.Replace(amountMatched[1], ".", "", -1)
	}

	return ""
}

func getName(matchExpression string, body string) string {
	r, _ := regexp.Compile(matchExpression)
	whoWired := r.FindStringSubmatch(body)
	if whoWired != nil {
		return strings.Title(strings.ToLower(whoWired[1]))
	}

	return ""
}

func matchScotiabank(body string) KioskbotTypes.BankWire {
	matches := KioskbotTypes.BankWire{
		Amount: getAmount(`por un monto de \$([\d\.]+)`, body),
		Name:   getName(`Sr\(a\) (.+) ha instruido`, body),
		Bank:   "Scotiabank",
	}

	return matches
}

func matchBancoBCI(body string) KioskbotTypes.BankWire {
	matches := KioskbotTypes.BankWire{
		Amount: getAmount(`Monto transferido[\:\s]+\$([\d\.]+)`, body),
		Name:   getName(`Titular de la cuenta de origen[\:\s]+(.+)[\s\n]+`, body),
		Bank:   "Banco BCI",
	}

	return matches
}

func matchBancoSantander(body string) KioskbotTypes.BankWire {
	matches := KioskbotTypes.BankWire{
		Amount: getAmount(`Monto[\s\n]+de la operación[\s\n]+[\:\$\s]+([\d\.]+)`, body),
		Name:   getName(`Le informamos que hoy, .+, nuestro\(a\) cliente[\s\n]+(.+) ha instruído`, body),
		Bank:   "Banco Santander",
	}

	return matches
}

func matchBancoChile(body string) KioskbotTypes.BankWire {
	matches := KioskbotTypes.BankWire{
		Amount: getAmount(`Monto[\:\s\n]+[\$\s]+([\d\.]+)`, body),
		Name:   getName(`(?:Le informamos que nuestro\(a\) cliente|Le informamos que) (.+) (?:ha efectuado|le ha transferido)`, body),
		Bank:   "Banco Chile",
	}

	return matches
}

func Email() {
	r, err := imapreader.NewReader(imapreader.Options{
		Addr:     "imap.gmail.com",
		Username: os.Getenv("KB_EMAIL"),
		Password: os.Getenv("KB_EMAIL_PASSWORD"),
		TLS:      true,
		Timeout:  60 * time.Second,
		MarkSeen: true,
	})
	if err != nil {
		panic(err)
	}

	if err := r.Login(); err != nil {
		panic(err)
	}
	defer r.Logout()

	// Search for all the emails in "all mail" that are unseen
	// read the docs for more search filters
	messages, err := r.List(imapreader.GMailAllMail, imapreader.SearchUnseen)
	if err != nil {
		panic(err)
	}

	for _, v := range messages {
		fromSlice := v.Header["X-Original-Sender"]
		if fromSlice == nil {
			fromSlice = v.Header["From"]
			r, _ := regexp.Compile("<(.+)>")
			fromSlice[0] = r.FindStringSubmatch(fromSlice[0])[1]
		}
		from := fromSlice[0]

		parsedBody := strings.NewReader(string(v.Body))
		doc, err := goquery.NewDocumentFromReader(parsedBody)
		if err != nil {
			panic(err)
		}

		var wire KioskbotTypes.BankWire
		switch from {
		case "serviciodetransferencias@bancochile.cl":
			wire = matchBancoChile(doc.Text())
		case "mensajes@santander.cl":
			wire = matchBancoSantander(doc.Text())
		case "transferencias@bci.cl":
			wire = matchBancoBCI(doc.Text())
		case "informaciones@scotiabank.cl":
			wire = matchScotiabank(doc.Text())
		default:
			log.Printf("We couldnt find how to parse an email " + from)
			continue
		}

		Notify(wire)
	}
}
