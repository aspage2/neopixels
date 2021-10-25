from flask import Flask, request, abort
import board
import neopixel

app = Flask(__name__)



@app.route("/set/", methods=["POST"])
def set():
    print(request.data)
    print(request.headers)
    data = request.json

    global pixels
    color = data.get("color")
    if color is None:
        abort(400, "request must include 'color' prop")
    
    pixels.fill(color)
    pixels.show()

    return "", 201


@app.route("/off/", methods=["POST"])
def off():
    global pixels

    pixels.fill([0,0,0])
    pixels.show()
    return "", 201


pixels = neopixel.NeoPixel(board.D18, 39, brightness=80, auto_write=False)


if __name__ == "__main__":
    app.run(host="0.0.0.0", port="5000")
