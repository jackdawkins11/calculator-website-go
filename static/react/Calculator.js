'use strict';

function CalculatorScreen(props) {
  return (
    <div className="calculator-top">
      {props.screenValue}
    </div>
  );
}

function CalculatorButton(props) {
  let classNames = "calculator-button " + props.classNames;
  return (
    <div className={classNames}
      onClick={props.onClick} >
      {props.calculatorButtonValue}
    </div>
  );
}

/* Buttons on calculator. Stored as [ column1, column2, column3, column4 ] */
let calculatorButtonValues =
  [['+', '7', '4', '1', '0'],
  ['-', '8', '5', '2', '.'],
  ['*', '9', '6', '3', 'c'],
  ['/', '=']];


/*
Renders the calculator.
*/

function CalculatorRender(props) {
  let columns = [];
  for (let c = 0; c < 4; c++) {
    let column = [];
    for (let r = 0; c < 3 && r < 5 || c ===3 && r < 2; r++) {
      let button = calculatorButtonValues[c][r];
      let classNames = "";
      if( c === 3 && r == 1 ){
        classNames +=  "equal-button";
      }
      if( r === 0 ){
        classNames += " operator-button";
      }
      column.push(
        <CalculatorButton
          key={button}
          classNames={classNames}
          calculatorButtonValue={button}
          onClick={() => props.buttonClick(button)} />
      );
    }
    columns.push(<div className="calculator-column" key={c}>{column}</div>);
  }
  return (
    <div className="calculator-area">
      <div className="calculator">
        <CalculatorScreen screenValue={props.screenValue} />
        <div className="calculator-bottom">
          {columns}
        </div>
      </div>
    </div>
  );
}

/*
Calculator class.
It keeps track of its state with the 'inputs' array.

len( inputs ) == 0  ==>  Waiting for first number
len( inputs ) == 1  ==>  In the process of entering the first number
len( inputs ) == 2  ==>  The first number and the operation have been entered
len( inputs ) == 3  ==>  Entering the second number

Digits and '.' can be entered in all 4 states. Ops can only
be entered in state 1. '=' can only be entered in state 3. 'c'
can be entered anytime.

*/

class Calculator extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      screen: 0
    };
    this.inputs = [];
    this.MAXSCREENDIGITS = 10;
  }

  /*
  Takes as input the character entered and updates state
  */
  buttonPushed(button) {
    let inputStep = this.inputs.length;
    if (isDigit(button)) {
      if (inputStep === 0) {
        this.inputs.push(button);
        this.setState({ screen: this.inputs[0] });
      } else if (inputStep === 1) {
        this.inputs[0] += button;
        this.setState({ screen: this.inputs[0] });
      } else if (inputStep === 2) {
        this.inputs.push(button);
        this.setState({ screen: this.inputs[2] });
      } else if (inputStep === 3) {
        this.inputs[2] += button;
        this.setState({ screen: this.inputs[2] });
      }
    } else if (isDecPoint(button)) {
      if (inputStep === 0) {
        this.inputs.push(button);
        this.setState({ screen: this.inputs[0] });
      } else if (inputStep === 1 && !this.inputs[0].includes('.')) {
        this.inputs[0] += button;
        this.setState({ screen: this.inputs[0] });
      } else if (inputStep === 2) {
        this.inputs.push(button);
        this.setState({ screen: this.inputs[2] });
      } else if (inputStep === 3 && !this.inputs[2].includes('.')) {
        this.inputs[2] += button;
        this.setState({ screen: this.inputs[2] });
      }
    } else if (isOp(button)) {
      if (inputStep === 1) {
        this.inputs.push(button);
      }
    } else if (isEqu(button)) {
      if (inputStep === 3) {
        let calculation = doCalculation(this.inputs).toPrecision(this.MAXSCREENDIGITS);
        calculation = parseFloat(calculation);
        this.postCalculation(this.inputs, calculation);
        this.inputs = [];
        this.setState({ screen: calculation });
      }
    } else if (isClear(button)) {
      this.inputs = [];
      this.setState({ screen: 0 });
    }
  }

  /*
  Sends a calculation to the server.
  */
  postCalculation(inputs, output) {
    let x = String(inputs[0]),
      op = inputs[1],
      y = String(inputs[2]),
      val = String(output),
      date = getDateTimeString();
    let calculation = new URLSearchParams();
    calculation.append("x", x);
    calculation.append("op", op);
    calculation.append("y", y);
    calculation.append("val", val);
    calculation.append("date", date);
    fetch("/backend/AddCalculation", {
      method: "POST",
      body: calculation
    }).then((response) => response.json())
      .then(result => {
        if (!result.error) {
          //do something
        } else {
          this.error("There was an error");
        }
      }).catch((reason) => {
        this.error("There was an error " + reason);
      });
  }
  error(message) {
    console.log(message);
  }

  render() {
    return <CalculatorRender
      screenValue={this.state.screen}
      buttonClick={(c) => this.buttonPushed(c)} />
  }

}

function isDigit(str) {
  return str.length === 1 && /[0-9]/.test(str);
}

function isDecPoint(str) {
  return str === '.';
}

function isOp(str) {
  return str.length === 1 && "+-*/".includes(str);
}

function isEqu(str) {
  return str === '=';
}

function isClear(str) {
  return str === 'c';
}

/*
Does the calculation on the given inputs array.
inputs is of the form [ number, operation, number ]
*/

function doCalculation(inputs) {
  let op = inputs[1];
  let x = parseFloat(inputs[0]);
  let y = parseFloat(inputs[2]);
  if (op === '+') {
    return x + y;
  } else if (op === '-') {
    return x - y;
  } else if (op === '*') {
    return x * y;
  } else if (op === '/') {
    return x / y;
  }
}

/*
Generates a string representation of the current time
in a MySQL datetime formate.
*/
function getDateTimeString() {
  let date = new Date();
  return date.getFullYear()
    + "-" + String(date.getMonth() + 1).padStart(2, "0")
    + "-" + String(date.getDate()).padStart(2, "0")
    + " " + String(date.getHours()).padStart(2, "0")
    + ":" + String(date.getMinutes()).padStart(2, "0")
    + ":" + String(date.getSeconds()).padStart(2, "0");
}

export { Calculator };