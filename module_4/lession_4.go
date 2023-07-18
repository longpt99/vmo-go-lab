// main.go
package module4

import (
	"encoding/json"
	"log"
	"net/http"
)

func Lesson_4() {
	http.HandleFunc("/", HandlerRequest)
	log.Println("Server is starting")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Error creating server %s\n", err.Error())
	}
}

func HandlerRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	})
}
