package sms

import (
	"fmt"

	"net/url"
	"strings"
)

// TODO At present we only support sending via twilio
//  we should put services into adapters, and load the relevant one
// FIXME: this should create a type, which we attach these vars/functions to
// rather than using package variables (see view)

// The SMS service user (must be set before first sending)
var user string

// The SMS service secret key/password (must be set before first sending)
var secret string

// The phone no to use as a from value
var from string

// Setup sets the user and secret for use in sending mail (possibly later we should have a config etc)
func Setup(u string, s string, f string) {
	user = u
	secret = s
	from = f
}

// Send sends an SMS
func Send(recipient string, message string) error {

	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", user)
	data := url.Values{}
	data.Set("From", from)
	data.Set("To", recipient)
	data.Set("Body", message)
	body := data.Encode()

	fmt.Printf("SENDING SMS:%s BODY:%s\n", endpoint, body)

	request, err := http.NewRequest("POST", endpoint, strings.NewReader(body))
	if err != nil {
		return err
	}
	request.SetBasicAuth(user, secret)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		fmt.Sprintf("Error sending SMS:%v\n", response)
		return fmt.Errorf("SMS send failed with status:%d", response.StatusCode)
	}

	return nil

}
