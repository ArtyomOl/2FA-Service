# 2FA Service on Go

This is a microservice designed for two-factor authentication using time-based one-time password (TOTP)

### API:
POST:  
localhost:8080/api/check {"login": _ ,"code": _ } - to check the correctness of the time-based code  
localhost:8080/api/users/get {"login": _ ,"password": _ } - to verify the password is correct ang get the gain to generate TOTP  
localhost:8080/api/users/add {"login": _ ,"password": _ } - to add new user  

### Technologies:
`language` ~ `Golang`  
`database` ~ `MySQL`