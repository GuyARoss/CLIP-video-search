import leveldb
from io import BytesIO

class SceneEmbeddingIndex():
    def __init__(self, path: str) -> None:
        self.level_instance = leveldb.LevelDB(path)

    def get_tensor_by_id(self, id: str):
        b = self.level_instance.Get(bytes(id,'utf-8'))
        return BytesIO(b)

    def insert(self, sceneID, data):
        self.level_instance.Put(sceneID.encode('utf-8'), data)