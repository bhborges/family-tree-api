# Family Tree API

The Family Tree API is a tool that allows developers to create and manage family trees in their applications. Through this API, it is possible to add, remove, and search for individuals, as well as establish relationships between them.
Endpoint

The base endpoint for the API is:

```javascript
https://api.family-tree.com/v1/
```

## Authentication

The API requires authentication through a valid access token. This token should be sent in all requests as an HTTP header Authorization with the value Bearer {token}. To obtain an access token, please contact the API administrator.
Available Endpoints
Creating Individuals
POST /individuals

Creates a new individual in the family tree.

Parameters:

* **'name'** (string, required): the full name of the `.

Returns:

* **'id'** (integer): the ID of the created individual.
* **'name'** (string): the full name of the individual.

Example request:


```json
POST /individuals
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "John Doe",
}
```
Example response:

```json

HTTP/1.1 200 OK
Content-Type: application/json

{
  "id": 1,
  "name": "John Doe",
}
```
Listing Members
GET /members

Lists all members in the family tree.

Returns:

  A list of members objects, each with the following fields:

* **'id'** (integer): the ID of the member.
* **'name'** (string): the full name of the member.

Example request:

```json

GET /members
Authorization: Bearer {token}
```
Example response:

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "members":[
    {
      "name": "Phoebe",
      "relationships": [
          {
            "name": "Martin",
            "relationship": "parent",
          },
          {
            "name": "Anastasia",
            "relationship": "parent",
          }
        ],
    },
    {
      "name": "Martin",
      "relationships": [],
    },
    {
      "name": "Anastasia",
      "relationships": [],
    }
  ]
}
```
