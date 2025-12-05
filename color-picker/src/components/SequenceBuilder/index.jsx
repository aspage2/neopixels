import "./index.css";

import { useEffect, useState } from "react";
import { Color, textColor } from "../../colorutils";
import { getStatus, setStatus } from "../../api";
import {Sequence} from "../../pattern";
import { ColorPreview, HueSlider, SVSlider } from "../../colorPicker";

const DEFAULT_SEQUENCE = [
		[255, 0, 0],
		[252, 186, 3],
		[207, 207, 0],
		[0, 255, 0],
		[0, 0, 255],
		[85, 0, 171],
		[187, 0, 250]
].map(Color.fromArr);

function Customizer({color: initialColor, onSubmit=()=>{}}) {
  useEffect(() => {
		setColor(color);
		setSelectedColor(color);
	} , [initialColor]);

	const [selectedColor, setSelectedColor] = useState(new Color({}));
	const [color, setColor] = useState(initialColor);
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
			<ColorPreview color={color} setColor={c => {setColor(c); setSelectedColor(c);}}/>
			<div style={{margin: 8}}/>
			<HueSlider height={40} hue={sh}
				onHSVChange={hsv => {
					setColor(color.setHSV(hsv));
				}}
				onHSVCommit={hsv => {
					const nc = selectedColor.setHSV(hsv);
					setSelectedColor(nc);
				}}
			/>
			<div style={{margin: 8}}/>
			<SVSlider width={340} height={240} hue={h} sv={{s: ss, v: sv}} onChange={onSVChange} 
			onRelease={hsv => {
					const nc = selectedColor.setHSV(hsv);
					setSelectedColor(nc);
			}}/>
			<br/>
			<button onClick={() => onSubmit(color)}>Submit</button>
		</div>
  </div>;

}

function EditOverlay({
	color,
	onSubmit=()=>{},
	exit=()=>{}
}) {
	return <div className="edit-overlay">
		<button onClick={exit} className="exit-button">X</button>
		<Customizer color={color} onSubmit={onSubmit} />
	</div>
}

function ColorButton({
	c = Color.fromArr([0, 0, 0]),
	onDelete=()=>{},
	onEditSelect=()=>{},
	onUp=()=>{},
	onDown=()=>{},
	className="",
}) {
	const style = {
		background: c.hex(),
		color: textColor(c.rgb()),
	};
	return <div
		draggable={true}
		className={`${className} color-button`}
		style={style}
	>
		<button style={style} className="control-button" onClick={onDelete}>X</button>
		<button style={style} className="control-button" onClick={onEditSelect}>Edit</button>
		<button style={style} onClick={onUp} className="control-button">^</button>
		<button style={style} onClick={onDown} className="control-button">v</button>
		<span style={{marginLeft: 30}}>{c.hex()}</span>
	</div>
}

export function SequenceBuilder() {
	const [seq, setSeq] = useState(DEFAULT_SEQUENCE);
	const [loading, setLoading] = useState(true);

	useEffect(() => {
		async function doit() {
			const data = await getStatus();
			if (data.type === "SEQUENCE") {
				setSeq(data.data.colors.map(Color.fromArr));
			} 
			setLoading(false);
		}
		doit();
	}, []);

	const [editState, setEditState] = useState({
		editing: false, target: null,
	});

	function swapTwo(i, j) {
		const ci = seq[i];
		const cj = seq[j];
		const newSeq = seq.map((curr, ind) => {
			if (ind === i) {
				return cj;
			} else if (ind === j) {
				return ci;
			}
			return curr;
		});
		setSeq(newSeq);
	}

	function setOne(c, i) {
		const newSeq = seq.map((d, j) => {
			if (i === j) {
				return c;
			}
			return d;
		});
		setSeq(newSeq);
	}

	function drop(i) {
		const newSeq = seq.filter((_, j) => j !== i);
		setSeq(newSeq);
	}

	function prepend(c) {
		setSeq([c, ...seq]);
	}

	function submit() {
		const s = Sequence.fromArrs(seq.map(c => c.rgb()));
		setStatus(s.getPattern());
	}

	return <div>
		{editState.editing && <EditOverlay 
			color={seq[editState.target]}
			onSubmit={(c) => {setOne(c, editState.target); setEditState({editing: false});}}
			exit={() => setEditState({editing: false, target: null})}
		/>}
		<h1>Sequence</h1>
		<button onClick={() => submit(seq)}>Submit</button>
		<button onClick={() => prepend(Color.fromArr([0x50, 0x50, 0x50]))}>+</button>
		{loading ? "" : 
			seq.map((c, i) => {
				return <div className="colorbutton-row" key={i}>
					<ColorButton 
						c={c} 
						onDelete={() => drop(i)}
						onEditSelect={() => setEditState({editing: true, target: i})}
						onUp={() => {
							if (i === 0) {
								return;
							}
							swapTwo(i, i-1);
						}}
						onDown={() => {
							if (i + 1 >= seq.length) {
								return;
							}
							swapTwo(i, i+1);
						}}
					/>
				</div>;
			})
		}
	</div>;
}
