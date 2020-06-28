# go-todo-backend

Go Todo Backend Example Using Modular Project Layout for Product Microservice. It's suitable as starting point for a medium to larger project.

This example uses [Chi](https://github.com/go-chi/chi) for http router and [REL](https://github.com/Fs02/rel) for database access.

Feature:

- Modular Project Structure.
- Full example including tests.
- Docker deployment.
- Compatible with [todobackend](https://www.todobackend.com/specs/index.html).

## Installation

### Prerequisite

1. Install [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) for database migration.
2. Install [mockery](https://github.com/vektra/mockery#installation) for interface mock generation.

### Running

1. Prepare `.env`.
    ```
    cp .env.sample .env
    ```
2. Create a database and update `.env`.
2. Prepare database schema.
    ```
    make migrate
    ```
3. Build and Running
    ```
    make
    ```

## Project Structure

```
.
├── api
│   ├── handler
│   │   ├── todos.go
│   │   └── [other handler].go
│   └── middleware
│       └── [other middleware].go
├── bin
│   ├── api
│   └── [other executable]
├── cmd
│   ├── api
│   │   └── main.go
│   └── [other cmd]
│       └── main.go
├── db
│   ├── schema.sql
│   └── migrations
│       └── [migration file]
├── todos
│   ├── todo.go
│   ├── create.go
│   ├── update.go
│   ├── delete.go
│   ├── service.go
│   └── todostest
│       ├── todo.go
│       └── service.go
├── [other domain]
│   ├── [entity a].go
│   ├── [business logic].go
│   ├── [other domain]test
│   │   └── service.go
│   └── service.go
└── [other client]
    ├── [entity b].go
    ├── client.go
    └── [other client]test
        └── client.go
```

This project structure is based on a modular project structure, with loosely coupled dependencies between domain, Think of making libraries under a single repo that only exports certain functionality that used by other service and http handler. One of domain that present in this example is todos.

Loosely coupled dependency between domain is enforced by avoiding the use of shared entity package, therefore any entity struct should be included inside it's own respective domain. This will prevent cyclic dependency between entity. This shouldn't be a problem in most cases, becasause if you encounter cyclic dependency, there's huge chance that the entity should belongs to the same domain.

For example, consider three structs: user, transaction and transaction items. transaction and its transaction items might need cyclic dependency and items doesn't works standalone (items without transaction should not exists), thus it should be on the same domain.
In the other hand, user and transaction shouldn't require cyclic dependency, transaction might have a user field in the struct, but user shouldn't have a slice of transaction field, therefore it should be on a separate domain.

### Domain vs Client

Domain and Client folder is very similar, the difference is client folder doesn't actually implement any business logic (service), but instead a client that calls any internal/external API to works with the domain entity.
