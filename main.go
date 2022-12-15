package main

import (
	"net/http"

	"github.com/GuyARoss/clip-video-search/internal/engine"
	"github.com/GuyARoss/clip-video-search/internal/web"
	"github.com/GuyARoss/clip-video-search/pkg/inference"
	"github.com/GuyARoss/clip-video-search/pkg/queue"
	"github.com/GuyARoss/clip-video-search/pkg/zmqpool"
)

func main() {
	q := queue.New()

	pool := zmqpool.New([]string{
		"tcp://localhost:5550",
		"tcp://localhost:5551",
		"tcp://localhost:5552",
	})
	server := inference.New(pool)

	engine, err := engine.NewDefaultEngine(server)
	if err != nil {
		panic(err)
	}

	go queue.LongLivedIterator(q, engine)

	engineController := web.NewEngineController(engine, q)

	http.HandleFunc("/search", engineController.TopNFromText)
	http.HandleFunc("/insert", engineController.Insert)

	http.ListenAndServe(":3000", nil)
}
