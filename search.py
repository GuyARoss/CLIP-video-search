from concurrent.futures import ThreadPoolExecutor
from functools import partial
from clip.frame_text_processor import FrameProcessor
from indexes.video_scene import VideoSceneIndex, VideoSceneRecord
from indexes.scene_embedding import SceneEmbeddingIndex
from indexes.video_metadata import VideoMetadataIndex

class Search():
    def __init__(self, scene_index: VideoSceneIndex, scene_embedding_index: SceneEmbeddingIndex, video_metadata_index: VideoMetadataIndex) -> None:
        self.scene_index = scene_index
        self.scene_embedding_index = scene_embedding_index
        self.video_metadata_index = video_metadata_index

        self.clip = FrameProcessor()

    def _scene_averages(self, r: VideoSceneRecord, input: str):
        tensors = [self.scene_embedding_index.get_tensor_by_id(id) for id in r.sceneIDs]
        prob = self.clip.frame_text_processor(tensors, input)

        return (prob, r.videoID)

    def from_input(self, input: str, top_n=1):
        records = self.scene_index.collect_records()        
        f = partial(self._scene_averages, input=input)

        with ThreadPoolExecutor() as thread_pool:
            results = thread_pool.map(f, records)    
            sorted = list(results)
            sorted.sort(key=lambda x: x[0], reverse=True)
            
            results = []
            for s in sorted[:top_n]:            
                data = self.video_metadata_index.get_by_id(s[1])
                
                results.append({
                    'video_id': s[1],
                    'score': s[0],
                    'video_uri': data['VideoURI']
                })

        return results