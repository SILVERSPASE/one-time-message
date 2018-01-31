package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Welcome!  ˆ_ˆ")
	if askForConfirmation("Do you want to open link via browser? y/n") {
		openbrowser("http://localhost:8000/messages")
	}
	router := mux.NewRouter()
	router.HandleFunc("/messages", GetMessages).Methods("GET")
	router.HandleFunc("/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/messages/{link}", GetMessage).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

var messages []Message

func askForConfirmation(str string) bool {
	fmt.Println(str)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	if response == "y" {
		return true
	} else {
		return false
	}

}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
