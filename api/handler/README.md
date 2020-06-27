# handler

This package contains handler that handle http request for each endpoints.
Every handler should be simple and light, it's should only responsible for decoding request, calling business service, and encoding the response.

It's recommended to avoid implementing any business logic directly in handler, including writes to the database.
Even if the logic seems simple at the beginning, implementing the logic directly in the handler might trigger tech-debt when additional requirement comes and other engineer just added the implementation directly in handler without moving it to a specific service.
