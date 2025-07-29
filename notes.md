todo:
- generate unique short urls
    - how do we garuantee that concurrently generated short urls are unique:
        - pre generate short url codes
        - partition the namespace of short url codes by the number of servers
        - use insert into on conflict to prevent concurrent updates
            - this is the simplest solution, use this with random generation
- store mapping between short urls and long urls
- serve short url redirects quickly
- analytics:
    - how many times has that short url been clicked
    - how many times has that long url been clicked across short urls

- when making the dockerfile for url shortener, do I have to copy in the url_shortener
  directory name for it to build properly?

- context:
    - a way to pass request scoped variables accross api boundaries and between processes
    - can be used to store lifetime information about a task / request
    - https://go.dev/blog/context



- Next steps:
    - [x] understand what context is?
    - [x] add the database client to the application, pass it to the handlers
    - [x] add two routes:
        - GET /{short_url_id}
        - POST /register
    - [x] add database support to those two routes
    - [x] add health check
    - [ ] add caching
    - [ ] add horizontal scaling
    - [ ] add load balancing
    - [ ] add metrics
    - [ ] add log aggregation
    - [ ] add tracing with opentelemetry
    - [ ] add admin dashboard
        - analytics?
    - [ ] add frontend
    - [ ] add auth?
    - testing types
        - unit + integration tests with code coverage
        - formal specification
            - tla+ or P
        - property testing
            - testing/quick
        - code coverage guided property based testing
            - not yet implemented in golang