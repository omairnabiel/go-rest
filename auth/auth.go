package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
)

var goCache *cache.Cache
var sessionCache *cache.Cache

// Will get this from environment variables in real world application
var jwtKey = []byte("secret_key")

// SignUpRequest maps to signup payload user sends in request body
type SignUpRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// LoginRequest maps to login api call params
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse object to send after successful login of user
type LoginResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// LogoutRequest maps to logout api call params
type LogoutRequest struct {
	Email string `json:"email"`
}

// Signup function to execute on signup route. Checks if user doesn't exist then add the user in the DB
func Signup(ctx *gin.Context) {
	var body SignUpRequest

	// bind request body to JSON
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// check if user exists in Database (goCache being used as database here)
	_, found := goCache.Get(body.Email)

	// if user already exists notify the user else create the user in Database
	if found {
		ctx.JSON(http.StatusConflict, gin.H{"error": "User Already Exists"})
		return
	}

	encrypted, err := hashPassword(body.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed creating user. Please try again"})
		return
	}
	user := body
	user.Password = string(encrypted)
	goCache.Set(body.Email, user, cache.NoExpiration)
	ctx.JSON(http.StatusOK, gin.H{"success": "User Created Successfully"})
}

// Login function to execute on signup route. Checks if user doesn't exist then add the user in the DB
func Login(ctx *gin.Context) {
	var creds LoginRequest

	// bind request creds to JSON
	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// check if user exists in Database
	cachedUser, found := goCache.Get(creds.Email)

	// map the cache return string varialbe to  SignUpRequest User Type

	// if user is not found in DB. Throw error and return
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User doesn't exist. Please signup"})
		return
	}

	user, ok := cachedUser.(SignUpRequest)

	// if mapping fails panic
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// if the user exists but password doesn't match. Throw error and return
	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))

	if passErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Incorrect Password"})
		return
	}

	// generate token with Email as a prop in claims
	token, err := generateToken(creds.Email)

	// if token generation fails. Throw an error and return
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login. Please Try Again"})
		return
	}

	resp := &LoginResponse{Email: creds.Email, Token: token, Name: user.Name}
	sessionCache.Set(creds.Email, token, cache.DefaultExpiration)
	ctx.JSON(http.StatusOK, gin.H{"user": resp})
}

// Logout route
func Logout(ctx *gin.Context) {
	var body LogoutRequest

	// bind request body to JSON
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user is logged in and has a session running
	_, found := sessionCache.Get(body.Email)

	// if user is found in session cache logout the use else notify the user that it's already logged out
	if found {
		sessionCache.Delete(body.Email)
		ctx.JSON(http.StatusOK, gin.H{"success": "User logged out successfully"})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "You're already logged out"})
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

// InitializeCache sets the db and session cache
func InitializeCache() {

	// initialize cache with default cache retention time of 20 mins and cleanup time 20 mins
	goCache = cache.New(20*time.Minute, 20*time.Minute)
	sessionCache = cache.New(20*time.Minute, 10*time.Minute)
}
