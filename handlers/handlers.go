package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zusyed/qlik/dal"
	"github.com/zusyed/qlik/models"
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
