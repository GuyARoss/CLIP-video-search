import torch
import uuid

def save_tensor(t: torch.Tensor) -> str:
    path = f'/tmp/{uuid.uuid4()}'
    torch.save(t, path)

    return path
