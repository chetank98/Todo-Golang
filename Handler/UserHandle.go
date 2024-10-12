package Handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo/Database/dbHelper"
	"todo/Middleware"
	"todo/Models"
	"todo/Utils"
)

var Access bool

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating the user")

	var user Models.Users

	if r.Body == nil {
		Utils.RespondError(w, http.StatusBadRequest, nil, "enter sufficient data")
	}
	decErr := json.NewDecoder(r.Body).Decode(&user)
	if decErr != nil {
		Utils.RespondError(w, http.StatusBadRequest, decErr, "failed to decode the data")
	}

	already, alreadyErr := dbHelper.AlreadyUser(user.Email)
	if alreadyErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, nil, "Failed to check user identity")
	}

	if already {
		Utils.RespondError(w, http.StatusBadRequest, nil, "User already in database")
	}

	hashedPassword, hasErr := Utils.HashPassword(user.Password)
	if hasErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
		return
	}

	saveErr := dbHelper.RegisterUser(user.UserName, user.Email, hashedPassword)
	if saveErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to save user")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"user created successfully"})

}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	var login Models.UserLogin

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		Utils.RespondError(w, http.StatusBadRequest, err, "invalid payload")
	}

	if r.Body == nil {
		Utils.RespondError(w, http.StatusBadRequest, nil, "enter the required the data")
	}

	userID, name, Email, checkErr := dbHelper.LoginCheck(login.Email, login.Password)
	if checkErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, checkErr, "failed to find user")
		return
	}

	sessionID, crtErr := dbHelper.SessionGenerated(userID)
	if crtErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, crtErr, "failed to create user session")
		return
	}

	token, genErr := Utils.GenerateJWT(userID, Email, name, sessionID)
	if genErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, genErr, "failed to generate token")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{"user logged in successfully", token})

}

func Logout(w http.ResponseWriter, r *http.Request) {
	userCtx := Middleware.UserContext(r)
	sessionID := userCtx.SessionID

	saveErr := dbHelper.DeleteUserSession(sessionID)
	if saveErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to delete user session")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"user logged out successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Deleting the user")

}
