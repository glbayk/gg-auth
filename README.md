# gg-auth - Generalized JWT Bearer Token Authentication and Authorization Boilerplate

[Gin](https://github.com/gin-gonic/gin) + [Gorm](https://github.com/go-gorm/gorm) + [GoDotEnv](https://github.com/joho/godotenv)

## Overview

It is not a full fledged authentication and authorization package. It is a boilerplate that can be used to start a project. All the endpoints are open to the world. You can add your own endpoints and middleware to protect them.

> **Note:** The user model has been kept simple. It is recommended to add more fields to the user model or extend it via different models and relations.

## Endpoints

### Open to the world:

**Health Check**

- **GET** `/ping` - check if the server is up and running

**Auth**

- **POST** `/auth/register` - register a new user
- **POST** `/auth/login` - login a user and get access and refresh tokens
- **GET** `/auth/refresh-token` - get a new access token using a refresh token
- **GET** `/auth/reset-password` - reset a user's password

### Protected

**User**

- **GET** `/user/profile` - get the user's profile information (requires access token)

## Getting started

1. Go install the package:

```bash
go install github.com/glbayk/gg-auth
```

2. Copy the `.env.example` file to `.env` and replace the values with your own:
3. Build the project:

```bash
go build
```

4. Run the project:

```bash
./gg-auth
```

## QnA (Questions and Answers)
