# go-api

This api is a security-focused system that combines role-based authentication, streamlined CRUD operations (Create, Read, Update, Delete), and efficient user management. It empowers administrators to control access based on user roles while providing a seamless experience for data manipulation and user interactions.

## Architecture Overview

```bash
├───cmd
├───docker
└───pkg
    ├───auth
    ├───db
    ├───entities
    ├───errors
    └───server
```

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`PORT`, `DB_URI`, `DB_NAME`, `JWT_SECRET`

## Run Locally

#### Make sure Docker, Git and Go is installed.

Clone the project

```bash
  git clone https://github.com/Merloss/go-api.git
```

Go to the project directory

```bash
  cd go-api
```

Run services

```bash
  docker-compose -f /docker/docker-compose.yml up -d
```

Install dependencies

```bash
  go mod tidy
```

Start the server

```bash
  go run .\cmd\main.go
```

## API Reference

#### Register

```http
  POST /api/auth/register
```

| Field name | Type     | Description                 |
| :--------- | :------- | :-------------------------- |
| `username` | `string` | **Required**. Your username |
| `password` | `string` | **Required**. Your password |

#### Login

```http
  POST /api/auth/login
```

| Field name | Type     | Description                 |
| :--------- | :------- | :-------------------------- |
| `username` | `string` | **Required**. Your username |
| `password` | `string` | **Required**. Your password |

#### Get all posts

```http
  GET /api/posts
```

#### Create post (editor only)

```http
  POST /api/posts
```

| Field name    | Type     | Description                    |
| :------------ | :------- | :----------------------------- |
| `title`       | `string` | **Required**. Post title       |
| `description` | `string` | **Required**. Post description |

#### Edit post (editor only)

```http
  POST /api/posts/:id
```

| Field name    | Type     | Description                    |
| :------------ | :------- | :----------------------------- |
| `title`       | `string` | **Required**. Post title       |
| `description` | `string` | **Required**. Post description |
| `status`      | `string` | **Admin only** Post status     |

#### Delete post (admin only)

```http
  DELETE /api/posts/:id
```

| Parameter | Type     | Description           |
| :-------- | :------- | :-------------------- |
| `id`      | `string` | **Required**. Post ID |

#### Get pending posts (admin only)

```http
  GET /api/posts/pending
```

#### Edit user (admin only)

```http
  POST /api/users/:id/edit
```

| Field name | Type       | Description   |
| :--------- | :--------- | :------------ |
| `username` | `string`   | User username |
| `roles`    | `string[]` | User roles    |
