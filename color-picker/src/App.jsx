import "./App.css";
import { useEffect, useState } from 'react';
import { Color } from './colorutils';
import { ColorPreview, SVSlider, HueSlider } from './colorPicker';
import { getStatus, setStatus } from './api';
import { PresetPage } from './components/PresetPage';
import { useNavigate, Routes, NavLink, Route, BrowserRouter } from 'react-router';

function setSolid(color) {
  return setStatus({
    type: "SOLID",
    data: {color: color.rgb().map(Math.trunc)}
  })
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

function Customizer() {
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
		<h2>Builder</h2>
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
  </div>;

}

function App() {
	function Rest() {
	const navigate = useNavigate();
	useEffect(() => navigate("/"), []);
		return <><div id="header">
			<NavLink index className={({isActive}) => isActive ? "nav-link active" : "nav-link"} to="/">Presets</NavLink>
			<NavLink className={({isActive}) => isActive ? "nav-link active" : "nav-link"} to="/builder">Builder</NavLink>
		</div>
		<Routes>
			<Route path="/" element={<PresetPage />}/>
			<Route path="/builder" element={<Customizer />}/>
		</Routes></>;
	}
	return <div id="root">
		<BrowserRouter> 
		 <Rest/>
	</BrowserRouter></div>;
}

export default App;
