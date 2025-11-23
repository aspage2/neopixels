import { Color } from "./colorutils";

export class Solid {
	constructor(color) {
		this.color = color;
	}
	getPattern() {
		return {
			type: "SOLID",
			data: {color: this.color.rgb()}
		};
	}

	static fromArr(rgb) {
		return new Solid(Color.fromArr(rgb));
	}
}

export class Sequence {
	constructor(colors) {
		this.colors = colors;
	}
	getPattern() {
		return {
			type: "SEQUENCE",
			data: {
				colors: this.colors.map(c => c.rgb()),
			}
		};
	}
	static fromArrs(rgbs) {
		return new Sequence(rgbs.map(Color.fromArr));
	}
}

