ToDo List
Objective
Create a simple RESTful API in Go that implements a ToDo list. The API should be able to handle the routes and actions below. The list can be persisted in any way. You can use a flat file, a database or just store the list in memory. Please only use the standard Go library. Tests are encouraged. Please create a GitHub repository to share the code with us. It is encouraged that you make several commits so that we can see your thought process as you solve this problem. You are allowed to add anything else to the app as long as all the requirements are met.

List ToDos
Route - GET /todos
Returns a list of ToDo objects
Example Response:
[{
“id”: 1,
“name”: “Feed Dog”,
“complete”: false
},
{
“id”: 2,
“name”: “Wash the dishes”,
“complete”: true
}]

Get ToDo By ID
Route - GET /todos/id (id is an integer)
Returns a ToDo object or a 404 error if the id does not exist
Example Request: GET /todos/1
Example Response:
{
“id”: 1,
“name”: “Feed Dog”,
“complete”: false
}

Create ToDo
Route - POST /todos
Creates a new ToDo object and returns it
Example Request: POST /todos
{
“name”: “Feed Dog”,
“complete”: false
}
Example Response:
{
“id”: 1,
“name”: “Feed Dog”,
“complete”: false
}
Update ToDo
Route - PUT /todos/id (id is an integer)
Updates the ToDo that has the ID of from the route
Returns a 404 if the todo does not exist
Example Request: PUT /todos/1
{
“complete”: true
}
Example Response: 204 with no body
Delete ToDo
Route - DELETE /todos/id (id is an integer)
Deletes the ToDo that has the ID from the route
Returns a 404 if the ToDo does not exist
Example Request: DELETE /todos/1
Example Response: 204 with no body
