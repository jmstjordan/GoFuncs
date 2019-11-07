package gecko

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

type User struct {
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	Email        string `json:"Email"`
	EmployeeType string `json:"EmployeeType"`
	UserId       string `json:"UserId"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	user := r.URL.Query()["user"]

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "gecko-time")
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	ref := client.Doc("users/" + user[0])
	docsnap, err := ref.Get(ctx)
	if err != nil {
		log.Fatalf("Unable to read: %v", err)
	}
	dataMap := docsnap.Data()
	fmt.Println(dataMap)

	w.Header().Set("Content-Type", "application/json")
	returnUser := User{
		FirstName:    dataMap["FirstName"].(string),
		LastName:     dataMap["LastName"].(string),
		Email:        dataMap["Email"].(string),
		EmployeeType: dataMap["EmployeeType"].(string),
		UserId:       dataMap["UserId"].(string),
	}
	json.NewEncoder(w).Encode(returnUser)
}
