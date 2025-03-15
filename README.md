# 2FA Service on Go

This is a microservice designed for two-factor authentication using time-based one-time password (TOTP)

### API:
POST:  
localhost:8080/api/check --data {"login": \<string\> ,"code": \<int\> } - to check the correctness of the time-based code  
Response: {"status": bool, "message": string}  
  
localhost:8080/api/users/get --data {"login": \<string\> ,"password": \<string\> } - to verify the password is correct ang get the gain to generate TOTP  
Response: {"status": bool, "code": string}  
  
localhost:8080/api/users/add --data {"login": \<string\> ,"password": \<int\> } - to add new user  

### Operating principle:
Using the API, you can check the current totp correctness, knowing the user's login.  
This feature is blocked after three unsuccessful attempts, and you'll have to wait 30 seconds to continue.  
  
The TOTP is generated based on the current time and a secret code that is generated randomly for each user when it is added to the database.  
  
Usernames, passwords, and secret codes of users are stored encrypted in the postgresql database.  
Usernames and passwords stored as a hash functions obtained by the sha256 encryption algoritm.  
The secret code is stored in base32 format.


### Technologies:
`language` ~ `Golang`  
`database` ~ `PostgreSQL`, `Redis`  