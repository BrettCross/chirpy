# API Reference

Base URL: `http://localhost:8080/api`

---

## 1. Get All Chirps

Retrieves a list of recent chirps.

- **URL:** `/chirps`
- **Method:** `GET`
- **Success Response:**
  - **Code:** `200 OK`
  - **Body:**
    ```json
    [
      {
        "id": "1",
        "body": "Hello world!",
        "author_id": "42"
      }
    ]
    ```

---

## 2. Create a Chirp

Creates a new chirp entry.

- **URL:** `/chirps`
- **Method:** `POST`
- **Headers:**
  - `Content-Type: application/json`
  - `Authorization: Bearer <token>`
- **Request Body:**
  ```json
  {
    "body": "This is a new chirp"
  }
  ```
- **Responses:**
  - **Code:** `201 Created`
    ```json
    {
      "id": "2",
      "body": "This is a new chirp",
      "author_id": "42"
    }
    ```
  - **Code:** `400 Bad Request`
    ```json
    {
      "error": "Chirp is too long"
    }
    ```
  - **Code:** `401 Unauthorized`
    ```json
    {
      "error": "Invalid token"
    }
    ```