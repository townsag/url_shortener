todo:
- generate unique short urls
    - how do we garuantee that concurrently generated short urls are unique:
        - pre generate short url codes
        - partition the namespace of short url codes by the number of servers
        - use insert into on clonflict to prevent concurrent updates
            - this is the simplest solution, use this with random generation
- store mapping between short urls and long urls
- serve short url redirects quickly

- when making the dockerfile for url shortener, do I have to copy in the url_shortener
  directory name for it to build properly?