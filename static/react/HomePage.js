import { Calculator } from './Calculator.js';
import { MessageBox } from './MessageBox.js';

/*
Renders the home page.
*/

function HomePage(props) {
    return (
        <div className="homepage">
            <div className="homepage-top-bar">
                <div className="sign-out-button" onClick={props.signOut}>
                    Sign Out
                </div>
            </div>
            <div className="homepage-main">
                <Calculator />
                <MessageBox />
            </div>
        </div>
    );
}

export {HomePage};