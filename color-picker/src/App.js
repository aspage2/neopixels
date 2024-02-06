import './App.css';
import { useEffect, useState } from 'react';

function rgbToHsv(r, g, b) {
  r /= 255; g /= 255; b /= 255;

  let max = Math.max(r, g, b), min = Math.min(r, g, b);
  let h, s, v = max;

  let d = max - min;
  s = max === 0 ? 0 : d / max;

  if (max === min) {
    h = 0; // achromatic
  } else {
    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
      default: throw Error("What?");
    }

    h /= 6;
  }

  return [ h, s, v ];
}

function hsvToRgb(h, s, v) {
  let r, g, b;

  let i = Math.floor(h * 6);
  let f = h * 6 - i;
  let p = v * (1 - s);
  let q = v * (1 - f * s);
  let t = v * (1 - (1 - f) * s);

  switch (i % 6) {
    case 0: r = v; g = t; b = p; break;
    case 1: r = q; g = v; b = p; break;
    case 2: r = p; g = v; b = t; break;
    case 3: r = p; g = q; b = v; break;
    case 4: r = t; g = p; b = v; break;
    case 5: r = v; g = p; b = q; break;
    default: throw Error("not good");
  }

  return [ r * 255, g * 255, b * 255 ];
}

function hex2RGB(colorString) {
  return [
    parseInt(colorString.substring(1, 3), 16),
    parseInt(colorString.substring(3, 5), 16),
    parseInt(colorString.substring(5, 7), 16),
  ]
}

function getStatus() {
  return fetch("http://192.168.1.74:5000/status/", {mode: "cors"})
    .then(resp => resp.json())
}

function setStatus(data) {
    return fetch(
        "http://192.168.1.74:5000/status/",
        {
            mode: "cors",
            method: "POST",
            body: JSON.stringify(data),
            headers: {
                "Content-Type": "application/json",
            }
        },
    );
}

function rgb2Hex([r, g, b]) {
  return "#"+Math.trunc(r).toString(16).padStart(2, "0")
    + Math.trunc(g).toString(16).padStart(2, "0")
    + Math.trunc(b).toString(16).padStart(2, "0");
}

function textColor(rgb) {
    const brightness = Math.round(((parseInt(rgb[0]) * 299) +
    (parseInt(rgb[1]) * 587) +
    (parseInt(rgb[2]) * 114)) / 1000);
    if (brightness > 125) {
      return "black";
    }
    return "white";
}

function App() {
  useEffect(() => {
    getStatus().then(res => {
      setColor(res.data.color)
    });
  }, []);
  const [color, setColor] = useState([58,115,83]);
  const colorHex = rgb2Hex(color);
  const colorHsl = rgbToHsv(...color);
  return <>
    <p>{(colorHsl[2]*100).toFixed(2)}% Brightness</p>
    <div style={{
      background: colorHex || "black",
      width: "400px",
      height: "80px",
      textAlign: "center"
    }}>
      <p style={{
        lineHeight: "80px",
        color: textColor(color),
        fontWeight: "bold",
        fontSize: "1.1rem"
        }}>{colorHex}</p>
    </div>
    <input 
      type="color" 
      value={colorHex} 
      onChange={e => {
        setColor(hex2RGB(e.target.value));
      }}
    />
    <button onClick={() => 
      setStatus({type:"SOLID", data: {color: color.map(Math.trunc)}})}
      >Submit</button>
    <button onClick={() => {
      let [h, s, l] = colorHsl;
      let newL = l - 0.05;
      if (newL < 0) {
        newL = 0
      }
      let newColor = hsvToRgb(h, s, newL);
      setStatus({type:"SOLID", data: {color: newColor.map(Math.trunc)}}).then(
        () => setColor(newColor)
      );
    }}>
      - Brightness
    </button>
    <button onClick={() => {
      let [h, s, l] = colorHsl;
      let newL = l + 0.05;
      if (newL > 1) {
        newL = 1.
      }
      let newColor = hsvToRgb(h, s, newL);
      setStatus({type:"SOLID", data: {color: newColor.map(Math.trunc)}}).then(
        () => setColor(newColor)
      );
    }}>
      + Brightness
    </button>
    <button onClick={() => {
      setStatus({type:"SOLID", data: {color: [0, 0, 0]}}).then(
      () => setColor([0, 0, 0])
      );
    }}>OFF</button>
  </>;
}

export default App;