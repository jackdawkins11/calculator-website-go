'use strict';

function SignUpUsername(props) {
    return (
        <div className="auth-content">
            <div className="auth-content-label">
                Username
            </div>
            <input className="auth-content-input" onChange={props.onChange} ></input>
        </div>
    );
}

function SignUpPassword(props) {
    return (
        <div className="auth-content">
            <div className="auth-content-label">
                Password
            </div>
            <input type="password" className="auth-content-input" onChange={props.onChange} ></input>
        </div>
    );
}

function SignUpPasswordConfirm(props) {
    return (
        <div className="auth-content">
            <div className="auth-content-label">
                Confirm Password
            </div>
            <input type="password" className="auth-content-input" onChange={props.onChange} ></input>
        </div>
    );
}

function SignUpMessage(props) {
    let classNames = "auth-content-message " + props.classNames;
    return (
        <div className="auth-content">
            <div className={classNames}>
                {props.message}
            </div>
        </div>
    );
}

function SignUpButton(props) {
    return (
        <div className="auth-content">
            <div className="auth-content-button" onClick={props.onClick} >
                Sign Up
            </div>
        </div>
    );
}

/*
SignUp component. Renders a form for signing up. Handles logic for signing up.
Displays a message when signing up fails or succeeds.
*/
class SignUp extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            username: "",
            password: "",
            confirmPassword: "",
            message: "",
            messageClassNames: ""
        };
    }
    render() {
        let passwordStengthMessage = "*Passwords must be at least 10 characters long"
            + " and contain at least 1 uppercase letter, lowercase letter, number and symbol from"
            + " !@#$%^&*()";
        return (
            <div className="sign-up">
                <SignUpUsername onChange={(event) => this.changeUsername(event)} />
                <SignUpPassword onChange={(event) => this.changePassword(event)} />
                <SignUpPasswordConfirm onChange={(event) => this.changeConfirmPassword(event)} />
                <SignUpMessage message={this.state.message}
                    classNames={this.state.messageClassNames}/>
                <SignUpButton onClick={() => this.signUp()} />
                <SignUpMessage message={passwordStengthMessage}
                    classNames="password-strength-message" />
            </div>
        );
    }
    changeUsername(event) {
        this.setState({
            username: event.target.value
        });
    }
    changePassword(event) {
        this.setState({
            password: event.target.value
        });
    }
    changeConfirmPassword(event) {
        this.setState({
            confirmPassword: event.target.value
        });
    }
    /*
    Tries to create an account with the given credentials.
    */
    signUp() {
        let username = this.state.username,
            password = this.state.password,
            confirmPassword = this.state.confirmPassword;
        if (username === "" || password === "") {
            this.invalidInput("All fields must be used");
        } else if (password !== confirmPassword) {
            this.invalidInput("Passwords must match");
        } else {
            this.signUpRequest(username, password);
        }
    }
    /*
    Requests the server to make an account with the given credentials.
    Displays a message on failure or success.
    */
    signUpRequest(username, password) {
        let credentials = new URLSearchParams();
        credentials.append("username", username);
        credentials.append("password", password);
        fetch("/backend/CreateAccount", {
            method: "POST",
            body: credentials
        }).then((response) => response.json())
            .then(result => {
                if( !result.error && result.createdAccount ){
                    this.success();
                }else if( !result.error && !result.createdAccount ){
                    this.invalidInput( result.message );
                }else{
                    this.error();
                }
            }).catch((reason) => {
                this.error();
            });
    }
    success() {
        this.setState({
            message: "Successfully created account",
            messageClassNames: "message-green"
        });
    }
    invalidInput(message) {
        this.setState({
            message: message,
            messageClassNames: "message-red"
        });
    }
    error() {
        this.setState({
            message: "There was an error signing up",
            messageClassNames: "message-red"
        });
    }
}

export { SignUp };