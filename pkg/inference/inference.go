package inference

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/GuyARoss/clip-video-search/pkg/zmqpool"
)

type PixelScene struct {
	LocalPath string `json:"local_path"`
	Scene     struct {
		StartFrameNum int `json:"start_frame_num"`
		EndFrameNum   int `json:"end_frame_num"`
	} `json:"scene"`
}

type SceneFeatures struct {
	NumOfScenes     int          `json:"num_of_scenes"`
	ClipPixelScenes []PixelScene `json:"clip_pixel_scenes"`
}

type FullVideoFeatures struct {
	SceneFeatures SceneFeatures `json:"scene_features"`
	VideoDuration float32       `json:"video_duration"`
}

type InferenceImplementation interface {
	VideoFeatures(string) (*FullVideoFeatures, error)
	FrameTextProcessor([]string, string) (float64, error)
}

type InferenceServer struct {
	ipc zmqpool.ZMQPoolImplementation
	m   *sync.Mutex
}

type InferenceOperations string

const (
	VideoFeatures      InferenceOperations = "video_features"
	FrameTextProcessor InferenceOperations = "frame_text_processor"
)

func (s *InferenceServer) FrameTextProcessor(tensorPath []string, text string) (float64, error) {
	ipcResponse, err := s.ipc.Send(string(FrameTextProcessor), strings.Join(tensorPath, " "), text)
	if err != nil {
		return 0.0, err
	}

	score, _ := strconv.ParseFloat(ipcResponse, 32)

	return score, nil
}

func (s *InferenceServer) VideoFeatures(uri string) (*FullVideoFeatures, error) {
	ipcResponse, err := s.ipc.Send(string(VideoFeatures), uri)
	if err != nil {
		return nil, err
	}

	features := &FullVideoFeatures{}
	json.Unmarshal([]byte(ipcResponse), features)

	return features, nil
}

func New(ipc zmqpool.ZMQPoolImplementation) InferenceImplementation {
	return &InferenceServer{
		ipc: ipc,
	}
}
