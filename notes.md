todo:
- generate unique short urls
    - how do we garuantee that concurrently generated short urls are unique:
        - pre generate short url codes
        - partition the namespace of short url codes by the number of servers
        - use insert into on clonflict to prevent concurrent updates
            - this is the simplest solution, use this with random generation
- store mapping between short urls and long urls
- serve short url redirects quickly
- analytics:
    - how many times has that short url been clicked
    - how many times has that long url been clicked across short urls

- when making the dockerfile for url shortener, do I have to copy in the url_shortener
  directory name for it to build properly?

- contexts:
    - a way to pass request scoped state between functions

- Next steps:
    - understand what context is?
    - add the database client to the appication, pass it to the handlers
    - add two routes:
        - GET /{short_url_id}
        - POST /register
    - add database support to those two routes
    - add health check
    - add caching
    - add load balancing
    - add horizontal scaling
    - add metrics
    - add log agregation
    - add tracing with opentelemetry
    - add admin dashboard
        - analytics?
    - add frontend
    - add auth?