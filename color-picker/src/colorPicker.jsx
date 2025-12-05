import "./colorPicker.css";
import { useEffect, useRef, useState } from "react";
import { Color, hsvToRgb, textColor, hex2RGB } from "./colorutils";

function canvasDraw(canvas, hue, width, height) {
    const ctx = canvas.getContext("2d", {willReadFrequently: true});

    const imgData = ctx.createImageData(width, height);
    const data = imgData.data;

    for (let y = 0; y < height; y++) {
      for (let x = 0; x < width; x++) {
        const s = x / width;
        const v = 1 - y / height;
        const [r, g, b] = hsvToRgb(hue, s, v);

        const i = (y * width + x) * 4;
        data[i] = r;
        data[i + 1] = g;
        data[i + 2] = b;
        data[i + 3] = 255;
      }
    }
    ctx.putImageData(imgData, 0, 0);
}

export function SVSlider({ 
	hue = 200,
	width = 300,
	height = 200,
	sv = {s: 1, v: 1},
	onChange = () => {},
	onRelease = () => {},
}) {
  const canvasRef = useRef(null);
	function onCommit(e) {
		onRelease({
			s: e.coords.x / width, 
			v: (height-e.coords.y) / height
		});
  }
	function oc(e) {
		onChange({
			s: e.coords.x / width, 
			v: (height-e.coords.y) / height
		});
	}
	const [dragState, mouseDown, touchStart] = useCanvasDrag(canvasRef, oc, onCommit);
	
  // Draw the HSV rectangle
  useEffect(() => {
	  canvasDraw(canvasRef.current, hue, width, height);
  }, [hue, width, height]);


  let diameter = 30;
	let offset = {x: -15, y: -15}; 
	if (dragState.type === "touch") {
		diameter = 60;
		offset.x = -30;
		offset.y = -100;
	}
  return (
    <div className="svslider">
      <canvas
        ref={canvasRef}
        width={width}
        height={height}
        style={{cursor: dragState.type !== null ? "none" : "default"}}
				onMouseDown={mouseDown}
				onTouchStart={touchStart}
      />

      {dragState.type !== null && (
        <div
					className="preview"
          style={{
            top: dragState.coords.y + offset.y,
            left: dragState.coords.x + offset.x,
            width: diameter,
            height: diameter,
            background: `rgb(${dragState.sample[0]}, ${dragState.sample[1]}, ${dragState.sample[2]})`,
          }}
        />
      )}
			<div
				className="selection"
				style={{
					top: height - sv.v * height - 5,
					left: sv.s * width - 5,
				}}
			/>
		</div>
  );
}

export function HueSlider({
	hue=0,
	sat=1,
	val=1,
	width=300,
	height=30,
	onHSVChange=()=>{},
	onHSVCommit=()=>{},
}) {
	const canvasRef = useRef(null);
	const [dragState, mouseDown, touchStart] = useCanvasDrag(
		canvasRef, 
		(e) => {
			const h = 360 * (e.coords.x / width)
			onHSVChange({h});
		},
		(e) => {
			const h = 360 * (e.coords.x / width)
			onHSVCommit({h});
		}
	);

	useEffect(() => {
		const canv = canvasRef.current
		const ctx = canv.getContext("2d", {willReadFrequently: true});

    const imgData = ctx.createImageData(width, height);
    const data = imgData.data;

		for (let x = 0; x < width; x++) {
			const hue = 360 * x / width;
			const [r, g, b] = hsvToRgb(hue, sat, val);
			for (let y = 0; y < height; y++) {
				const idx = (y * width + x) * 4  
				data[idx] = r;
				data[idx + 1] = g;
				data[idx + 2] = b;
				data[idx + 3] = 255;
			}
		}
    ctx.putImageData(imgData, 0, 0);

	}, [sat, val, width, height]);

	return <div className="hueslider">
		<canvas 
			ref={canvasRef}
			width={width}
			height={height}
		style={{borderRadius: 8}}
			onMouseDown={mouseDown}
			onTouchStart={touchStart}
		/>
		{dragState.type !== null && <div
			className="preview"
			style={{
				top: -30,
				left: dragState.coords.x-15,
				background: `rgb(${dragState.sample[0]}, ${dragState.sample[1]}, ${dragState.sample[2]})`,

			}}
		/>}
		<div
			className="selection"
			style={{
				top: 0,
				boxSizing: "border-box",
				height: height,
				left: width * (hue / 360),
			}}
		/>
	</div>;
}

function clampCanvasCoordinates(canvas, x, y) {
	const rect = canvas.getBoundingClientRect()
	return [
		clamp(x - rect.left, 0, canvas.width-1),
		clamp(y - rect.top, 0, canvas.height-1),
	];
}

function useCanvasDrag(canvasRef, onChange, onCommit) {
	const [dragState, setDragState] = useState({
		type: null,
		coords: {x: 0, y: 0},
		sample: [0, 0, 0],
	});

	function setDragging(e, typ) {
		const canvas = canvasRef.current;
		const [x, y] = clampCanvasCoordinates(
			canvas, e.clientX, e.clientY,
		);
		const [r, g, b] = canvas.getContext("2d", {willReadFrequently: true}).getImageData(x, y, 1, 1).data;
		const st = {
			type: typ,
			coords: {x, y}, 
			sample: [r, g, b],
		};
		onChange(st);
		setDragState(st);
	}
	
	function touchStart(e) {
		setDragging(e.touches[0], "touch");
		function touchMove(e) {
			setDragging(e.touches[0], "touch");
		}
		function touchEnd(e) {
			const canvas = canvasRef.current;
			const [x, y] = clampCanvasCoordinates(
				canvas, e.changedTouches[0].clientX, e.changedTouches[0].clientY,
			);
			const [r, g, b] = canvas.getContext("2d", {willReadFrequently: true}).getImageData(x, y, 1, 1).data;
			var st = {
				type: "touch",
				coords: {x, y},
				sample: [r, g, b],
			};
			onCommit(st);
			st.type = null;
			setDragState(st);
			document.removeEventListener("touchmove", touchMove);
		}
		document.addEventListener("touchmove", touchMove);
		document.addEventListener("touchend", touchEnd, {once: true});
	}

	function mouseDown(e) {
		setDragging(e, "mouse");
		function mouseMove(e) {
			setDragging(e, "mouse");
		}
		function mouseUp(e) {
			const canvas = canvasRef.current;
			const [x, y] = clampCanvasCoordinates(
				canvas, e.clientX, e.clientY,
			);
			const [r, g, b] = canvas.getContext("2d", {willReadFrequently: true}).getImageData(x, y, 1, 1).data;
			var st = {
				type: "mouse",
				coords: {x, y},
				sample: [r, g, b],
			};
			onCommit(st);
			st.type = null;
			setDragState(st);
			document.removeEventListener("mousemove", mouseMove);
		}
		document.addEventListener("mousemove", mouseMove);
		document.addEventListener("mouseup", mouseUp, {once: true});
	}
	return [dragState, mouseDown, touchStart];
}

function clamp(n, min, max) {
	return Math.min(Math.max(n, min), max);
}

export function ColorPreview({color, setColor=()=>{}}) {
  const hexColor = color.hex();
	const [value, setValue] = useState(hexColor);
  const _tc = textColor(color.rgb());

	useEffect(() => setValue(hexColor), [hexColor]);

  return <div style={{
    background: hexColor || "black",
    width: "100%",
    height: "80px",
    textAlign: "center"
  }}>
			<input 
				style={{
					lineHeight: "80px",
					textAlign: "center",
					color: _tc,
					fontWeight: "bold",
					fontSize: "1.1rem",
					background: "none",
					border: "none",
					outline: "none",
				}}
				type="text" 
				value={value}
				onChange={e => {
					const v = e.target.value;
					if (/^#[0-9A-Fa-f]{6}$/.test(v)){
						const x = hex2RGB(v);
						setColor(Color.fromArr(x));
					}
					setValue(v);
				}}
				onBlur={() => {
					setValue(hexColor);
				}}
			/>
  </div>
}

