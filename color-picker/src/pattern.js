import { useState, useEffect, useRef } from "react";
import isEmpty from "lodash/isEmpty";

class Solid {
    constructor() {
        this.type = "SOLID"
    }
    PreviewBar({data: {color}, ...props}) {
        let style = {backgroundColor: `rgb(${color.join(',')})`};
        if (props.hasOwnProperty("style")) {
            style = {
                ...props.style,
                ...style,
            };
        }
        return <div {...props} style={style} />
    }

    Selection({data, update}) {
        return <div>
        <ColorInput width="10" onUpdate={c => update({color: c})} value={data.color}/>
    </div>
    }
}

class Gradient {
    constructor() {
        this.type = "GRADIENT"
    }
    PreviewBar({data: {colors}, ...props}) {
        let style = {
            background: `linear-gradient(90deg, rgb(${colors[0].join(",")}), rgb(${colors[1].join(",")})`
        };
        if (props.hasOwnProperty("style")) {
            style = {
                ...props.style,
                ...style,
            };
        }
        return <div {...props} style={style} />
    }

    Selection({data: {colors}, update}) {
        return <div>
            <div><ColorInput
                width="10" 
                value={colors[0]} 
                onUpdate={c => update({colors: [c, colors[1]]})} 
            />
            <ColorInput
                width="10"
                value={colors[1]}
                onUpdate={c => update({colors: [colors[0], c]})}
            /></div>
        </div>
    }
}

class Unknown {
    PreviewBar({data, ...props}) {
        return <p>???</p>
    }

    Selection({data, ...props}) {
        return <p>???</p>
    }
}

function GetPattern({type}) {
    switch(type) {
    case "SOLID":
        return new Solid();
    case "GRADIENT":
        return new Gradient();
    default:
        return new Unknown();
    }
}

function validateRGBArray(value) {
    const arr = JSON.parse(value);
    if(!Array.isArray(arr)) {
        throw new TypeError(`${value} is not a JSON array`);
    }

    if(arr.length !== 3) {
        throw new TypeError(`${value} must have 3 elements (not ${arr.length})`)
    }
    
    if(!arr.every(x => 0 <= x && x <= 255)) {
        throw new TypeError(`${value} must have elements in the range [0, 255]`)
    }

    return arr
}

function ValidatedInput({value, validate, onUpdate, ...props}) {
    const ref = useRef();
    const [err, setErr] = useState("");
    useEffect(() => {ref.current.value = JSON.stringify(value)}, [value]);
    const noErr = isEmpty(err);
    const submit = () => {
        try {
            const validated = validate(ref.current.value);
            onUpdate(validated);
            setErr("");
        } catch(e) {
            setErr(e.message);
        }
    };
    return <div style={{display: "inline-block"}}>
        {noErr ? <></> : <p className="err-msg">{err}</p>}
        <input {...props} className={noErr ? "" : "err"} type="text" ref={ref} onBlur={submit} />
    </div>
}

function ColorInput(props) {
    return <ValidatedInput validate={validateRGBArray} {...props}/>
}

export {GetPattern};