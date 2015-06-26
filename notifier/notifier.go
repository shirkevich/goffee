package notifier

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goffee/goffee/data"
	"github.com/goffee/goffee/queue"
	"github.com/keighl/mandrill"
	"github.com/parnurzeal/gorequest"
	"bitbucket.org/ckvist/twilio/twirest"
)

var (
	exit           = make(chan bool)
	MandrillKey    string
	mandrillClient *mandrill.Client
	SlackUrl       string
	TwilioSid			 string
	TwilioToken		 string
	TwilioFromNumber string
	GroupEmail		 string
	Phones				 string
)

func Run() {
	go run()
}

func run() {
	mandrillClient = mandrill.ClientWithKey(MandrillKey)

	for {
		notifications := queue.FetchNotifications()
		for _, n := range notifications {
			checkId, err := strconv.ParseInt(n, 10, 64)

			check, err := data.FindCheck(checkId)
			if err != nil {
				continue
			}

			user, err := check.User()
			if err != nil {
				continue
			}

			sendMessage(check, user)
		}
	}
}

func sendMessage(c data.Check, u data.User) {
	var subject string

	if c.Success {
		subject = fmt.Sprintf("Up: %s (%d)", c.URL, c.Status)
	} else {
		subject = fmt.Sprintf("Down: %s (%d)", c.URL, c.Status)
	}

	html := `<strong>%s</strong>
  <br>
  <br>
  <p>Checked at %s by <a href='http://goffee.io/'>Goffee.io</a></p>`
	html = fmt.Sprintf(html, subject, c.UpdatedAt.Format(time.Kitchen))

	text := `%s\n\nChecked at %s by Goffee.io`
	text = fmt.Sprintf(text, subject, c.UpdatedAt.Format(time.Kitchen))

	message := &mandrill.Message{}
	message.AddRecipient(GroupEmail, "Alerts channel", "to")
	message.FromEmail = "noreply@excursiopedia.org"
	message.FromName = "Goffee Notifier"
	message.Subject = subject
	message.HTML = html
	message.Text = text

	responses, err := mandrillClient.MessagesSend(message)
	if err != nil {
		fmt.Printf("%s - %s\n", err, responses)
	} else {
		fmt.Printf("Notifying via email: %s\n", GroupEmail)
	}

	request := gorequest.New()
	jsonMessage := `{"text":"` + text + `", "channel":"#ops", "username":"goffee"}`
	_, _, errs := request.Post(SlackUrl).Send(jsonMessage).End()
	if errs != nil {
		fmt.Printf("%s", errs)
	} else {
		fmt.Printf("Sent notification to ops channel\n")
	}

	client := twirest.NewClient(TwilioSid, TwilioToken)

	msg := twirest.SendMessage{
					Text: text,
					To:   Phones,
					From: TwilioFromNumber}

	resp, err := client.Request(msg)
	if err != nil {
					fmt.Printf("%s", err)
	} else {
		fmt.Printf("Message to: %v %v", Phones, resp.Message.Status)
	}

}

func Wait() {
	<-exit
}
