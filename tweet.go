package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
)

var (
	httpClient *http.Client
	URL        = "https://api.twitter.com/2"
	cred       *Credentials
)

// Credentials stores all access/consumer tokens
// and secret keys needed for authentication against
// the twitter REST API.
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}
type Createtweetrequest struct {
	Text string `json:"text,omitempty"`
}

// CreateTweetData is the data returned when creating a tweet
type CreateTweetData struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type CreateTweetResponse struct {
	Tweet *CreateTweetData `json:"data"`
}

// return http.Client types to automatically authorize http.Request's
func (creds *Credentials) GetClientToken() (client *http.Client) {
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	return config.Client(oauth1.NoContext, token)
}

func init() {
	cred = &Credentials{
		ConsumerKey:       os.Getenv("TWITTER_SECRET_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_SECRET"),
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_SECRET"),
	}

}

func CreateTweet(text string) (response *CreateTweetResponse, err error) {
	endpoint := "/tweets"
	path := URL + endpoint
	tweet := Createtweetrequest{
		Text: text,
	}
	data, err := json.Marshal(&tweet)
	if err != nil {
		log.Println("Error encoding to json", err)
	}

	httpClient = cred.GetClientToken()
	res, err := httpClient.Post(path, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error making request: ", err)
	}

	fmt.Println(res.Status)
	// decode the response body to bytes
	respByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var resp CreateTweetResponse

	if err := json.Unmarshal(respByte, &resp); err != nil {
		log.Println("Error decoding response::: ", err)
	}

	return &resp, nil
}
