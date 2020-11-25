
import leds_pb2_grpc as led_rpc
import leds_pb2 as led_proto
import grpc
import time


def lerp(c1, c2, t):
    return tuple(int(t*(_b - _a) + _a) for _a, _b in zip(c1, c2))

def toint(c):
    return (c[2]&0xff) | ((c[1]&0xff)<<8) | ((c[0]&0xff)<<16)


def colors():
    c1 = (0, 255, 0)
    c2 = (0, 0, 255)

    i = 0
    while i < 1:
        yield led_proto.SetColorRequest(color=toint(lerp(c1, c2, i)))
        i += 0.01
        time.sleep(0.01)


with grpc.insecure_channel("192.168.2.25:50051") as c:
    stub = led_rpc.LedStripStub(c)
    stub.SetRealtime(colors())


