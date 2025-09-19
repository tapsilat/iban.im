# IBAN.im
Shorten, create and share memorable links for IBANS

**BIG REFACTORING STARTED**


## Purpose
Shorten IBAN numbers with url such as :
- iban.im/user/alias
- iban.im/fakturk/garanti


## Stacks

- Go
- GraphQL : [graphql-go](https://github.com/graph-gophers/graphql-go)
- ORM : [gorm](https://github.com/jinzhu/gorm)


## Features

- [x] New users should Sign Up & Sign In
- [x] Change a Password of user
- [x] Change a Profile of user
- [ ] Delete a Profile of user
- [x] Get Profile of user
- [x] New IBAN add for user
- [x] Update IBAN  for user
- [x] Delete IBAN  for user
- [x] Get IBAN's of user
- [x] When adding new IBAN check if is it exist with same name (we can add with different names)
- [x] A user should add iban to only itself

## How to Run

### Initialize DB

1. Create a database

```shell
postgres=# CREATE DATABASE ibanim;
```

2. Create a user as owner of database

```shell
postgres=# CREATE USER ibanim WITH ENCRYPTED PASSWORD 'ibanim';

postgres=# ALTER DATABASE ibanim OWNER TO ibanim;
```

3. Grant all privileges to user for the database

```shell
postgres=# GRANT ALL PRIVILEGES ON DATABASE ibanim TO ibanim;
```

4. Configure the db in `db.go`

```go
// ConnectDB : connecting DB
func ConnectDB() (*DB, error) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=ibanim dbname=ibanim password=ibanim sslmode=disable")

	if err != nil {
		panic(err)
	}

	return &DB{db}, nil
}
```

or with Docker

> host address should be edited to `host.docker.internal` to connect a host interface.

```go
// ConnectDB : connecting DB
func ConnectDB() (*DB, error) {
	db, err := gorm.Open("postgres", "host=host.docker.internal port=5432 user=ibanim dbname=ibanim password=ibanim sslmode=disable")

	if err != nil {
		panic(err)
	}

	return &DB{db}, nil
}
```

### Initial Migration

```shell
$ go run ./migrations/init.go
```

or with Docker

```
$ docker build -t ibanim .
$ docker run --rm ibanim migrate
```

This will generate the `users` table in the database as per the User Model declared in `./model/user.go`

### Run the server

```shell
$ go run server.go
```

or with Docker

```
$ docker run --rm -d -p 8080:8080 ibanim
```

### GraphQL Playground

Connect to http://localhost:8080

### Authentication : JWT

You need to set the Http request headers `Authorization`: `{JWT_token}`

## Usage

### Sign Up

```graphql
mutation {
  signUp(
    email: "test@test.com"
    password: "12345678"
    firstName: "graphql"
    lastName: "go"
    handle:"test"
  ) {
    ok
    error
    user {
      id
      handle
      email
      firstName
      lastName
      bio
      avatar
      createdAt
      updatedAt
    }
  }
}
```

### Sign In

```graphql
mutation {
  signIn(email: "test@test.com", password: "12345678") {
    ok
    error
    token
  }
}
```

### Change a Password

```graphql
mutation {
  changePassword(password: "87654321") {
    ok
    error
    user {
      id
      handle
      email
      firstName
      lastName
      bio
      avatar
      createdAt
      updatedAt
    }
  }
}
```

### Change a Profile

```graphql
mutation {
  changeProfile(bio: "Go developer", avatar: "go-developer.png") {
    ok
    error
    user {
      id
      handle
      email
      firstName
      lastName
      bio
      avatar
      createdAt
      updatedAt
    }
  }
}
```

### Get my profile

```graphql
query {
  getMyProfile {
    ok
    error
    user {
      id
      handle
      email
      firstName
      lastName
      bio
      avatar
      createdAt
      updatedAt
    }
  }
}
```

### Add new Iban

```graphql
mutation {
  ibanNew(text:"TR320010009999901234567890",password:"fatih",handle:"fakturk"){
    ok
    error
    iban{
      id
      handle
      text
      password
      createdAt
      updatedAt
    }
  }
}
```

### Update Iban

```graphql
mutation {
  ibanUpdate(text:"TR420010009999901234567891",password:"fatih",handle:"garanti"){
    ok
    error
    iban{
      id
      handle
      text
      password
      createdAt
      updatedAt
    }
  }
}
```

### Get User IBANs

```graphql
query {
  getMyIbans {
    ok
    error
    iban {
       id
      handle
      text
      password
      createdAt
      updatedAt
      ownerId
    }
  }
}
```

