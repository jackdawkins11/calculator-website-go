'use strict';

/*
Renders a div containing the input element
for the username
*/

function SignInUsername(props) {
    return (
        <div className="auth-content">
            <div className="auth-content-label">
                Username
            </div>
            <input className="auth-content-input" onChange={props.onChange} ></input>
        </div>
    );
}

/*
Renders a div containing the input element
for the password
*/

function SignInPassword(props) {
    return (
        <div className="auth-content">
            <div className="auth-content-label">
                Password
            </div>
            <input type="password" className="auth-content-input" onChange={props.onChange} ></input>
        </div>
    );
}

/*
Renders a div containing a message
*/

function SignInMessage(props) {
    return (
        <div className="auth-content">
            <div className="auth-content-message message-red">
                {props.message}
            </div>
        </div>
    );
}

/*
Renders a div containing the button for signing in.
*/

function SignInButton(props) {
    return (
        <div className="auth-content">
            <div className="auth-content-button" onClick={props.onClick} >
                Sign In
            </div>
        </div>
    );
}

/*
The SignIn component. Renders a div containing a form for signing in.

Handles signing in logic. When a sign in is successfull, calls props.signIn()
to change to the HomePage component.
*/

class SignIn extends React.Component {
    constructor(props) {
        super(props);
        this.state = ({
            username: "",
            password: "",
            message: ""
        });
    }
    render() {
        return (
            <div className="auth-box-main">
                <SignInUsername onChange={(event) => this.changeUsername(event)} />
                <SignInPassword onChange={(event) => this.changePassword(event)} />
                <SignInMessage message={this.state.message} />
                <SignInButton onClick={() => this.signIn()} />
            </div>
        );
    }
    /*
    Updates this.state.username
    */
    changeUsername(event) {
        this.setState({
            username: event.target.value
        });
    }
    /*
    Updates this.state.password
    */
    changePassword(event) {
        this.setState({
            password: event.target.value
        });
    }

    /*
    Tries to start a session with the given credentials. On failure,
    displays a message. On success, calls props.signIn() to switch to
    the home page.
    */
    signIn() {
        let username = this.state.username;
        let password = this.state.password;
        let credentials = new URLSearchParams();
        credentials.append("username", username);
        credentials.append("password", password);
        fetch("/backend/StartSession", {
            method: "POST",
            body: credentials
        }).then((response) => response.json())
            .then(result => {
                if (!result.error) {
                    if (result.hasSession) {
                        this.props.signIn();
                    } else {
                        this.invalidCredentials();
                    }
                } else {
                    this.errorMessage();
                }
            }).catch((reason) => {
                this.errorMessage();
            });
    }
    /*
    Displays a message.
    */
    invalidCredentials() {
        this.setState({
            message: "Invalid username and password."
        });
    }
    /*
    Displays a message.
    */
    errorMessage() {
        this.setState({
            message: "There was an error signing in."
        });
    }
}

export { SignIn };