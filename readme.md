# CLIP Video Search
[CLIP (Contrastive Language–Image Pre-training)](https://openai.com/blog/clip/) is a technique "which efficiently learns visual concepts from natural language supervision.". CLIP has found applications in [stable diffusion](https://github.com/huggingface/diffusers/tree/main/examples/community#clip-guided-stable-diffusion).

This repository aims act as a POC in exploring the ability to use CLIP for video search using natural language outlined in the article found [here](https://medium.com/@guyallenross/using-clip-to-build-a-natural-language-video-search-engine-6498c03c40d2). 

Please check out the `master` branch of this repo, to find a more optimized example of this code. 


## Usage
### Dependencies
- python >= 3.8


This repository includes the indexes for the 3 example videos included in this repo. 
### Indexing/ Searching a video
1. To index a video, use `python3 insert <path_to_video>`.
2. To search for a video, use `python3 search <natural search>`.
