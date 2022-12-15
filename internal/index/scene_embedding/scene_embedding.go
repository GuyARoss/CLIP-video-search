package sceneembedding

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type SceneEmbeddingIndex interface {
	InsertSceneEmbedding(sceneID string, tensor []byte)
	GetTensorBySceneID(string) []byte
}

type LevelSceneEmbeddingIndex struct {
	db *leveldb.DB
}

func (i *LevelSceneEmbeddingIndex) InsertSceneEmbedding(sceneID string, tensor []byte) {
	i.db.Put([]byte(sceneID), tensor, nil)
}

func (i *LevelSceneEmbeddingIndex) GetTensorBySceneID(sceneID string) []byte {
	b, _ := i.db.Get([]byte(sceneID), nil)

	return b
}

func New(path string) (SceneEmbeddingIndex, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelSceneEmbeddingIndex{
		db: db,
	}, nil
}
