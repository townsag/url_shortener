## Purpose:
- this document serves as a record of the software design decisions that have been made in this project so that in the future I make look back at it and understand what I was thinking in the past
    - I also find that writing down the decision elucidates the decision process 

## Caching:
- use write around LRU caching
    - LRU:
        - use lru caching because it is simple and because we cannot accurately guess when url mapping will be accessed
        - will consider using lfu instead of lru caching
    - write around caching:
        - writes go to the database and invalidate the value in the cache for that read
        - reads read from the cache
            - if the cache is cold then reads read from the database and write the fresh value from the database to the cache
    - write around caching simplifies the consistency model and the way that the application interacts with the cache
        - don't have to worry about drift between database and cache as I would with write back caching
        - don't have to write custom caching logic like I would with write through

## Testing:
- use testcontainers to find a middle-ground between unit and integration testing
    - prevent the work necessary to mock out the redis / postgres clients
    - tests working with the actual postgres container will be more meaningful

## Observability: (8/27/25)
- use Opentelemetry for metrics:
    - opentelemetry metrics has not been stable for very long
    - opentelemetry sdk natively supports exporting metrics to prometheus
        - how does this work, is it push based?
    - opentelemetry has first class support in grafana tempo
    - grafana tempo has metric <-> trace collation feature using exemplars
        - really like the idea of associating traces with metrics
- use opentelemetry for traces:
    - I do not know of a better solution
    - use grafana Tempo oss as backend for traces
        - has aforementioned integration with grafana for metrics