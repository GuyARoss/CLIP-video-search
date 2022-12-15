package videometadata

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

type VideoMetadataIndex interface {
	InsertByVideoID(videoID string, data VideoMetaData)
	GetByVideoID(string) *VideoMetaData
}

type LevelVideoMetadataIndex struct {
	db *leveldb.DB
}

type VideoMetaData struct {
	VideoDuration float32
	VideoURI      string
}

func (i *LevelVideoMetadataIndex) InsertByVideoID(videoID string, data VideoMetaData) {
	b, _ := json.Marshal(data)

	i.db.Put([]byte(videoID), b, nil)
}

func (i *LevelVideoMetadataIndex) GetByVideoID(videoID string) *VideoMetaData {
	b, _ := i.db.Get([]byte(videoID), nil)

	md := &VideoMetaData{}
	json.Unmarshal(b, md)

	return md
}

func New(path string) (VideoMetadataIndex, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelVideoMetadataIndex{
		db: db,
	}, nil
}
