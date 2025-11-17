
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

  return [ h*360, s, v ];
}

function hsvToRgb(h, s, v) {
  let r, g, b;

  h = h % 360; // Ensure hue is within 0-360
  const c = v * s;            // Chroma
  const x = c * (1 - Math.abs(((h / 60) % 2) - 1));
  const m = v - c;

  if (h < 60) {
    r = c; g = x; b = 0;
  } else if (h < 120) {
    r = x; g = c; b = 0;
  } else if (h < 180) {
    r = 0; g = c; b = x;
  } else if (h < 240) {
    r = 0; g = x; b = c;
  } else if (h < 300) {
    r = x; g = 0; b = c;
  } else {
    r = c; g = 0; b = x;
  }

  // Convert to 0-255 range
  r = Math.round((r + m) * 255);
  g = Math.round((g + m) * 255);
  b = Math.round((b + m) * 255);

  return [ r, g, b ];
}

function hex2RGB(colorString) {
  return [
    parseInt(colorString.substring(1, 3), 16),
    parseInt(colorString.substring(3, 5), 16),
    parseInt(colorString.substring(5, 7), 16),
  ]
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

function adjustBrightness(rgb, by) {
    const colorHsv = rgbToHsv(...rgb);
    let [h, s, l] = colorHsv;
    let newL = l + by;
    if (newL > 1) {
      newL = 1.
    } else if (newL < 0) {
      newL = 0. 
    }
    return hsvToRgb(h, s, newL);
}


export class Color {
	constructor({r=0, g=0, b=0}) {
		this.r = r;
		this.g = g;
		this.b = b;
	}

	setRGB(rgb) {
		let newColor = {r: this.r, g: this.g, b: this.b};
		if (rgb.r !== undefined)
			newColor.r = rgb.r;
		if (rgb.g !== undefined)
			newColor.g = rgb.g;
		if (rgb.b !== undefined)
			newColor.b = rgb.b;
		return new Color(rgb);
	}

	setHSV(hsv) {
		let [h,s,v] = rgbToHsv(this.r, this.g, this.b);
		let newColor = hsvToRgb(
			hsv.h !== undefined ? hsv.h : h,
			hsv.s !== undefined ? hsv.s : s,
			hsv.v !== undefined ? hsv.v : v,
		);
		return Color.fromArr(newColor);
	}
	setHex(hex) {
		let [r, g, b] = hex2RGB(hex);
		return new Color({r, g, b});
	}
	hsv() {
		return rgbToHsv(this.r, this.g, this.b);
	}
	rgb() {
		return [this.r, this.g, this.b];
	}
	css() {
		return `rgb(${this.r}, ${this.g}, ${this.b})`
	}
	hex() {
		return rgb2Hex(this.rgb());
	}
	adjustBrightness(by) {
		let [, , v] = this.hsv();
		v += by;
		if (v > 1) v = 1;
		else if (v < 0) v = 0;
		return this.setHSV({v});
	}
	static fromArr([r, g, b]) {
		return new Color({r, g, b});
	}
}

export {textColor, rgb2Hex, rgbToHsv, hex2RGB, hsvToRgb, adjustBrightness}
