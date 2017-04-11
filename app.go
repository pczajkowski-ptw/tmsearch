package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Application stores main information needed to run the app
type Application struct {
	Name, AccessToken, Sid, BaseURL, AuthString string
	Languages                                   map[string]string
	Delay                                       time.Duration
}

// SetBaseURL sets base URL for API endpoint.
func (app *Application) SetBaseURL(url string) {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	app.BaseURL = url
}

// JSONDecoder decodes json to given interface, borrowed from SO.
func JSONDecoder(data io.ReadCloser, target interface{}) {
	decoder := json.NewDecoder(data)

	err := decoder.Decode(target)
	if err != nil {
		log.Printf("Error reading json: %v", err)
	}
}

// LoadLanguages loads languages from languages.json to map.
func (app *Application) LoadLanguages() {
	data, err := os.Open("./html/languages.json")
	defer data.Close()
	if err != nil {
		log.Fatalf("Error reading languages: %v", err)
	}

	app.Languages = make(map[string]string)
	JSONDecoder(data, &app.Languages)
}

// CheckLanguage checks if language is in the map.
func (app Application) CheckLanguage(language string) bool {
	_, ok := app.Languages[language]
	if !ok {
		return false
	}

	return true
}

// Login logs into the API and sets AuthString.
func (app *Application) Login() {
	credentials, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		log.Fatalf("Error reading credentials: %v", err)
	}

	loginURL := app.BaseURL + "auth/login"

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(credentials))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Error logging in: %v", resp.Status)
	}
	defer resp.Body.Close()

	JSONDecoder(resp.Body, &app)

	app.AuthString = "?authToken=" + app.AccessToken
	log.Println(app.AuthString, resp.Status)
}
