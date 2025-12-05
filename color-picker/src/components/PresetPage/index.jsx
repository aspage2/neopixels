
import "./PresetPage.css";

import { Solid, Sequence } from "../../pattern";
import { hex2RGB } from "../../colorutils";
import { PresetButton } from "../PresetButton";
import { setStatus } from "../../api";

const PRESETS = {
//  "OFF": Solid.fromArr([0, 0, 0]),
	"\u{273F}": Solid.fromArr(hex2RGB("#586593")),
	"Mintel": Solid.fromArr(hex2RGB("#fedb00")),
	"BB <GO>": Solid.fromArr(hex2RGB("#FFA028")),
	"\u{1F485}": Sequence.fromArrs([
		[255, 0, 0],
		[252, 186, 3],
		[207, 207, 0],
		[0, 255, 0],
		[0, 0, 255],
		[85, 0, 171],
		[187, 0, 250]
	]),
	"Warm White 🎄": Solid.fromArr([244,187,113]),
	"\u{1F341}": Sequence.fromArrs([
		[218, 97, 21],
		[218, 178, 52],
		[222, 1, 84],
		[129, 83, 24]
	]),
};

export function PresetPage() {
    return <>
			<h2>Presets</h2>
			<div className="preset-page-content">
				{
					Object.entries(PRESETS).map(([colorName, pattern]) => {
						return <PresetButton 
							key={colorName} 
							setPattern={setStatus}
							setting={pattern}
							msg={colorName}
						  className={"preset-page-button"}
						/>
					}) 
				}
			</div>
		</>;
}

