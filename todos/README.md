# todos

Contains domain related entities and business logic implementations. The business functionality should be exported using `Service` interface that contains necessary functions to work with the entity.

Every domain/client should have it's own testing package (`todostest`) that can be used to mock the functionality of this package, usualy generated using external tools like `mockery`.
