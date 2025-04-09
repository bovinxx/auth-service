# auth-service
Authentication and Authorization Service built with gRPC, Protocol Buffers, and Go

---

## 📦 API Overview

### `AuthService`
- `Login`: Authenticate via username/password
- `GetAccessToken`: Exchange refresh token for access token
- `GetRefreshToken`: Rotate refresh token

### `UserService`
- `Create`, `Get`, `Update`, `Delete`: Manage users
- Uses HTTP annotations for REST compatibility

### `AccessService`
- `Check`: Check if user has access to an endpoint

---

## 🛠️ Technologies

- **Go** (1.24)
- **gRPC**, **Protocol Buffers**
- **grpc-gateway** for REST API
- **JWT** for tokens

---
