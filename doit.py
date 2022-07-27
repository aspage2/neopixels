import time
import requests
from cv2 import VideoCapture, destroyWindow, imwrite, waitKey

URL = "http://192.168.2.16:4000"

def take(fn: str):
    cam = VideoCapture(0)
    result, img = cam.read()
    
    if not result:
        raise Exception("can't take image")
    imwrite(f"images/{fn}", img)

def set(n: int):
    requests.post(URL, json={"n": n}).raise_for_status()
def off():
    requests.post(URL, json={}).raise_for_status()
def all():
    requests.post(URL, json={"n": 1000}).raise_for_status()

off()
time.sleep(0.5)
take("off.png")
for i in range(100):
    set(i)
    time.sleep(0.5)
    take(f"image{i}.png")
all()
time.sleep(0.5)
take("all.png")
