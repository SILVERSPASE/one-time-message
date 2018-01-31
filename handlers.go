package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if len(messages) > 0 {
		json.NewEncoder(w).Encode(messages)
	} else {
		fmt.Fprint(w, "curl -d '{\"Message\": \"Hello World\"}' -H \"Content-Type: application/json\" -X POST http://localhost:8000/messages")
	}
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message Message
	_ = json.NewDecoder(r.Body).Decode(&message)
	message.PersistentLink = uuid.Must(uuid.NewV4())
	message.OneTimeLink = uuid.Must(uuid.NewV4())
	message.SeenDate = time.Time{}                       // prevent of setting time while creating message
	secret := []byte("a_very_very_very_very_secret_key") // 32 bytes
	plaintext := []byte(message.Message)
	ciphertext, err := encrypt(secret, plaintext) //[]uint8
	if err != nil {
		log.Fatal(err)
	}
	message.CipherMessage = string(ciphertext)
	messages = append(messages, message)
	json.NewEncoder(w).Encode(messages)
	var boblink string = "http://localhost:8000/messages/" + message.OneTimeLink.String()
	json.NewEncoder(w).Encode(boblink)
	fmt.Println(boblink)
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for key, item := range messages {
		link, err := uuid.FromString(params["link"])
		if err != nil {
			fmt.Fprint(w, "Link error")
		}
		if item.OneTimeLink == link {
			if item.SeenDate.IsZero() { // if not seen date - show message
				var bobMessage string
				bobMessage = item.CipherMessage
				keys, ok := r.URL.Query()["secret"]
				if !ok || len(keys) < 1 {
					json.NewEncoder(w).Encode(bobMessage)
					fmt.Fprint(w, "Add secret to url(?secret=a_very_very_very_very_secret_key)")
					return
				}
				secret := []byte(keys[0])
				result, err := decrypt(secret, []byte(item.CipherMessage))
				if err != nil {
					log.Fatal(err)
				}
				bobMessage = string(result)
				json.NewEncoder(w).Encode(bobMessage)
				messages[key].SeenDate = time.Now() //set seen date
			} else {
				fmt.Fprint(w, "Link has been expired")
			}
			return
		} else if item.PersistentLink == link {
			json.NewEncoder(w).Encode(item)
		} else {
			fmt.Fprint(w, "Error")
		}
	}
}
