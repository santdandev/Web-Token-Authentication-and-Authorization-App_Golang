# Sample Web Token authentication and authorization App

## This demos the use of sqlite3 db for signup, signin along with web token for authentication and authorization.

## To get dependencies

go get 

build
go build

run
./jwt

usage

POST http://localhost:80/signup

{
  "username": "xyz",
  "password": "xyz"
}


POST http://localhost:80/signin

{
  "username": "xyz",
  "password": "xyz"
}

GET http://localhost:80/checktoken

Sample Cookie :
"jwt":
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InNpZGRodSIsImV4cCI6MTYwMDI3ODU2NH0.Xm_exWEmGxX06K5XLltlZnl6-LCHlwtsD_h9cTtUEZI"