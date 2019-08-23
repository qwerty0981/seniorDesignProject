package main


import (
	"fmt"
	"log"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
)

type Room struct{
	status		string
	command		string
}

type statusMessage struct {
	RoomNumber	string	`json:"roomNumber"`
	Status		string	`json:"status,ommitempty"`
}

type commandMessage struct {
	Command		string	`json:"command"`
	RoomNumber	string	`json:"roomNumber"`
}

type statusResponse struct {
	Status		string	`json:"status"`
}

type commandResponse struct {
	Command		string	`json:"command"`
}

var rooms = map[string]*Room{"1": &Room{status:"none", command:"none"}}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/status", getStatus).Methods("GET")
	router.HandleFunc("/status", setStatus).Methods("POST")
	router.HandleFunc("/command", getCommand).Methods("GET")
	router.HandleFunc("/command", setCommand).Methods("POST")

	log.Fatal(http.ListenAndServe(":9090", router))
	fmt.Println("Done")
}

func setStatus(w http.ResponseWriter, r *http.Request) {
	var m statusMessage

	if r.Body == nil {
		fmt.Println("Empty body")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		text, _ := json.Marshal(map[string]string{"message": "Invalid Request"})
		w.Write(text)
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(400)
		text, _ := json.Marshal(map[string]string{"message": "Invalid Request"})
		w.Write(text)
		return
	}

	log.Println("Setting status for room number: " + m.RoomNumber + " to \"" + m.Status + "\"")

	rooms[m.RoomNumber].status = m.Status
	w.WriteHeader(200)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)

	m := statusMessage{}

	err := d.Decode(&m)
	if err != nil {
		w.WriteHeader(400)
		text, _ := json.Marshal(map[string]string{"message": "Invalid Request"})
		w.Write(text)
		return
	}
	response := statusResponse{Status: rooms[m.RoomNumber].status}
	res, _  := json.Marshal(response)
	w.WriteHeader(200)
	w.Write(res)
}

func getCommand(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)

	m := commandMessage{}

	err := d.Decode(&m)
	if err != nil {
		w.WriteHeader(400)
		text, _ := json.Marshal(map[string]string{"message": "Invalid Request"})
		w.Write(text)
		return
	}
	response := commandResponse{Command: rooms[m.RoomNumber].command}
	res, _  := json.Marshal(response)

	log.Println("Command for room number " + m.RoomNumber + " was read. Resetting it to 'None'")
	rooms[m.RoomNumber].command = "none"

	w.WriteHeader(200)
	w.Write(res)

}

func setCommand(w http.ResponseWriter, r *http.Request) {
	var m commandMessage

	if r.Body == nil {
		fmt.Println("Empty body")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		text, _ := json.Marshal(map[string]string{"message": "Invalid Request"})
		w.Write(text)
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(400)
		text, _ := json.Marshal(map[string]string{"message": "Invalid Request"})
		w.Write(text)
		return
	}

	log.Println("Setting commnad for room number: " + m.RoomNumber + " to \"" + m.Command + "\"")

	rooms[m.RoomNumber].command = m.Command
	w.WriteHeader(200)

}
