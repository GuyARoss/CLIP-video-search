from transformers import CLIPProcessor, CLIPModel
import torch

class FrameProcessor:
    def __init__(self) -> None:
        self.device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

        self.model = CLIPModel.from_pretrained("openai/clip-vit-base-patch32").to(self.device)
        self.processor = CLIPProcessor.from_pretrained("openai/clip-vit-base-patch32")

    def text_probability_from_tensor_paths(self, serial_image_tensor_paths: str, text: str) -> str:
        image_tensor_paths = serial_image_tensor_paths.split(' ')

        avg_sum = 0.0
        for image_tensor_path in image_tensor_paths:
            image_tensor = torch.load(image_tensor_path)
            image_tensor = torch.load(image_tensor_path).to(self.device)

            inputs = self.processor(text=text, return_tensors="pt", padding=True).to(self.device)

            inputs['pixel_values'] = image_tensor    
            outputs = self.model(**inputs)

            logits_per_image = outputs.logits_per_image    
            probs = logits_per_image.squeeze()

            avg_sum += probs.item()

        return str(avg_sum / len(image_tensor_paths))