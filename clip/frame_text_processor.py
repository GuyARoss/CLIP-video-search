import typing
import torch
from transformers import CLIPProcessor, CLIPModel

class FrameProcessor:
    def __init__(self) -> None:
        self.device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

        self.model = CLIPModel.from_pretrained("openai/clip-vit-base-patch32").to(self.device)
        self.processor = CLIPProcessor.from_pretrained("openai/clip-vit-base-patch32")

    def frame_text_processor(self, tensors: typing.List[torch.Tensor], text: str) -> float:
        avg_sum = 0.0
        for tensor in tensors:
            image_tensor = torch.load(tensor)
            image_tensor = image_tensor.to(self.device)

            inputs = self.processor(text=text, return_tensors="pt", padding=True).to(self.device)

            inputs['pixel_values'] = image_tensor   
            outputs = self.model(**inputs)

            logits_per_image = outputs.logits_per_image    
            probs = logits_per_image.squeeze()

            avg_sum += probs.item()

        return avg_sum / len(tensors)