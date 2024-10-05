package Handle

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"todo/Database"
	"todo/Models"
	"todo/Utils"
)

/*
	1) Create separate function for parseBody, EncodeJSON, DecodeJSON, ResPondJSON and call it
	2) send errors to frontend in a structured manner
	3) use proper status code and send in responses
*/

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// No need of this line
	fmt.Println("Creating the user")

	var user Models.Users

	if r.Body == nil {
		w.Write([]byte("Insert the information for registartion"))
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {

		// send proper error messages with status code in a JSON format

		fmt.Println("error")
	}

	Pass, hasErr := Utils.HashPassword(user.Password)
	if hasErr != nil {
		json.NewEncoder(w).Encode("enter the password")
		fmt.Println("enter the password")
	}
	user.Password = Pass

	fmt.Println(user)

	// need to make a separate function in a file and place it to the dbHelper folder
	// check first whether exists already or not then proceed further

	queryString := `INSERT INTO regisuser (username, email, password)  VALUES ($1,$2,$3)`

	fmt.Println(queryString)

	res, err := Database.DBConnection.Query(queryString, user.UserName, user.Email, user.Password)
	if err != nil {
		log.Fatal(err)
		fmt.Println("error in db execution")
		return
	}

	/*
		Its a very bad thing you are inserting the data first and then checking later whether the field is blank or not
	*/

	if user.UserName == "" || user.Email == "" || user.Password == "" {
		json.NewEncoder(w).Encode("enter the data")
	}

	json.NewEncoder(w).Encode("User Created Sucessfully")

	fmt.Println(res)

}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	fmt.Println("User Login")

	var user Models.Users

	// create a struct model for the details

	var storedEmail string
	var storedPassword string

	if r.Body == nil {
		json.NewEncoder(w).Encode("Enter the creditionals")
	}

	// need to check arcvhived_at is null or not
	// need to create session_id and JWT token here

	Query := `Select email,password from regisuser where email=$1`

	err := Database.DBConnection.QueryRow(Query, user.Email).Scan(&storedEmail, &storedPassword)
	if err != nil {
		if err != nil {
			json.NewEncoder(w).Encode("No user found with this email")
		}
		json.NewEncoder(w).Encode("error in searching")
	}

	res := Utils.CheckPassword(user.Password, storedPassword)

	if res == true && user.Email == storedEmail {

	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Deleting the user")

}
