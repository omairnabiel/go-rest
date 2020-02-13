## **GO HTTP SERVER**

This http api server uses go [gin](https://github.com/gin-gonic/gin) framework.

These include the user auth on-boarding routes :

>1. `/signup    POST`
>2. `/login    POST`
>3. `/logout   POST`


`/signup` route takes following body as params to add the user in cache (goCache is being used a database)
```
{
	"user" : "MyName",
	"email" : "myemail@gmail.com",
	"password" : "myPassword"
}
```

`/login` you can login once a user is created using email and password

`/logout` logout the user with the token you've recieved in the login response. Token is to inserted in Headers Authorization Bearer

`verifyToken()` middleware is added in and is for now being only used for log out.

> NOTE : This project uses goCache for Database. Any database can be used but this is just for experimental purposes

###### > Make sure go is installed in your machine

To run the project clone the repo and cd into it

Run

`$ go run main.go`