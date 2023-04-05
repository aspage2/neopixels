import './App.css';
import { GetPattern } from "./pattern";
import { useLifeSpan } from './hooks';
import { getStatus, setStatus } from './util';
import {useState, useEffect} from "react";
import map from "lodash/map";
import {useCookies} from "react-cookie"


function FavList({onChoose, onRemove, favs}) {
    const favItem = (data, i) => {
        const pattern = GetPattern(data);
        return <div 
            className="fav-item" 
            onClick={()=>onChoose(i)}
        >
            <pattern.PreviewBar className="preview" data={data.data} />
            <button onClick={()=>onRemove(i)}>X</button>
        </div>
    };

    return <div>{map(favs, favItem)}</div>
}

const PATTERNS = {
    SOLID: new Solid(),
    GRADIENT: new Gradient(),
}

function App() {
    const [stat, setStat] = useState({type: "SOLID", data: {color: [0, 0, 0]}})
    const [err, setErr] = useLifeSpan(3000, "");
    useEffect(() => {
        getStatus().then(setStat)
    }, []);

    const [cookies, _setFavs] = useCookies(["fav-patterns"]);
    const setFavs = fs => _setFavs("fav-patterns", fs, {
        expires: new Date(Date.now() + 10 * 3600000 * 24 * 365),
    });

    const pattern = GetPattern(stat);
    const favs = cookies["fav-patterns"] || [];

    useEffect(() => {
        if (favs["fav-patterns"] === undefined) {
            setFavs([])
        }
    }, []);

    return <div className="App">
        <div>
            <label>
                <input type="radio" checked={stat.type === "SOLID"} onClick={() => setStat({type: "SOLID", data:{color: [0, 0, 0]}})}/>
                SOLID
            </label>
            <label>
                <input type="radio" checked={stat.type === "GRADIENT"} onClick={() => setStat({type: "GRADIENT", data:{colors:[[0,0,0], [255, 0, 0]]}})}/>
                GRADIENT
            </label>
        </div>
        <div>
            <pattern.Selection 
                data={stat.data} 
                update={data => setStat({
                    type: pattern.type,
                    data
                })}
            />
            <pattern.PreviewBar data={stat.data} className="preview"/>
        </div>
        <button onClick={() => setStatus(stat).catch(reason => setErr(reason.message))}>Submit</button>
        <button onClick={() => getStatus().then(setStat)}>Refresh</button>
        <button onClick={() => setFavs(favs.concat([stat]))}>Add Favorite</button>
        <h1>{err}</h1>
        <FavList
            favs={favs}
            onChoose={i => {
                setStat(favs[i]);
            }}
            onRemove={i => {
                let fs = Array.from(favs)
                fs.splice(i, 1);
                setFavs(fs)
            }}
            />
    </div>;
}

export default App;
