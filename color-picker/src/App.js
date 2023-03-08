import './App.css';
import {useState, useRef, useEffect} from "react";
import clamp from "lodash/clamp";

function useSelectBar(coordPos) {
    const [coord, setCoord] = useState(coordPos || 0);
    const ref = useRef();
    const [mouseDown, setMouseDown] = useState(false);

    const onMouseMove = ev => {
        if (!mouseDown) {
            return
        }
        const bounds = ref.current.getBoundingClientRect();
        setCoord(clamp(ev.clientX - bounds.left, 0, 500));
    };

    const onMouseDown = () => {
        setMouseDown(true);
    };

    const onMouseUp = () => {
        setMouseDown(false);
    }

    const onTouchMove = ev => {
        if (!mouseDown) {
            return
        }
        console.log("touch move");
        const bounds = ref.current.getBoundingClientRect();
        setCoord(clamp(ev.changedTouches[0].clientX - bounds.left, 0, 500));
    };

    const onTouchStart = ev => {
        console.log("touch start");
        console.log(ev)
        setMouseDown(true);
    }

    const onTouchEnd = ev => {
        setMouseDown(false);
        console.log("touch end");
        console.log(ev);
    }

    return [coord, {
        ref,
        onMouseDown,
        onTouchStart,
    }, {
        onTouchEnd,
        onTouchMove,
        onMouseMove,
        onMouseUp,
        style: {display: mouseDown ? 'inline-block' : 'none'}
    }];
}

function hslToRGB(h, s, l) {
    // Must be fractions of 1
    s /= 100;
    l /= 100;

    let c = (1 - Math.abs(2 * l - 1)) * s,
        x = c * (1 - Math.abs((h / 60) % 2 - 1)),
        m = l - c/2,
        r = 0,
        g = 0,
        b = 0;
    
    if (h < 60) {
        r = c; g = x; b = 0;  
    } else if (h < 120) {
        r = x; g = c; b = 0;
    } else if (h < 180) {
        r = 0; g = c; b = x;
    } else if (h < 240) {
        r = 0; g = x; b = c;
    } else if (h < 300) {
        r = x; g = 0; b = c;
    } else if (h < 360) {
        r = c; g = 0; b = x;
    }
    r = Math.round((r + m) * 255);
    g = Math.round((g + m) * 255);
    b = Math.round((b + m) * 255);
    return [r, g, b];
    
}

function App() {

    const [hueCoord, hueEventProps, hueOverlayProps] = useSelectBar(0);
    const [satCoord, satEventProps, satOverlayProps] = useSelectBar(500);
    const [lightCoord, lightEventProps, lightOverlayProps] = useSelectBar(250);

    const hue = (360 * hueCoord) / 500;
    const sat = (100 * satCoord) / 500;
    const light = (100 * lightCoord) / 500;

    const sendColor = () => {
        fetch(
            "http://192.168.2.16:5000/status/",
            {
                method: "POST",
                body: JSON.stringify({
                    type: "SOLID",
                    color: hslToRGB(hue, sat, light)
                }),
                headers: {
                    "Content-Type": "application/json",
                }
            },
        );
    };

    return <div className="App">

        <div id="hsl-gradient" className="color-bar" {...hueEventProps}>
            <div className="color-bar-select" style={{left: `${hueCoord}px`}}/>
        </div><br/>
        <div className="secret-overlay" {...hueOverlayProps} />

        <div id="sat-gradient" className="color-bar" style={{
            background: `linear-gradient(90deg, hsl(${hue}, 0%, ${light}%), hsl(${hue}, 100%, ${light}%))`
        }} {...satEventProps}>
            <div className="color-bar-select" style={{left: `${satCoord}px`}}/>
        </div><br/>
        <div className="secret-overlay" {...satOverlayProps}/>
        <div id="lightness-gradient" className="color-bar" style={{
            background: `linear-gradient(90deg, hsl(${hue}, ${sat}%, 0%), hsl(${hue}, ${sat}%, 50%), hsl(${hue}, 100%, 100%))`
        }} {...lightEventProps}>
            <div className="color-bar-select" style={{left: `${lightCoord}px`}}/>
        </div><br/>
        <div className="secret-overlay" {...lightOverlayProps}/>
        <div id="preview" style={{background: `hsl(${hue}, ${sat}%, ${light}%)`}}/>
        <button id="submit" onClick={sendColor}>Submit</button>
    </div>;
}

export default App;
