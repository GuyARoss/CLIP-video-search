import typing
import leveldb

class VideoSceneRecord():
    def __init__(self, videoID: str, sceneIDs: typing.List[str]):
        self.videoID = videoID
        self.sceneIDs = sceneIDs


class VideoSceneIndex():
    def __init__(self, path) -> None:
        self.level_instance = leveldb.LevelDB(path)

    def collect_records(self) -> typing.List[VideoSceneRecord]:
        r = []
        for k, v in self.level_instance.RangeIter():
            record = VideoSceneRecord(k.decode('utf-8'), str(v.decode('utf-8')).split(','))
            r.append(record)

        return r

    def insert(self, videoID, sceneIds: typing.List[str]):
        b = ",".join(sceneIds)
        self.level_instance.Put(videoID.encode('utf-8'), b.encode('utf-8'))