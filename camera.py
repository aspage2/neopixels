from cv2 import VideoCapture, destroyWindow, imwrite, waitKey

cam = VideoCapture(0)

result, img = cam.read()
if not result:
    print("can't take image")

imwrite("image.png", img)

