package main

import (
    "net/http"
    "encoding/json"
    "log"
)

func healthz(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    err := json.NewEncoder(w).Encode(map[string]string{
        "status":   "ok",
        "instance": "dsl-node-01",
    })

    if err != nil {
        log.Printf("write health response: %v", err)
    }
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /healthz", healthz)
    err := http.ListenAndServe("127.0.0.1:8080", mux)
    
    if err != nil {
        log.Fatalf("server failed: %v", err)
    }
}