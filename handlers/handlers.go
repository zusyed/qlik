package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"log"

	"github.com/gorilla/mux"
	"github.com/zusyed/qlik/dal"
	"github.com/zusyed/qlik/models"
)

const (
	jsonContent = "application/json; charset=UTF-8"
)

type (
	//Handler responds to HTTP requests
	Handler struct {
		db dal.MessageDB
	}
)

//NewHandler creates a Handler object with the specified db
func NewHandler(db dal.MessageDB) *Handler {
	return &Handler{
		db: db,
	}
}

//GetMessages gets all the messages
func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.db.GetMessages()
	if err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered an error getting messages from database")

		return
	}

	if messages == nil {
		// return an empty array if there are no records instead of returning nil
		messages = []models.Message{}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(messages); err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered an error encoding JSON response")

		return
	}
}

//CreateMessage adds the specified message
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var err error

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse request body: %s", err)

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("Encountered an error closing the request body: %+v", err)
		}
	}()

	var message models.Message

	if err := json.Unmarshal(body, &message); err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not parse JSON: %s", err)

		return
	}

	m, err := h.db.AddMessage(message)
	if err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not insert message in the database: %s", err)

		return
	}

	w.Header().Set("Content-Type", jsonContent)
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(m); err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not encode message: %s", err)

		return
	}
}

//DeleteMessage removes the message with the specified id
func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int

	idStr := mux.Vars(r)["id"]
	if id, err = strconv.Atoi(idStr); err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id must be an integer")

		return
	}

	err = h.db.DeleteMessage(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		if err.Error() == dal.NotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Could not find record with id %d", id)

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not delete message with id %d: %s", id, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

//GetMessage gets the message with the specified id
func (h *Handler) GetMessage(w http.ResponseWriter, r *http.Request) {
	var err error
	var id int

	idStr := mux.Vars(r)["id"]
	if id, err = strconv.Atoi(idStr); err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "id must be an integer")

		return
	}

	message, err := h.db.GetMessage(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		if err.Error() == dal.NotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Could not find record with id %d", id)

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not get message with id %d: %s", id, err)

		return
	}

	message.IsPalindrome = isPalindrome(message.Body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(message); err != nil {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Encountered an error encoding JSON response")

		return
	}
}

//isPalindrome returns true if the specified string is a palindrome; otherwise returns false
func isPalindrome(s string) bool {
	n := len(s)
	for i := 0; i < (n / 2); i++ {
		if s[i] != s[n-i-1] {
			return false
		}
	}

	return true
}
