# API Payment Transfer

API Payment Transfer is simulation API for transfer money for registered user.

## Step of Using API Payment Transfer

---

Get the repository :

1. Clone the repository
1. running the repository with:

```go
go run "main.go"
```

## Testing API

---

there 3 API to use:

1. Login
2. Payment
3. Logout

Please use Postman for testing the API

### 1. Login

You can access the endpoint with below endpoint:

```
POST -- /login
```

the require of this endpoint is:
| Name | Value | Description |
| --- | --- | --- |
| Application/json | <br>- username<br>- password | fill **username** and **password** with data in user.json |

after request endpoint **/login** you will get:

```json
{
  "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJmYXR1ciIsImZ1bGxuYW1lIjoiZmF0dXIiLCJub3JlayI6IjEyMzQ1Njc4OSIsImJhbGFuY2UiOjI1MDAwMDAsImNyZWF0ZWRfYXQiOiIyMDIyLTAxLTAxIiwidXBkYXRlZF9hdCI6IjIwMjItMDEtMTMgMTc6MjE6MDkuNDA4MzA5OCArMDcwMCArMDcgbT0rMTYuMTk1NDg0ODAxIiwiZXhwIjoxNjQyMDg4MTQzfQ.z4QZ-nqb7RmnriTmqN4V03jLhuUPMTSI_HSsel-nhT4",
  "data": {
    "id": 1,
    "username": "ijat",
    "fullname": "ijat",
    "norek": "123456789",
    "balance": 2500000,
    "created_at": "2022-01-01",
    "updated_at": "2022-01-13 17:21:09.4083098 +0700 +07 m=+16.195484801",
    "exp": 1642088143
  },
  "error": false,
  "expires": "2022-01-13T22:35:43.0751525+07:00",
  "message": "Success Login"
}
```

Copy Token value for access payment with Authorization with type Bearer Token

### 2. Payment

You can access the endpoint with below endpoint:

```
POST -- /payment
```

the require of this endpoint is:
| Name | Value | Description |
| --- | --- | --- |
|Authorization|Bearer Token|Paste the Token from request **/Login** endpoint|
| Application/json | <br>- nominal<br>- to | fill **nominal** with the amount you want transfer and **to** for the account number (Norek) from user.json |

after request endpoint **/payment** you will get:

```json
{
  "data": {
    "id": 32,
    "user_id": 1,
    "from": "123456789",
    "to": "987654321",
    "Nominal": 500000,
    "created_at": "2022-01-13 22:57:23.8038722 +0700 +07 m=+82.052272101",
    "updated_at": "2022-01-13 22:57:23.8038722 +0700 +07 m=+82.052272101"
  },
  "error": false,
  "message": "success transfer payment"
}
```

### 3. Logout

You can access the endpoint with below endpoint:

```
GET -- /Logout
```

after request endpoint **/logout** you will get:

```json
{
  "error": false,
  "message": "User Logout"
}
```
