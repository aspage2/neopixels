import './App.css';
import { useEffect, useState } from 'react';
import { Color, textColor, hex2RGB } from './colorutils';
import { ColorPreview, SVSlider, HueSlider } from './colorPicker';

async function getStatus() {
  const resp = await fetch("http://192.168.0.232:5000/status/", { mode: "cors" });
  return await resp.json();
}

function setStatus(data) {
	return fetch(
		"http://192.168.0.232:5000/status/",
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

function setSolid(color) {
  return setStatus({
    type: "SOLID",
    data: {color: color.rgb().map(Math.trunc)}
  })
}

function setSequence(colors) {
	return setStatus({
		type: "SEQUENCE", 
		data: { colors }
	});
}

function BrightnessAdjust({color, setColor}) {
  const [, , brightness] = color.hsv();
  return <div style={{display: "block", margin: "10px"}}>
    <button onClick={() => 
      setColor(color.adjustBrightness(-0.05))
    }>- Brightness</button>
    <span>{(brightness*100).toFixed(2)}%</span>
    <button onClick={() =>
      setColor(color.adjustBrightness(0.05))
    }>+ Brightness</button>
  </div>
}

const PRESETS = {
  "OFF": [0, 0, 0],
	"\u{273F}": hex2RGB("#586593"),
	"Mintel": hex2RGB("#fedb00"),
	"BB <GO>": hex2RGB("#FFA028"),
};

function App() {
  useEffect(() => {
    getStatus().then(res => {
	  if (res.type === "SEQUENCE")
		  console.log("lol");
	  else {
			const newColor = Color.fromArr(res.data.color);
		  setColor(newColor);
		  setSelectedColor(newColor);
		}
    });
  }, []);

	const [selectedColor, setSelectedColor] = useState(new Color({}));
	const [color, setColor] = useState(new Color({}));
	const [sh, ss, sv] = selectedColor.hsv();
	const [h] = color.hsv()

	function onSVChange(newSV) {
		setColor(color.setHSV(newSV));
	}

  return <div style={{maxWidth: 400}}>
		<div style={{
			display: "block",
			textAlign: "center",
		}}>
			<ColorPreview color={color} setColor={c => {setColor(c); setSelectedColor(c); setSolid(c)}}/>
			<div style={{margin: 8}}/>
			<HueSlider height={40} hue={sh}
				onHSVChange={hsv => {
					setColor(color.setHSV(hsv));
				}}
				onHSVCommit={hsv => {
					const nc = selectedColor.setHSV(hsv);
					setSelectedColor(nc);
					setSolid(nc);
				}}
			/>
			<div style={{margin: 8}}/>
			<SVSlider width={340} height={240} hue={h} sv={{s: ss, v: sv}} onChange={onSVChange} 
			onRelease={hsv => {
					const nc = selectedColor.setHSV(hsv);
					setSelectedColor(nc);
					setSolid(nc);
			}}/>
			<br/>
			<button onClick={() => setSolid(color)}>Submit</button>
			<BrightnessAdjust 
				color={color}
				setColor={newColor => setSolid(newColor).then(() => {setColor(newColor);setSelectedColor(newColor)})}
			/>
		</div>
    <h3>Presets</h3>
		<div>
			{
				Object.entries(PRESETS).map(([colorName, rgb]) => {
					const c = Color.fromArr(rgb);
					return <button 
						key={colorName} 
						onClick={()=>setSolid(c).then(() => {
							setColor(c); setSelectedColor(c)
						})}
						style={{
							"background": c.hex(),
							"color": textColor(rgb),
							"marginRight": 8,
						}}
					>{colorName}</button>
				}) 
			}
			<button onClick={() => {
			setSequence([[255,0,0],[252,186,3],[207,207,0],[0,255,0],[0,0,255],[85,0,171],[187,0,250]]);
		}} className="rainbow-border">{"\u{1F485}"}</button>
		</div>
  </div>;
}

export default App;
