import "./PresetButton.css";
import { textColor } from "../../colorutils";
import {Solid, Sequence} from "../../pattern";

export function PresetButton({msg, setting, setPattern=()=>{}, className }) {
	let props = {
		className: "preset-button",
		style: {}
	};
	if (className !== undefined)
		props.className += ` ${className}`;

	if (setting instanceof Sequence) {
		props.className += " pattern-sequence";
		props.style.color = "white";
		let hexes = setting.colors.map(c => c.hex());
		hexes = hexes.concat([hexes[0]]);
		props.style["--background-pattern"] = `conic-gradient(from 0deg, ${hexes})`;
	} else if (setting instanceof Solid) {
		props.className += " pattern-solid";
		props.style.background = setting.color.hex();
		props.style.color = textColor(setting.color.rgb());
	}
	return <button 
		onClick={() => setPattern(setting.getPattern())}
	  {...props}
	>{msg}</button>
}
