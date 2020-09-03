package tests

import (
	"fmt"
	"github.com/dataxpe/goinsta/v2"
	"log"
	"os"
	"testing"
	"github.com/tcnksm/go-input"
)

func TestImportAccount(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("logged into Instagram as user '%s'", insta.Account.Username)
}

func TestLogin(t *testing.T) {
	insta, err := getRandomAccount()
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := insta.Login(); err != nil {
		switch v := err.(type) {
		case goinsta.ChallengeError:
			err := insta.Challenge.Process(v.Challenge.APIPath)
			if err != nil {
				log.Fatalln(err)
			}

			ui := &input.UI{
				Writer: os.Stdout,
				Reader: os.Stdin,
			}

			query := "What is SMS code for instagram?"
			code, err := ui.Ask(query, &input.Options{
				Default:  "000000",
				Required: true,
				Loop:     true,
			})
			if err != nil {
				log.Fatalln(err)
			}

			err = insta.Challenge.SendSecurityCode(code)
			if err != nil {
				log.Fatalln(err)
			}
		}

		log.Fatalln(err)
	}
	defer insta.Logout()

	t.Logf("logged into Instagram as user '%s'", insta.Account.Username)
}

func TestGenerateEncPw(t *testing.T) {
	enc, err := goinsta.GenerateEncPassword(
		"PASSWORD",
		"169",
		"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF1dmtwUmcvNlpQV0hPM1lEclpxegpCN1NGNHkwL0Nhc1dDME9yNG9wYkMrOEs0MDJrMEZZRWEwQ1AyRE5VdjBrZTB4UFVYdnFwUXFhaERXSUlvOUI0Clk3RlFEOS9kbjZoWHk1SU1CUy9URnZ1em9ZNVpKUlhVd21yVVE5QW5oUnVLMGtUSTVjUjI3eXBYeDI3cjBwV0sKWXpROUdmMm9hTmg3N3lLTzB1RUtiTUZGUkFNMFBqQzRjUHNEU0RWVUNNeUt5cy94NVJCa3MxaVlLdFlQaEw0cQp6Vys3bHROYXpGY09QcThYOW1lR01XbVdaL3ZNWG1GeHM2QWp6ZGsreEt5aC9QSFRQQW1neklKcW9xNE1xN2l5Cm9ZaHpmVXo0M2RWb3lEd1BhNm5LcmpoMFdESlljWHFnTmZuVnNLbmQ2cXBnU3pvendOb0s0bFMyS0pxUExmU2YKZlFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==",
		)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	fmt.Printf("enc: %s\n",enc)
}

func TestGenerateJazoest(t *testing.T) {
	fmt.Printf("jazo: %s\n",goinsta.GenerateJazoest("e7f7c586-f884-42e2-af9a-b813a2b4cebb"))
}