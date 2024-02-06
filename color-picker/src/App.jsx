import './App.css';
import { useEffect, useState } from 'react';
import { rgb2Hex, rgbToHsv, textColor, hex2RGB, adjustBrightness} from './colorutils';

async function getStatus() {
  const resp = await fetch("http://192.168.1.74:5000/status/", { mode: "cors" });
  return await resp.json();
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

function setSolid(rgb) {
  return setStatus({
    type: "SOLID",
    data: {color: rgb.map(Math.trunc)}
  })
}

function ColorPreview({color}) {
  const _tc = textColor(color);
  const hexColor = rgb2Hex(color);

  return <div style={{
    background: hexColor || "black",
    width: "400px",
    height: "80px",
    textAlign: "center"
  }}>
    <p style={{
      lineHeight: "80px",
      color: _tc,
      fontWeight: "bold",
      fontSize: "1.1rem"
      }}>{hexColor}</p>
  </div>
}

function BrightnessAdjust({color, setColor}) {
  const [, , brightness]= rgbToHsv(...color);
  return <div style={{display: "block", margin: "10px"}}>
    <button onClick={() => 
      setColor(adjustBrightness(color, -0.05))
    }>- Brightness</button>
    <span>{(brightness*100).toFixed(2)}% Brightness</span>
    <button onClick={() =>
      setColor(adjustBrightness(color, +0.05))
    }>+ Brightness</button>
  </div>
}

const PRESETS = {
  "OFF": [0, 0, 0],
  "3000 K": hex2RGB("#cb905d"),
};

function App() {
  useEffect(() => {
    getStatus().then(res => {
      setColor(res.data.color)
    });
  }, []);
  const [color, setColor] = useState([58,115,83]);
  const colorHex = rgb2Hex(color);
  return <>
    <input 
      type="color" 
      value={colorHex} 
      onChange={e => {
        setColor(hex2RGB(e.target.value));
      }}
    />
    <button onClick={() => setSolid(color)}>Submit</button>
    <ColorPreview color={color} />
    <BrightnessAdjust 
      color={color}
      setColor={newColor => setSolid(newColor).then(() => setColor(newColor))}
    />
    <h3>Presets</h3>
    <div>{
      Object.entries(PRESETS).map(([colorName, rgb]) => 
        <button 
          key={colorName} 
          onClick={()=>setSolid(rgb).then(() => setColor(rgb))}
          style={{
            "background-color": rgb2Hex(rgb),
            "color": textColor(rgb),
          }}
        >{colorName}</button>
      ) 
    }</div>
  </>;
}

export default App;