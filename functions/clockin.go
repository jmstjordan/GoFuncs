package gecko

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

type Clock struct {
	UserId    string `json:"UserId"`
	ClockType string `json:"ClockType"`
}
type ClockPost struct {
	ClockTime time.Time              `json:"ClockIn"`
	ClockType string                 `json:"ClockType"`
	User      *firestore.DocumentRef `json:"User"`
}

func ClockIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Methot must be POST", http.StatusMethodNotAllowed)
		return
	}
	var clock Clock
	if err := json.NewDecoder(r.Body).Decode(&clock); err != nil {
		http.Error(w, "Unable to decode Body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "gecko-time")
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	usrRef := client.Doc("users/" + clock.UserId)
	_, _, err = client.Collection("clocks").Add(ctx, ClockPost{
		ClockType: clock.ClockType,
		ClockTime: time.Now(),
		User:      usrRef,
	})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to Create Document", http.StatusBadRequest)
		return
	}
	w.WriteHeader(201)
}
