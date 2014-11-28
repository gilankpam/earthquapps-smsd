package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	customerKey       string = "bestapp145"
	customerSecret           = "PMBNX"
	accessToken              = "ec127bae9f0f61d4cf8491238fb7b554"
	accessTokenSecret        = "c25e7fb09ba14effc32907385a51dbc7"
	smsEndPoint              = "http://sandbox.appprime.net/TemanDev/rest/sendSmsBulk/"
)

var (
	httpClient *http.Client = &http.Client{}
	client     oauth.Client = oauth.Client{
		Credentials: oauth.Credentials{
			Token:  customerKey,
			Secret: customerSecret,
		},
	}
	credentials    *oauth.Credentials = &oauth.Credentials{accessToken, accessTokenSecret}
	urlSmsEndPoint *url.URL
)

type smsBulk struct {
	sms sms `json:"smsBulk"`
}

type sms struct {
	number  string `json:"msisdn"`
	message string `json:"message"`
}

func newSmsBulk(number, message string) *smsBulk {
	return &smsBulk{
		sms{number, message},
	}
}

func main() {
	Call("081226906673", "test")
}

func Call(number, message string) {
	sms := newSmsBulk(number, message)
	smsJson, err := json.Marshal(sms)
	if err != nil {
		log.Fatalf("Error json marshal: %s", err)
	}
	body := bytes.NewReader(smsJson)
	req, err := http.NewRequest("POST", smsEndPoint, body)
	if err != nil {
		log.Fatal("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", client.AuthorizationHeader(credentials, "POST", urlSmsEndPoint, nil))
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal("Error making request: %v", err)
	}
	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(resBody))

}

func init() {
	var err error
	urlSmsEndPoint, err = url.Parse(smsEndPoint)
	if err != nil {
		log.Panicf("error parsing url: %v", err)
	}
}
