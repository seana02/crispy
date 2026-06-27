/* @refresh reload */
import { render } from 'solid-js/web';
import App from 'src/App';
import "src/styles/theme.css";
import "src/styles/layout.css";
import "src/styles/globals.css";
import "src/styles/ark.css";

render(() => <App />, document.getElementById('root') as HTMLElement);
