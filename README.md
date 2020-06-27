# go-todo-backend

Golang Todo Backend example with a complete project layout suitable as starting point for a larger project, It's built using [Chi](https://github.com/go-chi/chi) and [REL](https://github.com/Fs02/rel).

Demo: https://www.todobackend.com/client/index.html?https://go-todo-backend.herokuapp.com/

## Instalation

```
# Prepare .env
cp .env.sample .env

# Update dependencies
make dep

# Migrate
make migrate

# Build and Running
make
```

## Project Structure

```
.
├── api
│   └── handler
│   └── middleware
├── bin
├── cmd
│   └── api
│   └── [other cmd]
├── db
│   └── migrations
├── todos
├── [other domain]
└── [other client]
```

This project structure aims for a flat project structure, with loosely coupled dependency between domain. One of domain that present in this example is todos.

Loosely coupled dependency between domain is enforced by avoiding the use of shared entity package, therefore any entity struct should be included inside it's own respective domain. This will prevent cyclic dependency between entity.

In most cases, there shouldn't be a problem with this approach, as the case when you need cyclic dependency is mostly when the entity  is belongs to the same domain.

For example, consider three structs: user, transaction and transaction items: transaction and it's transaction items might need cyclic dependency and it doesn't works standalone, thus it should be on the same domain. In the other hand, user and transaction shouldn't require cyclic dependency, transaction might have a user field in the struct, but user shouldn't have a slice of transaction field, therefore it should be on a separate domain.
