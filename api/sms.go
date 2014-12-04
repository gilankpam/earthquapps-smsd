package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"io"
	// "io/ioutil"
	"log"
	// "net"
	"net/http"
	"net/url"
	// "time"
)

const (
	smsEndPoint string = "http://sandbox.appprime.net/TemanDev/rest"
)

var (
	httpClient *http.Client
	smsBulkUrl *url.URL
)

type Client struct {
	oauthClient *oauth.Client
	credential  *oauth.Credentials
}

func NewClient(cusKey, cusSecret, accessToken, accessTokenSecret string) *Client {
	client := new(Client)
	client.oauthClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  cusKey,
			Secret: cusSecret,
		},
	}
	client.credential = &oauth.Credentials{accessToken, accessTokenSecret}
	return client
}

type Response struct {
	Data struct {
		Code    string `json:"returnCode"`
		Message string `json:"returnMessage"`
	} `json:"data"`
}

type SmsBulk struct {
	Sms struct {
		Number  string `json:"msisdn"`
		Message string `json:"message"`
	} `json:"smsBulk"`
}

func (c *Client) SendSMSBulk(number, message string) (*Response, error) {
	endPoint := smsEndPoint + "/sendSmsBulk/"
	//This is ugly
	smsBulk := SmsBulk{Sms: struct {
		Number  string `json:"msisdn"`
		Message string `json:"message"`
	}{number, message}}
	jsonString, err := json.Marshal(smsBulk)
	if err != nil {
		log.Printf("Error json marshal: %s", err)
	}
	req, err := c.makeRequest(endPoint, bytes.NewReader(jsonString))
	if err != nil {
		log.Printf("Cannot create new request: %v", err)
		return nil, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer res.Body.Close()
	rres, err := getResponse(res)
	if err != nil {
		return nil, err
	}
	fmt.Println(rres)
	return rres, nil
}

func getResponse(res *http.Response) (*Response, error) {
	decoder := json.NewDecoder(res.Body)
	response := new(Response)
	err := decoder.Decode(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) makeRequest(url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.oauthClient.AuthorizationHeader(c.credential, "POST", smsBulkUrl, nil))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func init() {
	/*timeout := time.Duration(2 * time.Second)
	dialTimeout := func(network, address string) (net.Conn, error) {
		return net.DialTimeout(network, address, timeout)
	}*/

	httpClient = &http.Client{}
	smsBulkUrl, _ = url.Parse(smsEndPoint + "/sendSmsBulk/")
}
