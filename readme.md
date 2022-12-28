# CLIP Video Search
[CLIP (Contrastive Language–Image Pre-training)](https://openai.com/blog/clip/) is a technique _which efficiently learns visual concepts from natural language supervision_. CLIP has found applications in [stable diffusion](https://github.com/huggingface/diffusers/tree/main/examples/community#clip-guided-stable-diffusion).

This repository aims act as a POC in exploring the ability to use CLIP for video search using natural language outlined in the article found [here](https://medium.com/@guyallenross/using-clip-to-build-a-natural-language-video-search-engine-6498c03c40d2).

## Usage
### Dependencies
- [libzmq](https://github.com/zeromq/libzmq)
- python >= 3.8
- go >= 1.18

### Running
1. start up the inference zmq server found in the `./inference` directory `python3 zmq_server.py`.
2. start up the go server with `go run main.go`.

### Example
Before running this example, please ensure that your environment is correctly configured and the application is running without errors.

1. index the video clips provided by the `examples/videos`.
```bash
curl -X POST -d '{"videoURI": "<path_to_dir>/examples/videos/<video_name>.mp4" }' http://localhost:3000/insert 
```

__note__: it can take a moment for the video to become searchable.

2. then search for a video
```bash
curl -X POST -d '{"input": "a man cutting pepper", "maxResults": 1 }' http://localhost:3000/search
```

## TODO
- [ ] CLI to remove manual setup process
- [ ] ability to add dedicated inference machines (currently limited to same host)