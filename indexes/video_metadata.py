import leveldb
import json

class VideoMetadataIndex():
    def __init__(self, path: str) -> None:
        self.level_instance = leveldb.LevelDB(path)
    
    def get_by_id(self, id: str):
        b = self.level_instance.Get(bytes(id,'utf-8'))
        return json.loads(b.decode('utf-8'))

    def insert(self, videoID: str, data):
        b = json.dumps(data)

        self.level_instance.Put(videoID.encode('utf-8'), b.encode('utf-8'))