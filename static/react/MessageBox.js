'use strict';

/*
The following 2 functions render inner parts
of a message.
*/

function MessageLeft(props) {
    return (
        <div className="message-left">
            <div className="message-avatar">{props.avatarChar}</div>
        </div>
    );
}

function MessageRight(props) {
    return (
        <div className="message-right">
            <div className="message-header">
                <div className="message-name">{props.messageSenderName}</div>
                <div className="message-time">{props.messageTime}</div>
            </div>
            <div className="message-content">
                {props.messageContent}
            </div>
        </div>
    );
}

/*
Renders a single message.
*/
function Message(props) {
    return (
        <div className={"message " + props.parity}>
            <MessageLeft avatarChar={props.avatarChar} />
            <MessageRight messageSenderName={props.messageSenderName}
                messageTime={props.messageTime}
                messageContent={props.messageContent} />
        </div>
    );
}

/*
Renders the MessageBox.
*/
function MessageBoxRender(props) {
    let rows = props.messages.map((message, idx) => {
        let parity = idx % 2 == 0 ? 'even' : 'odd';
        return (<Message key={idx}
            avatarChar={message.avatarChar}
            messageSenderName={message.messageSenderName}
            messageTime={message.messageTime}
            messageContent={message.messageContent}
            parity={parity} />
        );
    });
    return (
        <div className="message-box"> {rows} </div>
    );
}

/*
MessageBox class.

Keeps an updated list of the last 10 calculations. Renders the
MessageBox.
*/
class MessageBox extends React.Component{
    constructor(props){
        super(props);
        this.state = {
            messages: []
        };
    }
    render(){
        return (
            <MessageBoxRender messages={this.state.messages} />
        );
    }
    /*
    Calls getMessages() every half second.
    */ 
    componentDidMount(){
        this.messageRefresher = setInterval( () => this.getMessages(), 500 );
    }
    componentWillUnmount(){
        clearInterval( this.messageRefresher );
    }
    /*
    Gets the last 10 calculations from the server.
    */
    getMessages(){
        fetch( "/backend/getLast10Calculations" )
            .then( response => response.json() )
            .then( result => {
                if( result.error ){
                    console.log( "Error getting calculations" );
                }else{
                    let messages = result.calculations.map( calculationToMessage );
                    this.setState({
                        messages: messages
                    });
                }
            }).catch( reason => {
                console.log( reason );
            });
    }
}

/*
Transforms a calculation from the database into a form suitable
for displaying.

The input, 'calculation', has the form
    calculation[ 'Username' ] = (string) the username associated with the calculation
    calculation[ 'Date' ] = (string) MySQL datetime string for when the calculation was made
    calculation[ 'X' ] = (string) the first number in the calculation
    calculation[ 'Op' ] = (string) the operation in the calculation
    calculation[ 'Y' ] = (string) the second number in the calculation
    calculation[ 'Val' ] = (string) the output of the calculation

The output is an object containing the following
    messageSenderName (string) the username associated with the calculation
    avatarChar (string) the first letter of the username
    messageTime (string) string representation of the time of day the calculation was made
    messageContent (string) string displaying the calculation

*/
function calculationToMessage( calculation ){
    let username = calculation[ 'Username' ];
    let time = new Date( calculation['Date'].replace(/-/g, '/') );
    let content = calculation[ 'X' ] + ' ' + calculation[ 'Op'] + ' ' + calculation[ 'Y']
        + " = " + calculation[ 'Val' ];
    let avatarChar = username[0].toUpperCase();
    let hours = time.getHours() % 12;
    if( hours == 0 ){
        hours = 12;
    }
    let minutes = String(time.getMinutes()).padStart(2, "0");
    let partOfDay = time.getHours() >= 12 ? "PM" : "AM";
    time = String( hours ) + ":" + String( minutes ) + " " + partOfDay; 
    return {
        messageSenderName: username,
        avatarChar: avatarChar,
        messageTime: time,
        messageContent: content
    };
}

export { MessageBox };
