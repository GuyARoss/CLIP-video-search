from urllib import request
import subprocess

def get_video_duration(filename: str) -> float:
    try:
        result = subprocess.run(["ffprobe", "-v", "error", "-show_entries",
                                "format=duration", "-of",
                                "default=noprint_wrappers=1:nokey=1", filename],
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT)
        return float(result.stdout)

    except:
        return 0

def generic_web_downloader(url: str, path: str, filename: str) -> str:
    try:
        if not url and ('https://' not in url or 'http://' not in url):
            raise Exception('video URL not found')
    
        with open(f'{path}/{filename}.mp4', 'wb') as video_file:
            video_file.write(request.urlopen(url).read())
        
        return f'{path}/{filename}.mp4'

    except Exception as e:
        raise Exception('error: web downloader: ', e.__str__())

