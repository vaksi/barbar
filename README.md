# Barbar Service 
Barbar Service is a project for demo of simple distribute authentication Service. 

## Project Structure
    Basicly we reffer to Domain Driven Design and Clean Architecture for this code.

    - cmd  // directory for handle command
        |- grpc_user.go
        |- ...
    - config // directory for handle config
        |- config.go
    - pkg // directory for internal utilities or wrapper 
        |- jwt
        |- middleware
        |- mongodb
        |- ...
    - domain // directory for internal services of domain service
        |- auth // domain auth service
            |- entity // is a directory of code for the entity domain
            |- infrastructure // directory of code for handle dependencies to other service
            |- repository // for handle data story (example the service should be insert to mongodb)
            |- usecase // business logic or domain logic handle on here 
            |- tranport // directory of code for handle transport layer example (grpc, http)
        |- users // domain user service
    - proto // is share file for proto onto the services
    - .env // is a configuration file for running local
    - .env_docker is a configuration file for running with docker
    - main.go is a root running program to call functions of the cases

## Tech Stack 
* Overral Basic Programming leaguage use golang 
* Datastore with MongoDB
* Communicaton with HTTP Rest API (for external connection) and GRPC (internal connection)
* Redis for caching model, for this case cache should be store valid Token
    
## Service Design
![alt text](doc/barbar_service.png)

## Rules 
After this program started without problem, we can hit the API for User-Service and Auth-Service.

* Register 
    - Email must be unique from other user has been registered 
    - Password will be encrypt on database
    - Register don't use Token `Authorization` 
* Update User, Delete User, Get User By ID, Get All User
    - Must use Token `Authorization`
    - Token should be expired on 1 hour
    - We refer to RESTFull API for users endpoint
    - If you have error invalid token/token expired please login again for getting new token
* Login
  - Don't use Token `Authorization`
  - Input valid email and password
  - Data should be store on auth database on mongo (for audit) and store on cache (redis)
* Logout
    - Must use Token `Authorization`
    - Data should be update the user record and also clear redis record by userId

## How to Run

### Install and Run Program
* install docker
* run docker build  
  ```shell
  $ docker build -t barbar .
  ```
* run docker compose 
  ```shell
  $ docker-compose up
  ```

### CURL 
* Register
```shell
curl --location --request POST 'localhost:3000/user-service/api/v1/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Darto Test 1",
    "phone": "081131234",
    "email": "test1@email.com",
    "password": "testing"
}'
```

* Login
```shell
curl --location --request POST 'localhost:3001/auth-service/api/v1/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "test1@email.com",
    "password": "testing"
}'
```
* Logout
```shell
curl --location --request GET 'localhost:3001/auth-service/api/v1/logout' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QxQGVtYWlsLmNvbSIsInVzZXJJZCI6ImQ2MjEyOGVlLTA3OWMtNDU5MS04MmQ0LTg3NWU0NmZkODI2ZCIsImV4cCI6MTYzNDQ5NTc5MSwianRpIjoiNTkyM2U1MmMtOWQxOS00ZjJiLWFkODctN2QwYjc2NWIzNmJmIiwic3ViIjoidXNlciJ9.jq-kYrL_Xo92hzFSH2DlntjPfU8YV9_fW2ezX3QZuUA' \
--data-raw ''
```

* Get User by Id
```shell
curl --location --request GET 'localhost:3000/user-service/api/v1/users/a80d3a5d-3452-4d6a-9ace-055bebcc9523' \
--header 'Authorization: Bearer  <token>'
```
* Get All User
  
* Update User
```shell
curl --location --request PUT 'localhost:3000/user-service/api/v1/users/a80d3a5d-3452-4d6a-9ace-055bebcc9523' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QxQGVtYWlsLmNvbSIsInVzZXJJZCI6ImQ2MjEyOGVlLTA3OWMtNDU5MS04MmQ0LTg3NWU0NmZkODI2ZCIsImV4cCI6MTYzNDQ5NTc5MSwianRpIjoiNTkyM2U1MmMtOWQxOS00ZjJiLWFkODctN2QwYjc2NWIzNmJmIiwic3ViIjoidXNlciJ9.jq-kYrL_Xo92hzFSH2DlntjPfU8YV9_fW2ezX3QZuUA' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Darto Test 1 Update",
    "phone": "081131233",
    "email": "test1@email.com",
    "password": "testing"
}'
```
* Delete User
```shell
curl --location --request DELETE 'localhost:3000/user-service/api/v1/users/e9d7353d-1c5c-4412-9a3b-62adb9bac02b' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3QyQGVtYWlsLmNvbSIsInVzZXJJZCI6ImU5ZDczNTNkLTFjNWMtNDQxMi05YTNiLTYyYWRiOWJhYzAyYiIsImV4cCI6MTYzNDUwNzc4NSwianRpIjoiYjAwMWZiNzEtMzY4OC00YWNmLTg0ZWMtYjMxODI3OWVjZjc5Iiwic3ViIjoidXNlciJ9.qW8JX9Ifnx1ztKdB3DMvbWCtcNkPcD_ZE5UE7bpQrJQ'
```

    
