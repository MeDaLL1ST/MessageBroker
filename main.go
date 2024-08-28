package main

import (
	"log"
	"net/http"
	"os"
	"selfmq/metrics"
	"selfmq/pkg"
)

func main() {
	go metrics.StartMonitor()
	http.HandleFunc("/add", pkg.AddHandler)
	http.HandleFunc("/subscribe", pkg.SubscribeHandler)
	http.HandleFunc("/list", pkg.ListHandler)
	http.HandleFunc("/info", pkg.InfoHandler)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("WS_PORT"), nil))
}
func init() {
	pkg.Loadenv()
}
