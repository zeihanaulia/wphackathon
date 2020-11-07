package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/zeihanaulia/instagram-scraper/config"
	"github.com/zeihanaulia/instagram-scraper/internal/engine"
	"github.com/zeihanaulia/instagram-scraper/internal/engine/ports"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	engine.Start(config.Configurations{
		SHEETY_API:         os.Getenv("SHEETY_API"),
		INSTAGRAM_USERNAME: os.Getenv("INSTAGRAM_USERNAME"),
		INSTAGRAM_PASSWORD: os.Getenv("INSTAGRAM_PASSWORD"),
	})

	response := ports.ResponseHashTag{
		SyncAt: time.Now().Format(time.RFC3339),
		Status: "Requested",
	}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(js))
}
