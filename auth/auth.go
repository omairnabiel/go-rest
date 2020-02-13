package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/omairnabiel/go-lang-starter/utils"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

var goCache *cache.Cache
var sessionCache *cache.Cache

// Will get this from environment variables in real world application
var jwtKey = []byte("secret_key")

// SignUpRequest maps to signup payload user sends in request body
type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// LoginRequest maps to login api call params
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// LoginResponse object to send after successful login of user
type LoginResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// LogoutRequest maps to logout api call params
type LogoutRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Signup function to execute on signup route. Checks if user doesn't exist then add the user in the DB
func Signup(ctx *gin.Context) {
	var body SignUpRequest

	// initialize a new validator
	v := validator.New()

	// bind request body to JSON
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	// pass the mapped request object to validtor to check it's validity
	errValid := v.Struct(body)

	// if validation error occurs send error to user
	if errValid != nil {
		var errors []string
		for _, e := range errValid.(validator.ValidationErrors) {
			log.Println("Errors", e.Field(), e.Tag(), e.Param())
			errors = append(errors, utils.ValidationMessage(e.Field(), e.Tag(), e.Param()))
		}
		var status []int
		status = append(status, http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessages(status, errors))
		return
	}

	// check if user exists in Database (goCache being used as database here)
	_, found := goCache.Get(body.Email)

	// if user already exists notify the user else create the user in Database
	if found {
		ctx.JSON(http.StatusConflict, utils.ErrorMessage(http.StatusConflict, utils.ErrUserExists))
		return
	}

	encrypted, err := hashPassword(body.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorMessage(http.StatusInternalServerError, utils.ErrInternalServerError))
		return
	}
	user := body
	user.Password = string(encrypted)
	goCache.Set(body.Email, user, cache.NoExpiration)
	ctx.JSON(http.StatusOK, utils.SuccessMessage(http.StatusOK, utils.SuccessUserCreated, nil))
}

// Login function to execute on signup route. Checks if user doesn't exist then add the user in the DB
func Login(ctx *gin.Context) {
	var creds LoginRequest

	// bind request creds to JSON
	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessage(http.StatusBadRequest, err.Error()))
	}

	v := validator.New()

	errValid := v.Struct(creds)

	if errValid != nil {
		var errors []string
		for _, e := range errValid.(validator.ValidationErrors) {
			log.Println("Errors", e)
			errors = append(errors, utils.ValidationMessage(e.Field(), e.Tag(), e.Param()))
		}
		var status []int
		status = append(status, http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessages(status, errors))
		return
	}
	// check if user exists in Database
	cachedUser, found := goCache.Get(creds.Email)

	// map the cache return string varialbe to  SignUpRequest User Type

	// if user is not found in DB. Throw error and return
	if !found {
		ctx.JSON(http.StatusNotFound, utils.ErrorMessage(http.StatusNotFound, utils.ErrUserDoesntExist))
		return
	}

	user, ok := cachedUser.(SignUpRequest)

	// if mapping is fails panic
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorMessage(http.StatusInternalServerError, utils.ErrInternalServerError))
		return
	}

	// if the user exists but password doesn't match. Throw error and return
	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))

	if passErr != nil {
		ctx.JSON(http.StatusForbidden, utils.ErrorMessage(http.StatusForbidden, utils.ErrIncorrectPassword))
		return
	}

	// generate token with Email as a prop in claims
	token, err := generateToken(creds.Email)

	// if token generation fails. Throw an error and return
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorMessage(http.StatusInternalServerError, utils.ErrInternalServerError))
		return
	}

	// form a response object
	var resp interface{}
	resp = &LoginResponse{Email: creds.Email, Token: token, Name: user.Name}

	// set user token and w.r.t user's email in cache for session
	sessionCache.Set(creds.Email, token, cache.DefaultExpiration)
	ctx.JSON(http.StatusOK, utils.SuccessMessage(http.StatusOK, utils.SuccessUserLogin, resp))
}

// Logout route
func Logout(ctx *gin.Context) {
	var body LogoutRequest

	// bind request body to JSON
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessage(http.StatusBadRequest, err.Error()))
		return
	}

	// check if user is logged in and has a session running
	_, found := sessionCache.Get(body.Email)

	// if user is found in session cache logout the use else notify the user that it's already logged out
	if found {
		sessionCache.Delete(body.Email)
		ctx.JSON(http.StatusOK, utils.SuccessMessage(http.StatusOK, utils.SuccessUserLoggedOut, nil))
	} else {
		ctx.JSON(http.StatusBadRequest, utils.ErrorMessage(http.StatusBadRequest, utils.ErrAlreadyLoggedOut))
	}
}

func hashPassword(password string) ([]byte, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

func generateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email":     email,
		"ExpiresAt": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// init function sets the db and session cache gets called at beginning of the execution
func init() {

	// initialize cache with default cache retention time of 20 mins and cleanup time 20 mins
	goCache = cache.New(20*time.Minute, 20*time.Minute)
	sessionCache = cache.New(20*time.Minute, 10*time.Minute)
}
