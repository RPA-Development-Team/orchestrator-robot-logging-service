# Robot Logging Service
This service receives and processes logs produced by robots on remote machines.

## Setting up
Clone the repository and make sure Go is installed, then run 
```go run main.go```

## Endpoints
The endpoints currently implemented are:
- POST /api/robot/login
- GET /rtlogs
### /api/robot/login Endpoint
This endpoint is used to authenticate the robot using the user's credentials. The request to be sent to the endpoint must look like the following example:
```json
{
    "username": "Example",
    "password": "pass"
}
```
The service then checks if the credentials are valid by communicating with the authentication service. If the credentials are valid, a one-time-use robot token is generated and sent to the robot for establishing a websocket connection with the server. The token is sent in the following format:
```json
{
    "token": "Example Token",
    "userId": "Example ID"
}
```
If invalid credentials were sent, the server responds with an HTTP BadRequest (status code 400) and the following JSON response:
```json
{
    "message": "Invalid credentials"
}
```
If any errors occur during the authentication process, the server responds with an HTTP BadRequest and the folowing JSON response:
```json
{
    "error": "Error details"
}
```
### /rtlogs Endpoint
This endpoint is concerened with establishing a websocket connection with the client (robot). After the robot obtains their token, a GET request must be sent to this endpoint with the token being included as a query parameter:
`/rtlogs?token="myToken"`. If the token is invalid, the server responds with status code 400 and the following JSON response:
```json
{
    "message": "Invalid token"
}
```
Subsequent communication between this service and the clients must be event based and must conform to the following format:
```json
{
    "eventType": "logEmitEvent",
    "payload": "myPayload"
}
```
### Supported Events
This is a list of all supported events:
- logEmitEvent - Sent by the client to the service to register logs
- logReceiveEvent - Sent by service to client to notify it of log reception
- errorMessageEvent - Sent by the service to the client to notify it of potential errors
### LogEmitEvent
The format for this event looks as follows:
```json
{
    "eventType": "logEmitEvent",
    "payload": {
        "logType": "ERROR",
        "name": "MessageBox",
        "status": "Running",
        "timestamp": "12345",
        "message": "this is a log entry",
        "robotId": 1
    }
}
```
If the sent log was successfully saved, the following response is sent to the client:
```json
{
    "eventType": "logReceiveEvent",
    "payload": {
        "message": "Log entry received.",
    }
}
```