# Chirpy

Chirpy is a simple microblogging platform written in Go. Users can register, log in, and create "chirps" that are associated with their accounts. This project serves as a personal practice application to explore web development in Go and is intended to be run locally.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Contact](#contact)

## Features

- **Local Deployment**: Run the application and database locally on your machine.
- **User Authentication**: Register and log in to create chirps.
- **Chirp Management**: Create, read, and delete chirps.
- **RESTful API**: Interact with the application via well-defined API endpoints.
- **Debug Mode**: Enable debug mode to reset the local database.

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) installed on your machine.

### Steps

1. **Clone the Repository**

   ```bash
   git clone https://github.com/Romasav/chirpy.git
   cd chirpy
   ```

2. **Set Up Environment Variables**

   Create a `.env` file in the root directory and add the following:

   ```env
   JWT_SECRET=N7wOWgJmqC4AbOUpzJcIcWM0lQTNxSfjU8ghpAI3jTKGauGbFx6RFfhKF7DEtpmuHjhQm9X94IHuZYg6AWRb0wio
   POLKA_KEY=f271c81ff7084ee5b99a5091b42d486e
   ```

   > **Note:** These are fake keys provided for testing purposes. You can use them as-is or generate your own.

3. **Build the Application**

   ```bash
   go build -o chirpy
   ```

4. **Run the Application**

   ```bash
   ./chirpy
   ```

   The server will start on port `8080` by default.

### Running in Docker

If you'd like to run **Chirpy** in Docker, follow these steps:

1. **Pull the Docker Image:**

   ```bash
   docker pull kavuunnn/chirpy:latest
   ```

2. **Run the Docker Container:**

   Once the image is pulled, you can run the application inside a Docker container:

   ```bash
   docker run -p 8080:8080 kavuunnn/chirpy
   ```

3. **Keeping Data Persistent Between Containers:**

   By default, the `database.json` file (used to store application data) will be created inside the container. However, this file will be lost when the container is removed. If you want to keep your data persistent between different containers, you can mount a Docker volume to store the `database.json` file outside the container.

   Use the following command to ensure your data is persisted:

   ```bash
   docker run -v chirpy-vol:/app/data -p 8080:8080 kavuunnn/chirpy
   ```

   In this case, Docker will create a volume called `chirpy-vol` to store the `database.json` file, and this volume will persist even if the container is removed. The file will be stored in a Docker-managed location that can be reused by other containers.

## Usage

### Local Deployment

- The application and its database (`database.json`) run locally on your machine.
- All data is stored locally, and the database is a simple JSON file.

### Command-Line Flags

- **Debug Mode**

  Enable debug mode to reset the local database by deleting the `database.json` file:

  ```bash
  ./chirpy -debug
  ```

### Accessing the Application

- Open your browser and navigate to `http://localhost:8080/app/` to access the application interface.
- All API endpoints are accessible under the `/api/` path.

### Authentication

- You must register and log in to create chirps.
- Authentication is handled via JWT tokens, which are provided upon successful login.

## API Endpoints

Below is a detailed description of each API endpoint available in the Chirpy application. This section provides information on the endpoint URLs, HTTP methods, request parameters, request bodies, and expected responses.

---

### Health Check

**Endpoint**

```
GET /api/healthz
```

**Description**

Checks the readiness of the Chirpy application. Useful for health monitoring.

**Request Headers**

- None

**Response**

- **Success (200 OK)**

  Returns a plain text message indicating the server is ready.

  **Example**

  ```
  OK
  ```

---

### Get All Chirps

**Endpoint**

```
GET /api/chirps
```

**Description**

Retrieves a list of all chirps. Supports filtering by `author_id` and sorting.

**Request Headers**

- None

**Query Parameters**

- `author_id` (integer, optional): Filters chirps by the author's user ID.
- `sort` (string, optional): Sorts the chirps. Accepts `asc` or `desc`. Default is `asc`.

**Example Request**

```
GET /api/chirps?author_id=1&sort=desc
```

**Response**

- **Success (200 OK)**

  Returns a JSON array of chirp objects.

  **Example**

  ```json
  [
    {
      "id": 2,
      "body": "Hello, world!",
      "author_id": 1,
      "created_at": "2023-10-01T12:34:56Z"
    },
    {
      "id": 1,
      "body": "My first chirp!",
      "author_id": 1,
      "created_at": "2023-10-01T12:00:00Z"
    }
  ]
  ```

- **Error Responses**

  - **500 Internal Server Error**

    ```json
    {
      "error": "Failed to load chirps"
    }
    ```

---

### Get Chirp by ID

**Endpoint**

```
GET /api/chirps/{chirpID}
```

**Description**

Retrieves a specific chirp by its ID.

**Request Headers**

- None

**URL Parameters**

- `chirpID` (integer, required): The ID of the chirp to retrieve.

**Example Request**

```
GET /api/chirps/1
```

**Response**

- **Success (200 OK)**

  Returns a JSON object of the chirp.

  **Example**

  ```json
  {
    "id": 1,
    "body": "My first chirp!",
    "author_id": 1,
    "created_at": "2023-10-01T12:00:00Z"
  }
  ```

- **Error Responses**

  - **404 Not Found**

    ```json
    {
      "error": "The chirp with id = 1 was not found"
    }
    ```

  - **500 Internal Server Error**

    ```json
    {
      "error": "Failed to convert chirpIDStr to int"
    }
    ```

---

### Create a New Chirp

**Endpoint**

```
POST /api/chirps
```

**Description**

Creates a new chirp associated with the authenticated user.

**Request Headers**

- `Authorization: Bearer {token}`

**Request Body**

A JSON object containing the chirp content.

- `body` (string, required): The content of the chirp.

**Example**

```json
{
  "body": "This is a new chirp!"
}
```

**Response**

- **Success (201 Created)**

  Returns the created chirp object.

  **Example**

  ```json
  {
    "id": 3,
    "body": "This is a new chirp!",
    "author_id": 1,
    "created_at": "2023-10-01T13:00:00Z"
  }
  ```

- **Error Responses**

  - **400 Bad Request**

    ```json
    {
      "error": "Invalid JSON"
    }
    ```

  - **401 Unauthorized**

    ```json
    {
      "error": "Authorization header is required"
    }
    ```

    ```json
    {
      "error": "Invalid or expired token"
    }
    ```

  - **500 Internal Server Error**

    ```json
    {
      "error": "Could not create chirp"
    }
    ```

---

### Delete a Chirp

**Endpoint**

```
DELETE /api/chirps/{chirpID}
```

**Description**

Deletes a chirp by ID. Only the author of the chirp can delete it.

**Request Headers**

- `Authorization: Bearer {token}`

**URL Parameters**

- `chirpID` (integer, required): The ID of the chirp to delete.

**Example Request**

```
DELETE /api/chirps/3
```

**Response**

- **Success (204 No Content)**

  The chirp was successfully deleted. No content is returned.

- **Error Responses**

  - **401 Unauthorized**

    ```json
    {
      "error": "Authorization header is required"
    }
    ```

    ```json
    {
      "error": "Invalid or expired token"
    }
    ```

  - **403 Forbidden**

    ```json
    {
      "error": "You can't delete chirps that were created by someone else"
    }
    ```

  - **404 Not Found**

    ```json
    {
      "error": "The chirp with id = 3 was not found"
    }
    ```

  - **500 Internal Server Error**

    ```json
    {
      "error": "Failed to convert chirpIDStr to int"
    }
    ```

---

### Register a New User

**Endpoint**

```
POST /api/users
```

**Description**

Registers a new user account.

**Request Headers**

- `Content-Type: application/json`

**Request Body**

- `email` (string, required): The user's email address.
- `password` (string, required): The user's password.

**Example**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response**

- **Success (201 Created)**

  Returns the created user's information.

  **Example**

  ```json
  {
    "id": 1,
    "email": "user@example.com",
    "is_chirpy_red": false
  }
  ```

- **Error Responses**

  - **400 Bad Request**

    ```json
    {
      "error": "Invalid JSON"
    }
    ```

    ```json
    {
      "error": "Couldn't create user"
    }
    ```

---

### User Login

**Endpoint**

```
POST /api/login
```

**Description**

Authenticates a user and returns an access token and refresh token.

**Request Headers**

- `Content-Type: application/json`

**Request Body**

- `email` (string, required): The user's email address.
- `password` (string, required): The user's password.

**Example**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response**

- **Success (200 OK)**

  Returns user information along with tokens.

  **Example**

  ```json
  {
    "id": 1,
    "email": "user@example.com",
    "is_chirpy_red": false,
    "token": "access_token_jwt",
    "refresh_token": "refresh_token_value"
  }
  ```

- **Error Responses**

  - **400 Bad Request**

    ```json
    {
      "error": "Invalid JSON"
    }
    ```

  - **401 Unauthorized**

    ```json
    {
      "error": "Incorrect Password"
    }
    ```

  - **404 Not Found**

    ```json
    {
      "error": "The user was not found"
    }
    ```

  - **500 Internal Server Error**

    ```json
    {
      "error": "Error getting a token"
    }
    ```

---

### Update User Information

**Endpoint**

```
PUT /api/users
```

**Description**

Updates the authenticated user's email and password.

**Request Headers**

- `Authorization: Bearer {token}`
- `Content-Type: application/json`

**Request Body**

- `email` (string, required): The new email address.
- `password` (string, required): The new password.

**Example**

```json
{
  "email": "newemail@example.com",
  "password": "newsecurepassword123"
}
```

**Response**

- **Success (200 OK)**

  Returns the updated user's information.

  **Example**

  ```json
  {
    "id": 1,
    "email": "newemail@example.com",
    "is_chirpy_red": false
  }
  ```

- **Error Responses**

  - **400 Bad Request**

    ```json
    {
      "error": "Invalid JSON"
    }
    ```

  - **401 Unauthorized**

    ```json
    {
      "error": "Authorization header is required"
    }
    ```

    ```json
    {
      "error": "Invalid or expired token"
    }
    ```

  - **500 Internal Server Error**

    ```json
    {
      "error": "Failed to update user"
    }
    ```

---

### Refresh Access Token

**Endpoint**

```
POST /api/refresh
```

**Description**

Generates a new access token using a valid refresh token.

**Request Headers**

- `Authorization: Bearer {refresh_token}`

**Response**

- **Success (200 OK)**

  Returns a new access token.

  **Example**

  ```json
  {
    "token": "new_access_token_jwt"
  }
  ```

- **Error Responses**

  - **401 Unauthorized**

    ```json
    {
      "error": "Authorization header is required"
    }
    ```

    ```json
    {
      "error": "The refresh token is invalid"
    }
    ```

  - **500 Internal Server Error**

    ```json
    {
      "error": "Error getting a token"
    }
    ```

---

### Revoke Refresh Token

**Endpoint**

```
POST /api/revoke
```

**Description**

Revokes a refresh token, effectively logging the user out.

**Request Headers**

- `Authorization: Bearer {refresh_token}`

**Response**

- **Success (204 No Content)**

  The refresh token was successfully revoked.

- **Error Responses**

  - **401 Unauthorized**

    ```json
    {
      "error": "Authorization header is required"
    }
    ```

    ```json
    {
      "error": "The refresh token is invalid"
    }
    ```

  - **500 Internal Server Error**

    ```json
    {
      "error": "Failed to delete refresh token"
    }
    ```

---

### Handle Polka Webhooks

**Endpoint**

```
POST /api/polka/webhooks
```

**Description**

Handles incoming webhooks from Polka. Currently supports upgrading a user to Chirpy Red.

**Request Headers**

- `Authorization: ApiKey {POLKA_KEY}`
- `Content-Type: application/json`

**Request Body**

- `event` (string, required): The event type. Currently supports `"user.upgraded"`.
- `data` (object, required): Event data.

  - `user_id` (integer, required): The ID of the user to upgrade.

**Example**

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": 1
  }
}
```

**Response**

- **Success (204 No Content)**

  The user was successfully upgraded.

- **Error Responses**

  - **400 Bad Request**

    ```json
    {
      "error": "Invalid JSON"
    }
    ```

  - **401 Unauthorized**

    ```json
    {
      "error": "Authorization header is required"
    }
    ```

    ```json
    {
      "error": "Incorrect key"
    }
    ```

  - **404 Not Found**

    ```json
    {
      "error": "The user was not found"
    }
    ```

---

### Admin Metrics

**Endpoint**

```
GET /admin/metrics
```

**Description**

Displays the number of times the file server has been accessed.

**Request Headers**

- None

**Response**

- **Success (200 OK)**

  Returns an HTML page with the metrics.

  **Example**

  ```html
  <html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 42 times!</p>
  </body>
  </html>
  ```

---

### Reset Metrics

**Endpoint**

```
POST /api/reset
```

**Description**

Resets the file server hit count to zero.

**Request Headers**

- None

**Response**

- **Success (200 OK)**

  Resets the metrics. No content is returned.

---

## Notes on Authentication

- **Access Token**

  - Used for authenticating standard API requests.
  - Obtained via the `/api/login` endpoint.
  - Included in the `Authorization` header as `Bearer {token}`.

- **Refresh Token**

  - Used to obtain a new access token when the current one expires.
  - Obtained via the `/api/login` endpoint.
  - Included in the `Authorization` header as `Bearer {refresh_token}` for `/api/refresh` and `/api/revoke` endpoints.

- **Polka Key**

  - A special API key used by Polka to authenticate webhook requests.
  - Included in the `Authorization` header as `ApiKey {POLKA_KEY}`.

---

## Configuration

- **Port Number**

  The server runs on port `8080` by default. You can change this by modifying the `port` variable in the `main.go` file:

  ```go
  const port = "8080"
  ```

- **Environment Variables**

  The application uses a `.env` file for configuration. Make sure to set the `JWT_SECRET` and `POLKA_KEY` as shown in the installation steps.

- **Database**

  - The application uses a local JSON file (`database.json`) to store data.
  - The database file is generated automatically upon running the application.
  - To reset the database, you can delete the `database.json` file manually or run the application in debug mode.

## Contact

*This project is a personal practice application intended to be run locally and is not intended for production use.*