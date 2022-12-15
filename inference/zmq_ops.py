from video_features import video_features
from clip.frame_text_processor import FrameProcessor
fp = FrameProcessor()

registry = {}
registry['video_features'] = video_features
registry['frame_text_processor'] = fp.text_probability_from_tensor_paths