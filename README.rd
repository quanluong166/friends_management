FRIEND_MANAGEMENT

+ Step to build
// Need to install postgres, docker first

1. make build
2. make run-postgres
3. make run

+ Init sample data
1. make migrate

+ User Stories

1. As a user, I need an API to create a friend connection between two email addresses.
The API should receive the following JSON request:
{
    "friends" : ["friend1@example.com", "friend2@example.com"]
}
The API should return the following JSON response on success:
{
    "success": true
}
Please propose JSON responses for any errors that might occur.

2. As a user, I need an API to retrieve the friends list for an email address.
The API should receive the following JSON request:
{
    "email" : "trendy@example.com"
}
The API should return the following JSON response on success:
{
    "success": true,
    "friends": [
        "mandy@example.com",
        "alameda@example.com",
        "bingo@example.com"
    ],
    "count": 3
}
Please propose JSON responses for any errors that might occur.


3. As a user, I need an API to retrieve the common friends list between two email addresses.
The API should receive the following JSON request:
{
    "friends" : ["bingo@example.com", "trendy@example.com"]
}
The API should return the following JSON response on success:
{
    "success": true,
    "friends": [
        "alameda@example.com"
    ],
    "count": 1
}
Please propose JSON responses for any errors that might occur.

4. As a user, I need an API to subscribe to updates from an email address.
Please note that "subscribing to updates" is NOT equivalent to "adding a friend connection".
The API should receive the following JSON request:
{
    "requestor": "abc@example.com",
    "target": "trendy@example.com"
}
The API should return the following JSON response on success:
{
    "success": true
}
Please propose JSON responses for any errors that might occur.

5. As a user, I need an API to block updates from an email address.
Suppose "andy@example.com" blocks "john@example.com":
• if they are connected as friends, then "andy" will no longer receive
notifications from "john"
• if they are not connected as friends, then no new friends connection can be added
The API should receive the following JSON request:
{
    "requestor": "ranger@example.com",
    "target": "trendy@example.com"
}
The API should return the following JSON response on success:
{
    "success": true
}
Please propose JSON responses for any errors that might occur.

6. As a user, I need an API to retrieve all email addresses that can receive updates from an email address.
Eligibility for receiving updates from i.e. "john@example.com":
    • has not blocked updates from "john@example.com", and
    • at least one of the following:
        o has a friend connection with "john@example.com"
        o has subscribed to updates from "john@example.com" o has been @mentioned in the update
The API should receive the following JSON request:
{
    "sender" : "trendy@example.com",
    "text": "something is about to happle mrbean@xyz.com luis@example.com"
}
The API should return the following JSON response on success:
{
    "success": true,
    "recipients": [
        "mandy@example.com",
        "alameda@example.com",
        "bingo@example.com",
        "adison@example.com",
        "lucas@example.com",
        "micky@example.com",
        "mrbean@xyz.com",
        "luis@example.com"
    ]
}

+ Implementation

- Services functions:

1.  AddFriend()  --> API for make friend connection
2.  AddSubscriber()  --> API for make subscriber connection
3.  ListFriend()  --> API get list friends of one email
4.	ListCommonFriends()  --> API get list common friends between two emails
5.	AddBlock()  --> API for add block
6.	GetListEmailCanReceiveUpdate()  -> API get list emails can receive update from an email


- Usecase functions:

1.  AddFriendship()  --> handle business logic of make friend connection
2.	ListFriendships() --> handle business logic of get list friends of one email
3.	ListCommonFriends()-->  handle business logic of get list common friends between two emails
4.	AddSubscriber()  --> handle business logic of add subscriber to one email
5.	AddBlock()  --> handle business logic of add block
6.	GetListEmailCanReceiveUpdate() --> handle business logic of get list email can receive update from one email

- Repo functions:

1.  CreateFriendRelationship() --> create friend connection records of two email
2.	UpdateToFriendship() --> update subscriber relation to friendship
3.	GetListBlockEmail() --> get all emails that block the input email
4.	GetListSubscriberEmail() --> get all emails that is a subscriber of the input email
5.	GetListFriendshipEmail() --> get all emails that is a friend of the input email
6.	AddSubscriber() --> create subscribe connection record of two email
7.	CreateBlockRelationship() --> create block connection record of two email
8.	CheckTwoUsersBlockedEachOther() --> check if one of two email block the other one
9.	CheckTwoUsersAreFriends() --> check if two email has friend connection
10.	CheckIfTheRequestorAlreadySubscribe() --> check if the requestor already a suscriber of target email
11.	DeleteRelationship() --> delete all relationship record of two email
