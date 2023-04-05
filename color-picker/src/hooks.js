
import {useState, useRef} from "react"

function useLifeSpan(t, zeroVal) {
    const [v, _setV] = useState(zeroVal);
    const timeoutRef = useRef(null);

    const set = val => {
        if (timeoutRef.current !== null) {
            clearTimeout(timeoutRef.current);
        }
        _setV(val)
        timeoutRef.current = setTimeout(() => {
            _setV(zeroVal);
            timeoutRef.current = null;
        }, t);
    };

    return [v, set]
}

export {useLifeSpan}