import cv2
import typing
import scenedetect as sd
from PIL import Image
import torch
from utils.tensor import save_tensor

from transformers import CLIPProcessor, CLIPModel

class FrameNumTimecode():
    def __init__(self, frame_num: int) -> None:
        self.frame_num = frame_num

class SceneFeatures:
    def __init__(self) -> None:
        self.clip_processor = CLIPProcessor.from_pretrained("openai/clip-vit-base-patch32")
        self.clip_model = CLIPModel.from_pretrained("openai/clip-vit-base-patch32")

    def collect_scenes_in_video(self, video_path: str) -> typing.List[typing.Tuple[sd.FrameTimecode, sd.FrameTimecode]]:
        video = sd.open_video(video_path)
        sm = sd.SceneManager()
        
        sm.add_detector(sd.ContentDetector(threshold=27.0))
        sm.detect_scenes(video)

        return sm.get_scene_list()

    def clip_embeddings(self, image: Image):
        inputs = self.clip_processor(images=image, return_tensors="pt", padding=True)
        input_tokens = {
            k: v for k, v in inputs.items()
        }
        return input_tokens['pixel_values']

    def clip_features_to_dic(self, num_of_scenes: int, clip_pixel_scenes: typing.List, scenes: typing.List[typing.Tuple[sd.FrameTimecode, sd.FrameTimecode]]) -> typing.Dict[str, any]:
        d = {}
        d['num_of_scenes'] = num_of_scenes        
        d['clip_pixel_scenes'] = [{
            'local_path': save_tensor(s['pixel_embeddings']),
            'scene': {
                'start_frame_num': scenes[s['scene_no']][0].frame_num,
                'end_frame_num': scenes[s['scene_no']][1].frame_num,
            }
        } for s in clip_pixel_scenes]
        return d

    def scene_features(self, video_path: str, no_of_samples: int = 3) -> typing.Dict:
        scenes = self.collect_scenes_in_video(video_path)

        cap = cv2.VideoCapture(video_path)
        scenes_frame_samples = []    
        for scene_idx in range(len(scenes)):
            scene_length = abs(scenes[scene_idx][0].frame_num - scenes[scene_idx][1].frame_num)
            every_n = round(scene_length/no_of_samples)
            local_samples = [(every_n * n) + scenes[scene_idx][0].frame_num for n in range(3)]
            
            scenes_frame_samples.append(local_samples)

        if len(scenes) == 0:
            # this could denote a single contiguous scene.
            frame_count = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
            if frame_count > 0:
                every_n = round(frame_count/no_of_samples)
                local_samples = [(every_n * n) for n in range(3)]
                scenes_frame_samples.append(local_samples)
                scenes = [(FrameNumTimecode(0),FrameNumTimecode(frame_count))]

        scene_clip_embeddings = []
        for scene_idx in range(len(scenes_frame_samples)):
            scene_samples = scenes_frame_samples[scene_idx]

            pixel_tensors = []
            for frame_sample in scene_samples:
                cap.set(1, frame_sample)
                ret, frame = cap.read()
                if not ret:
                    print('breaks oops', ret, frame_sample, scene_idx, frame)
                    break

                pil_image = Image.fromarray(frame)
                
                clip_pixel_values = self.clip_embeddings(pil_image)
                pixel_tensors.append(clip_pixel_values)

            scene_clip_embeddings.append({ 'pixel_embeddings': torch.mean(torch.stack(pixel_tensors), dim=0), 'scene_no': scene_idx })            

        return self.clip_features_to_dic(len(scenes), scene_clip_embeddings, scenes)
