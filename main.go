package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dghubble/oauth1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	urlPost = "https://api.twitter.com/1.1/statuses/update.json"
)

type (
	// Credential
	Credential struct {
		ConsumerKey    string `json:"consumer_key"`
		ConsumerSecret string `json:"consumer_secret"`
		AccessToken    string `json:"access_token"`
		AccessSecret   string `json:"access_secret"`
	}
	// DataCredential
	DataCredential struct {
		Credential []Credential `json:"credentials"`
	}
)

// NewConfig
func NewConfig(cKey, CScret, aToken, aSecret string) *Credential {
	return &Credential{
		ConsumerKey:    cKey,
		ConsumerSecret: CScret,
		AccessToken:    aToken,
		AccessSecret:   aSecret,
	}
}

// NewHttp
func (c Credential) NewHttp() *http.Client {
	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return httpClient
}

// GetConfigFile
func GetConfigFile(j chan<- Credential) {
	jsonFile, err := os.Open("account.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var cred DataCredential
	json.Unmarshal(byteValue, &cred)

	for _, y := range cred.Credential {
		j <- Credential{
			ConsumerKey:    y.ConsumerKey,
			ConsumerSecret: y.ConsumerSecret,
			AccessToken:    y.AccessToken,
			AccessSecret:   y.AccessSecret,
		}
	}
	close(j)
}

func main() {

	var status = flag.String("text", "halo", "Guys..")
	flag.Parse()

	jobs := make(chan Credential, 0)

	go func() {
		GetConfigFile(jobs)
	}()

	for g := range jobs {

		c := NewConfig(g.ConsumerKey,
			g.ConsumerSecret,
			g.AccessToken,
			g.AccessSecret)

		textString := strings.Replace(*status, " ", "%20", -1)
		path := urlPost + "?status=" + textString

		resp, err := c.NewHttp().Post(path,
			"application/json",
			nil)

		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			_, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(resp.Status)
		} else {
			fmt.Println(resp.Status)
		}
	}
}
