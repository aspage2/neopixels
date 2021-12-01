from typing import Callable, Sequence
from flask import Flask, request, abort
import board
import neopixel
import time

app = Flask(__name__)

def lerp(c1, c2, t):
    return tuple(int(t * y + (1 - t) * x) for x, y in zip(c1, c2))


Color = Sequence[int]
Pattern = Callable[[int], Color]


def solid(c: Color) -> Pattern:
    return lambda _: c


def sequence(cs: Sequence[Color]) -> Pattern:
    l = len(cs)
    return lambda i: cs[i % l]


def gradient(fst: Color, snd: Color) -> Pattern:
    l = len(pixels)
    return lambda i: lerp(fst, snd, i / l)


pattern = solid((0,0,0))


def transition(p1: Pattern, p2: Pattern):
    l = len(pixels)
    for t in range(21):
        for i in range(l):
            c = lerp(p1(i), p2(i), t / 20)
            pixels[i] = c
        pixels.show()
        time.sleep(0.01)


@app.route("/set/", methods=["POST"])
def set():
    global pattern 
    global pixels
    data = request.json

    if "color" in data:
        newpattern = solid(data["color"])
    elif "sequence" in data:
        newpattern = sequence(data["sequence"])
    elif "gradient" in data:
        newpattern = gradient(*data["gradient"])

    transition(pattern, newpattern)
    pattern = newpattern
    
    return "", 201


@app.route("/off/", methods=["POST"])
def off():
    global pixels
    global pattern
    newpattern = solid((0,0,0))
    transition(pattern, newpattern)
    pattern = newpattern
    return "", 201

pixels = neopixel.NeoPixel(
    board.D18,
    39, 
    brightness=0.3,
    auto_write=False,
    pixel_order="RGB",
)

