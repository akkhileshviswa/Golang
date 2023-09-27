package main

import (
	"broker/event"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
	"time"
)

type RequestPayload struct {
	Action    string           `json:"action"`
	Subscribe SubscribePayload `json:"submit,omitempty"`
}

type SubscribePayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	ID        int
	FirstName string `json:"first_name"`
	Name      string `json:"name"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Is_sent   bool      `json:"is_sent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Response struct {
	Result []User
}

// This function is used to handle the /subscribe request.
func (app *Config) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "subscribe":
		app.subscribe(w, requestPayload.Subscribe)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

// This function is used to make an rpc call to insert the record.
func (app *Config) subscribe(w http.ResponseWriter, subscribePayload SubscribePayload) {
	client, err := rpc.Dial("tcp", ":8083")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var response int
	err = client.Call("RPCServer.InsertRecord", subscribePayload, &response)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Account has been created",
	}

	logMessage := LogPayload{
		Name: "Insert",
		Data: "Record Created: " + strconv.Itoa(response),
	}

	app.logEvent(logMessage)
	app.writeJSON(w, http.StatusAccepted, payload)
}

// This function pushes the message to RabbitMQ for logging
func (app *Config) logEvent(msg LogPayload) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	message, err := json.MarshalIndent(&msg, "", "\t")
	if err != nil {
		log.Println(err)
		return err
	}

	err = emitter.Push(string(message), "log.INFO")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// This function sends mail to the users who are created newly, via cron job.
func (app *Config) SendMail() error {
	log.Println("Inside send Mail")
	client, err := rpc.Dial("tcp", ":8083")
	if err != nil {
		log.Println(err)
		return err
	}

	var response Response

	err = client.Call("RPCServer.GetMailUserRecords", struct{}{}, &response)
	if err != nil {
		return err
	}

	log.Println(response.Result)
	for _, item := range response.Result {
		if (item != User{}) {
			msg := MailPayload{
				ID:        item.ID,
				FirstName: item.FirstName,
				Name:      "AccountCreated",
				To:        item.Email,
				Subject:   "Account Created",
				Message:   "Thank you for signing up. You will start recieving Newsletter once we start publishing!!",
			}

			emitter, err := event.NewEventEmitter(app.Rabbit)
			if err != nil {
				return err
			}

			message, err := json.MarshalIndent(&msg, "", "\t")
			if err != nil {
				log.Println(err)
				return err
			}

			err = emitter.Push(string(message), "Mail")
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}
