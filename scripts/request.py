
import leds_pb2_grpc as led_rpc
import leds_pb2 as led_proto
import grpc

import click

def send(req):
    with grpc.insecure_channel("192.168.2.25:50051") as c:
        stub = led_rpc.LedStripStub(c)
        resp = stub.Set(req)

@click.group()
def set_cmd():
    pass


@set_cmd.command()
def off():
    v = 0
    req = led_proto.SetPatternRequest(
        solid=led_proto.SolidPattern(color=v)
    )
    send(req)


@set_cmd.command()
@click.argument("color")
def solid(color):
    v = int(color, base=16) & 0xffffff

    req = led_proto.SetPatternRequest(
        solid=led_proto.SolidPattern(color=v)
    )

    send(req)


@set_cmd.command()
@click.option("--spawn-chance", type=float, default=0.6, help="chance [0, 1] of a snowflake spawning")
def snowflake(spawn_chance):
    req = led_proto.SetPatternRequest(
        snowflake=led_proto.SnowflakePattern(spawn_chance=spawn_chance)
    )
    send(req)

@set_cmd.command()
@click.argument("colors", nargs=-1)
def breathe(colors):
    assert len(colors) > 1

    colors = [
        int(c, base=16) & 0xffffff
        for c in colors
    ]

    req = led_proto.SetPatternRequest(
        breathe=led_proto.BreathePattern(colors=colors)
    )
    send(req)

@set_cmd.command()
@click.argument("colors", nargs=-1)
@click.option("--stripe-size", type=int)
def stripe(stripe_size, colors):
    vals = [
        int(c, base=16) & 0xffffff
        for c in colors
    ]
    req = led_proto.SetPatternRequest(
        stripe=led_proto.StripePattern(
            stripe_len=stripe_size,
            colors=vals
        )
    )
    send(req)

if __name__ == "__main__":
    set_cmd()
