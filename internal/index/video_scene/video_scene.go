package videoscene

import (
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

type VideoSceneRecord struct {
	VideoID  string
	SceneIDs []string
}

type RecordVisitor interface {
	Visit(VideoSceneRecord) error
}

type VideoSceneIndex interface {
	InsertVideoScenes(videoID string, sceneIds []string)
	ForEach(RecordVisitor) error
}

type LevelVideoSceneIndex struct {
	db *leveldb.DB
}

func (i *LevelVideoSceneIndex) InsertVideoScenes(videoID string, sceneIds []string) {
	i.db.Put([]byte(videoID), []byte(strings.Join(sceneIds, ",")), nil)
}

func (i *LevelVideoSceneIndex) ForEach(vis RecordVisitor) error {
	iter := i.db.NewIterator(nil, nil)
	wg := &sync.WaitGroup{}

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		r := VideoSceneRecord{
			VideoID:  string(key),
			SceneIDs: strings.Split(string(value), ","),
		}

		wg.Add(1)
		go func(g *sync.WaitGroup) {
			vis.Visit(r)
			wg.Done()
		}(wg)
		// err := vis.Visit(r)
		// if err != nil {
		// 	return err
		// }
	}
	wg.Wait()
	return nil
}

func New(path string) (VideoSceneIndex, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelVideoSceneIndex{
		db: db,
	}, nil
}
