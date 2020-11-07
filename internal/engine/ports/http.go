package ports

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeihanaulia/instagram-scraper/config"
	"github.com/zeihanaulia/instagram-scraper/internal/engine"
)

type HttpServer struct {
	Cfg config.Configurations
}

func NewHttpServer(cfg config.Configurations) HttpServer {
	return HttpServer{cfg}
}

type ResponseHashTag struct {
	SyncAt string `json:"sync_at,omitempty"`
	Status string `json:"status,omitempty"`
}

func (h HttpServer) SyncHashtag(w http.ResponseWriter, r *http.Request) {
	go func() {
		engine.Start(h.Cfg)
	}()

	response := ResponseHashTag{
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
