package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GuyARoss/clip-video-search/internal/engine"
	"github.com/GuyARoss/clip-video-search/pkg/queue"
)

type EngineWebController struct {
	engine *engine.Engine
	queue  queue.Queue
}

type SearchRequest struct {
	Input      string `json:"input"`
	MaxResults int    `json:"maxResults"`
}

func (c *EngineWebController) TopNFromText(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	req := &SearchRequest{}
	dec.Decode(req)

	if req.MaxResults == 0 {
		req.MaxResults = 3
	}

	response := c.engine.TopNFromText(req.Input, req.MaxResults)

	rout, _ := json.Marshal(response)
	w.Write(rout)
}

func (c *EngineWebController) Insert(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	req := &engine.VideoItem{}

	dec.Decode(req)
	c.queue.Add(req)

	w.Write([]byte(fmt.Sprintf(`{ "queueSize": %d }`, c.queue.Length())))
}

func NewEngineController(engine *engine.Engine, queue queue.Queue) *EngineWebController {
	return &EngineWebController{
		engine: engine,
		queue:  queue,
	}
}
