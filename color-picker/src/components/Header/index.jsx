
import "./Header.css";
import powerIcon from "./power-switch-svgrepo-com.svg";

import { NavLink } from "react-router";
import { setStatus } from "../../api";
import { Solid } from "../../pattern";

export function Header() {

	return <div id="site-header">
			<button 
				className="off-button"
				onClick={() => setStatus(Solid.fromArr([0, 0, 0]).getPattern())}
			><img src={powerIcon}/></button>
			<nav>
				<NavLink 
					index="true"
					className={({isActive}) => isActive ? "nav-link active" : "nav-link"} 
					to="/"
				>Presets</NavLink>
				<NavLink 
					className={({isActive}) => isActive ? "nav-link active" : "nav-link"} 
					to="/builder"
				>Solid</NavLink>
				<NavLink 
					className={({isActive}) => isActive ? "nav-link active" : "nav-link"} 
					to="/sequence"
				>Sequence</NavLink>
			</nav>
	</div>;
}
