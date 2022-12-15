package engine

import (
	"fmt"
	"io/ioutil"

	sceneembedding "github.com/GuyARoss/clip-video-search/internal/index/scene_embedding"
	videometadata "github.com/GuyARoss/clip-video-search/internal/index/video_metadata"
	sceneIndex "github.com/GuyARoss/clip-video-search/internal/index/video_scene"
	"github.com/GuyARoss/clip-video-search/pkg/inference"
	sortedlist "github.com/GuyARoss/clip-video-search/pkg/sorted_list"
	"github.com/GuyARoss/clip-video-search/pkg/util"
	"github.com/google/uuid"
)

type VideoItem struct {
	VideoURI string `json:"videoURI"`
}

type Engine struct {
	VideoSceneIndex     sceneIndex.VideoSceneIndex
	SceneEmbeddingIndex sceneembedding.SceneEmbeddingIndex
	VideoMetadataIndex  videometadata.VideoMetadataIndex
	inference           inference.InferenceImplementation
}

func (s *Engine) Next(next interface{}) {
	nextItem := next.(*VideoItem)

	if nextItem.VideoURI == "" {
		fmt.Println("empty video uri")
		return
	}

	response, err := s.inference.VideoFeatures(nextItem.VideoURI)
	if err != nil {
		fmt.Println(err)
		return
	}

	videoId := uuid.NewString()

	s.VideoMetadataIndex.InsertByVideoID(videoId, videometadata.VideoMetaData{
		VideoDuration: response.VideoDuration,
		VideoURI:      nextItem.VideoURI,
	})

	sceneIds := []string{}
	for _, pixelData := range response.SceneFeatures.ClipPixelScenes {
		sceneID := uuid.NewString()
		f, err := ioutil.ReadFile(pixelData.LocalPath)
		if err != nil {
			fmt.Println(err)
		}

		sceneIds = append(sceneIds, sceneID)
		s.SceneEmbeddingIndex.InsertSceneEmbedding(sceneID, f)
	}

	s.VideoSceneIndex.InsertVideoScenes(videoId, sceneIds)
	fmt.Println("scenes inserted")
}

type TransientResult struct {
	// TODO(feature): would also be beneficial to understand which scene is the best fit
	VideoID string
	Score   float64
}

type EngineResult struct {
	*TransientResult
	VideoURI string
}

type topNVisitor struct {
	*Engine
	input   string
	Results sortedlist.SortableList
}

func (s *topNVisitor) InferSceneScore(sceneIDs []string) (float64, error) {
	tensorPaths := make([]string, len(sceneIDs))

	for idx, sceneID := range sceneIDs {
		cachePath := util.TensorCachedPath(sceneID)
		if cachePath != "" {
			tensorPaths[idx] = cachePath
			continue
		}

		tensor := s.SceneEmbeddingIndex.GetTensorBySceneID(sceneID)
		tensorPath, err := util.TensorBytesToFile(sceneID, tensor)

		if err != nil {
			fmt.Println("cannot read tensor for", sceneID)
			continue
		}

		tensorPaths[idx] = tensorPath
	}

	return s.inference.FrameTextProcessor(tensorPaths, s.input)
}

func (t *topNVisitor) Visit(vis sceneIndex.VideoSceneRecord) error {
	v, err := t.InferSceneScore(vis.SceneIDs)
	if err != nil {
		return err
	}

	t.Results.MaybeAdd(v, &TransientResult{
		Score:   v,
		VideoID: vis.VideoID,
	})

	return nil
}

func (e *Engine) TopNFromText(text string, limit int) []EngineResult {
	visitorInstance := &topNVisitor{
		input:   text,
		Engine:  e,
		Results: sortedlist.New(limit),
	}

	err := e.VideoSceneIndex.ForEach(visitorInstance)
	if err != nil {
		fmt.Println("error propagated in record", err)
	}

	recommendations := make([]EngineResult, limit)

	for i, r := range visitorInstance.Results.Results() {
		if r != nil {
			unpacked := r.(*TransientResult)

			metadata := e.VideoMetadataIndex.GetByVideoID(unpacked.VideoID)
			recommendations[i] = EngineResult{
				TransientResult: unpacked,
				VideoURI:        metadata.VideoURI,
			}
		}
	}

	return recommendations
}

func NewDefaultEngine(inference inference.InferenceImplementation) (*Engine, error) {
	sceneIndex, err := sceneIndex.New("./dbs/scene_index")
	if err != nil {
		return nil, err
	}

	sceneEmbeddingIndex, err := sceneembedding.New("./dbs/scene_embedding_index")
	if err != nil {
		return nil, err
	}

	videoMetadataIndex, err := videometadata.New("./dbs/videometadata_index")
	if err != nil {
		return nil, err
	}

	return &Engine{
		inference:           inference,
		VideoSceneIndex:     sceneIndex,
		SceneEmbeddingIndex: sceneEmbeddingIndex,
		VideoMetadataIndex:  videoMetadataIndex,
	}, nil
}
