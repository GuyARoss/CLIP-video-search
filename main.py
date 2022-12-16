import sys
import uuid
from search import Search
from video_features import video_features
from indexes.video_scene import VideoSceneIndex
from indexes.scene_embedding import SceneEmbeddingIndex
from indexes.video_metadata import VideoMetadataIndex

scene_index = VideoSceneIndex("./dbs/scene_index")
scene_embedding_index = SceneEmbeddingIndex('./dbs/scene_embedding_index')
video_metadata_index = VideoMetadataIndex('./dbs/videometadata_index')

search_instance = Search(scene_index, scene_embedding_index, video_metadata_index)
def search(input: str) -> int:
    results = search_instance.from_input(input)

    for s in results:
        print('Result', s['video_id'], s['score'], s['video_uri'])

    return 0

def insert(video_path: str) -> int:
    features = video_features(video_path)
    video_id = str(uuid.uuid4())

    video_metadata_index.insert(video_id, {        
        'VideoURI': video_path,
    })

    scene_ids = []
    for f in features['scene_features']['clip_pixel_scenes']:
        scene_id = str(uuid.uuid4())
        scene_ids.append(scene_id)
        with open(f['local_path'], mode='rb') as file:
            content = file.read()
            
            scene_embedding_index.insert(scene_id, content)

    scene_index.insert(video_id, scene_ids)


if __name__ == "__main__":
    r = {
        'search': search,
        'insert': insert,
    }

    fn = r.get(sys.argv[1])
    if fn == None:
        Exception('invalid type')

    input = " ".join(sys.argv[2:])
    
    SystemExit(fn(input))