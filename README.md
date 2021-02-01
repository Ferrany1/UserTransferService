# UserTransferService CLI + http

Service allows to create/ get balance/ make a transaction/ update user_info/ delete user with basic auth

## HTTP handlers:
V1 route /api/v1
> /user
> * `POST /create`                          -  takes json in body (Example: {"user_data":{"username":"test1","email":"test1@test.com","password":"test1pass"},"user_balance":{"username":"test1"}})
> * `GET /create`                           -  takes json in body (Example: {"user_data":{"username":"test1","email":"test1@test.com","password":"test1pass"},"user_balance":{"username":"test1"}})
> * `POST /transfer/:receiver/:amount`      -  takes json in body (Example: {"user_data":{"username":"test1","email":"test1@test.com","password":"test1pass"},"user_balance":{"username":"test1"}})
> * `PUT /update`                           -  takes json in body (Example: {"user_data":{"username":"test1","email":"test1@test.com","password":"test1pass"},"user_balance":{"username":"test1"}})
> * `DELETE /delete`                        -  takes json in body (Example: {"user_data":{"username":"test1","email":"test1@test.com","password":"test1pass"},"user_balance":{"username":"test1"}})

## CLI commands:

> User
> * `Create`                                -  takes username, email, password as arguments                             (Example: create username:test1 email:test1@test.com password:test1pass)
> * `Balance`                               -  takes username, email, password as arguments                             (Example: balance username:test1 email:test1@test.com password:test1pass)
> * `Transfer`                              -  takes username, email, password, receiver, amount as arguments           (Example: transfer username:test1 email:test1@test.com password:test1pass receiver:test2 amount:5)
> * `Update`                                -  takes username, email, password, new_email, new_password as arguments    (Example: update username:test1 email:test1@test.com password:test1pass new_email:test1@test1.com)
> * `Delete`                                -  takes username, email, password as arguments                             (Example: delete username:test1 email:test1@test.com password:test1pass)