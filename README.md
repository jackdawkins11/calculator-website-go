# Calculator app
A webpage with a calculator and a list of recent calculations.

## Using it
The site is running live here: http://52.15.174.58/frontend
Create an account, sign in, and then use the calculator and see recent calculations made by any of the users.

## Host it yourself
The backend uses Golang and MySQL. The frontend uses React with JSX. To set up:
* Install golang. See https://golang.org/doc/install
* Install two golang packages: `go get -u github.com/go-sql-driver/mysql && go get github.com/gorilla/sessions`
* Clone this repository to your server
* Set up the database. Create a MySQL user and put its username and password into the `main.go` file in this repository. Run the install.sql script
* Set up a JSX preprocessor. You can do so by following the steps at the bottom of this page https://reactjs.org/docs/add-react-to-a-website.html
* Create a folder called `react_compiled` inside the `static` folder. Run the JSX preprocessor on all of the JavaScript files in the `static/react` folder and put the output into the `react_compiled` folder. Afterwards your directory structure should look like
```
calculator-website-go / main.go
calculator-website-go / CheckSession.go
.
.
calculator-website-go / static / react / app.js
calculator-website-go / static / react / AuthPage.js
.
.
calculator-website-go / static / react_compiled / app.js
calculator-website-go / static / react_compiled / AuthPage.js
.
.
```
* Run the command `go run *go`
