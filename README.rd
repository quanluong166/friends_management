# FRIENDS_MANAGEMENT
This project implements a simple backend system for handling friend management business logic of social web/application
## Prerequisites
- Docker: Make sure to install docker application in your system
- PostgreSQL: Install postgreSQL latest version
- Go: Install go lastest version

## How to run
- Start docker enviroment
- Start postgreSQL
- Open terminal at project directory and type make start-app

## Project structure
```sh
friends-management/
├── cmd/
│   └── server/ 
|   └── migrate/ 
├── internal/
│   ├── config/ 
│   ├── constant/            
│   ├── controller/ 
│   ├── db/ 
│   ├── handler/ 
│       ├── api/ 
│   ├── model/ 
│   ├── repository/ 
│   ├── routes/ 
├── pkg/
│   ├── helper/ 
│   ├── utils/ 
├── Dockerfile
├── Makefile
├── document
├── Readme.nd
├── go.mod
├── go.sum
```
1. /cmd contains the main application entry point files for the project
2. /internal contains private project code. it includes: config, constant, controller, db, handler, model, repository, routes packages
    i. /handler contains all implement functions of api interface. This layer will process the request, validate the input and call controller layer for handle business logic.
    ii. /controller contains all the business functions. This layer will handle business logic, interfact this database layer and return result to handler.
    iii. /reposiotry contains all the functions that is needed by the controller layer to interact with the database for crud actions
    iv. /model contains all the ORM models
    v. /config contains config file of the application
    vi. /db contains file to init database connection
    vii. /constants contains all constant use within the application
    viii. /routes contains all the api path of the application
3. /pkg contains code that is can use outside of internal logic
    i. /helper contains helper file
    ii. /utils contains utility function for string, array

## Database Schema

### UserRelationship Table
| Column Name      | Data Type     | Constraints                                                | Description                          |
|------------------|---------------|-------------------------------------------------------------|--------------------------------------|
| `id`             | `uint`        | Primary Key, Auto Increment                                 | Unique identifier                    |
| `requestor_email`| `varchar(255)`| Not Null                                                    | Email of the requestor               |
| `target_email`   | `varchar(255)`| Not Null                                                    | Email of the target                  |
| `type`           | `text`        | Check: 'FRIEND', 'BLOCK', 'SUBSCRIBER'                      | Type of relationship                 |
| `created_at`     | `timestamp`   | Auto-managed by GORM                                        | Record creation time                 |
| `updated_at`     | `timestamp`   | Auto-managed by GORM                                        | Last update time                     |

## APIs

1.Create friend connection: (POST /api/user/relationship/friend)
1.1 Request body
```
friends: array of two emails that need make friend connection
```
+ Example:
```
{
    "friends" : ["friend7@example.com", "friend8@example.com"]
}
```
1.2 Response body
+ Success:
```
{
    "success": true
}
```
+ invalid_one_of_two_email_input:
```
{
    "success": false,
    "message": "INVALID_EMAIL_INPUT"
}
```

+ invalid_one_of_two_email_missing:
```
{
    "success": false,
    "message": "AT_LEAST_TWO_EMAILS_ARE_REQUIRED"
}
```
+ fail_already_friend:
```
{
    "success": false,
    "message": "YOU_ALREADY_FRIEND"
}
```
+ fail_one_of_two_email_block_the_other:
```
{
    "success": false,
    "message": "ONE_OF_YOU_BLOCK_EACH_OTHER"
}
```
2.Retrieve friends by email: (POST /api/user/relationship/list)
2.1 Request body
```
email: the email address of user need to get list friendship
```
+ Example:
```
{
    "email" : "trendy@example.com"
}
```
2.2 Response body
+ Success:
```
{
    "success": true,
    "friends": [
        "mandy@example.com",
        "alameda@example.com",
    ],
    "count": 2
}
```
+ invalid_email_input:
```
{
    "success": false,
    "message": "INVALID_EMAIL_INPUT"
}
```
+ Fail:
```
{
    "success": false,
    "message": "DATABASE_ERROR"
}
```
3.Get commond friends: (POST /api/user/relationship/common-friends)
3.1 Request body
```
friends: array of two emails that need to get list commond friend
```
+ Example:
```
{
    "friends" : ["bingo@example.com", "trendy@example.com"]
}
```
3.2 Response body
+ Success:
```
{
    "success": true,
    "friends": [
        "mandy@example.com",
        "alameda@example.com",
    ],
    "count": 2
}
```
+ invalid_one_of_two_email_input:
```
{
    "success": false,
    "message": "INVALID_EMAIL_INPUT"
}
```

+ invalid_one_of_two_email_missing:
```
{
    "success": false,
    "message": "AT_LEAST_TWO_EMAILS_ARE_REQUIRED"
}
```
+ fail_one_of_two_email_block_the_other:
```
{
    "success": false,
    "message": "ONE_OF_YOU_BLOCK_EACH_OTHER"
}
```
4.Subscribe to updates: (POST /api/user/relationship/subscriber)
4.1 Request body
```
requestor: email of user needs to subscribe
target: email of user will get a subscriber
```
+ Example:
```
{
    "requestor": "micky@example.com",
    "target": "trendy@example.com"
}
```
4.2 Response body
+ Success:
```
{
    "success": true
}
```
+ invalid_one_of_two_email_input:
```
{
    "success": false,
    "message": "INVALID_EMAIL_INPUT"
}
```

+ invalid_one_of_two_email_missing:
```
{
    "success": false,
    "message": "AT_LEAST_TWO_EMAILS_ARE_REQUIRED"
}
```
+ fail_already_susbcribe:
```
{
    "success": false,
    "message": "YOU_ALREADY_SUBSCRIBED"
}
```
+ fail_one_of_two_email_block_the_other:
```
{
    "success": false,
    "message": "ONE_OF_YOU_BLOCK_EACH_OTHER"
}
```
5.Block updates: (POST /api/user/relationship/block)
5.1 Request body
```
requestor: email of user want to block
target: email of user will be blocked
```
+ Example:
```
{
    "requestor": "micky@example.com",
    "target": "trendy@example.com"
}
```
5.2 Response body
+ Success:
```
{
    "success": true
}
```
+ invalid_one_of_two_email_input:
```
{
    "success": false,
    "message": "INVALID_EMAIL_INPUT"
}
```

+ invalid_one_of_two_email_missing:
```
{
    "success": false,
    "message": "AT_LEAST_TWO_EMAILS_ARE_REQUIRED"
}
```
+ fail:
```
{
    "success": false,
    "message": "DATABASE_ERROR"
}
```
6.Get recipient: (POST /api/user/relationship/recipients)
6.1 Request body
```
requestor: the author user email of the update
text: content of the update
```
+ Example:
```
{
    "sender" : "trendy@example.com",
    "text": "something is about to happle mrbean@xyz.com luis@example.com"
}
```
6.2 Response body
+ Success:
```
{
    "success": true,
    "recipients": [
        "mandy@example.com",
        "mrbean@xyz.com",
        "luis@example.com"
    ]
}
```
+ invalid_sender_email_required:
```
{
    "success": false,
    "message": "SENDER_IS_REQUIRED"
}
```
+ invalid_sender_email_input:
```
{
    "success": false,
    "message": "INVALID_EMAIL_INPUT"
}
```
