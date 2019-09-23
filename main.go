package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	opsGinieContacts := readCSV("contacts.csv")

	contactType := []string{"voice", "sms"}

	for _, contact := range opsGinieContacts {
		for _, v := range contactType {
			if v == "voice" {
				createContact(contact.Email, v, contact.VoiceTo)
			} else {
				createContact(contact.Email, v, contact.SmsTo)
			}
		}
	}
}

func createContact(username string, method string, to string) error {
	url := fmt.Sprintf("https://api.opsgenie.com/v2/users/%s/contacts", username)

	httpClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	data := map[string]string{"method": method, "to": to}

	jsonValue, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "GenieKey xxxx-xxxx-xxxx-xxxx")
	req.Header.Set("Content-Type", "application/json")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	contactCreated := func() string {
		if res.StatusCode == 201 {
			return "success"
		}
		return "failed"
	}()

	bodyString := string(body)

	fmt.Println(fmt.Sprintf("status: %s, email: %s, respose: %s", contactCreated, username, bodyString))

	return nil
}

func readCSV(csvFileName string) []Contact {
	csvFile, _ := os.Open(csvFileName)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var contacts []Contact
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		contacts = append(contacts, Contact{
			Email:   line[0],
			VoiceTo: line[1],
			SmsTo:   line[2],
		})
	}

	return contacts
}

type Contact struct {
	Email   string
	VoiceTo string
	SmsTo   string
}
