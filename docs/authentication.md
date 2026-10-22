# Authentication

This API uses JSON Web Tokens (JWT) for authenticating requests to protected endpoints.

---

## 1. Obtaining a Token

To receive an access token, send a `POST` request to `/api/login` with your credentials:

- **Endpoint:** `POST /api/login`
- **Headers:** `Content-Type: application/json`
- **Request Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "yourpassword"
  }
  ```
- **Response:** `200 OK`
  ```json
  {
    "id": "1",
    "email": "user@example.com",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
  ```

---

## 2. Authenticating Requests

Include the token in the `Authorization` header using the `Bearer` scheme on any protected endpoint:

```http
Authorization: Bearer <your_token_here>
```

### Example Request using `curl`:

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -H "Content-Type: application/json" \
  -d '{"body": "Hello world"}'
```

---

## 3. Token Expiration

- Access tokens expire **1 hour** after issuance.
- Once expired, any request using that token will return `401 Unauthorized`.
- You must call `/api/login` (or `/api/refresh` if refresh tokens are supported) to obtain a new token.

---

## 4. Error Codes

| Status Code | Meaning | Cause |
| :--- | :--- | :--- |
| `401 Unauthorized` | Missing or Invalid Token | Token was omitted, malformed, expired, or signed with an invalid key. |
| `403 Forbidden` | Insufficient Permissions | Token is valid, but the user is not allowed to access the requested resource. |