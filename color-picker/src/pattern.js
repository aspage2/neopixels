
class Solid {
	constructor(color) {
		this.color = color;
	}
	getPattern() {
		return {
			type: "SOLID",
			data: {color: this.color}
		};
	}
}

class Sequence {
	constructor(colors) {
		this.colors = colors;
	}
	getPattern() {
		return {
			type: "SEQUENCE",
			data: {
				colors: this.colors,
			}
		};
	}
}

