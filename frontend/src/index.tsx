/* @refresh reload */
import { render } from 'solid-js/web';
import App from './App';
import "./styles/theme.css";
import "./styles/layout.css";
import "./styles/globals.css";
import "./styles/ark.css";

render(() => <App />, document.getElementById('root') as HTMLElement);
