# Go Backend API

This is a RESTful API built with Go (Chi router), PostgreSQL (pgx), and JWT Authentication.

## Base URL
`http://localhost:3333`

## Authentication
Authentication is handled via **JWT (JSON Web Tokens)**.
For protected routes, include the token in the request header:
```
Authorization: Bearer <your_token>
```

---

## Allowed Methods & Endpoints

### 🔓 Public Routes

#### 1. Login
*   **Method:** `POST`
*   **Endpoint:** `/login`
*   **Description:** Authenticate user and receive a JWT token.
*   **Request Body:**
    ```json
    {
      "email": "tarun@example.com",
      "password": "mypassword123"
    }
    ```
*   **Response:**
    ```json
    {
      "token": "eyJhbGciOiJIUzI1Ni..."
    }
    ```

#### 2. Register User
*   **Method:** `POST`
*   **Endpoint:** `/users`
*   **Description:** Register a new user account.
*   **Request Body:**
    ```json
    {
      "username": "tarun",
      "name": "Tarun Monga",
      "email": "tarun@example.com",
      "password": "mypassword123"
    }
    ```
*   **Response:**
    ```json
    {
      "id": "550e8400-e29b-41d4-a716-446655440000"
    }
    ```

### 🔒 Protected Routes (Requires Token)

#### 3. Get Current User (Me)
*   **Method:** `GET`
*   **Endpoint:** `/me`
*   **Description:** Get the profile of the currently logged-in user.
*   **Response:**
    ```json
    {
      "id": "uuid...",
      "username": "tarun",
      "name": "Tarun Monga",
      "email": "tarun@example.com",
      "profilePic": "",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
    ```

#### 4. Get All Users
*   **Method:** `GET`
*   **Endpoint:** `/users`
*   **Description:** Retrieve a list of all users.

#### 5. Update User
*   **Method:** `PUT`
*   **Endpoint:** `/users/{id}`
*   **Description:** Update a user's details. You can only update your own ID.
*   **Request Body:**
    ```json
    {
      "username": "tarun_new",
      "name": "Tarun Updated",
      "email": "tarun@example.com",
      "profilePic": "http://example.com/pic.jpg",
      "password": "newpassword123" // Optional
    }
    ```
