# gcache

gcache is a distributed cache and cache-filling library for golang.

# Usage and example log

```
❯ ./run.sh
2024/01/15 21:37:11 frontend server is running at http://localhost:9999
>>> start test
2024/01/15 21:37:13 http://localhost:8003 cache miss for local cache []
2024/01/15 21:37:13 http://localhost:8003 cache miss for local cache []
2024/01/15 21:37:13 http://localhost:8003 cache miss for local cache []
2024/01/15 21:37:13 INFO [Server http://localhost:8003] http://localhost:8003 pick peer http://localhost:8001
2024/01/15 21:37:13 INFO [Server http://localhost:8001]  method=GET path=/_gcache/scores/Jack
2024/01/15 21:37:13 http://localhost:8001 cache miss for local cache []
2024/01/15 21:37:13 INFO [Server http://localhost:8001] peer is self
2024/01/15 21:37:13 http://localhost:8001 fetch from getter func []
2024/01/15 21:37:13 http://localhost:8001 sync value to local cache []
2024/01/15 21:37:13 http://localhost:8001 returning value: [589]
client sending out value:"589"
2024/01/15 21:37:13 http://localhost:8003 sync value to local cache []
2024/01/15 21:37:13 http://localhost:8003 returning value: [589]
2024/01/15 21:37:13 http://localhost:8003 returning value: [589]
2024/01/15 21:37:13 http://localhost:8003 returning value: [589]
{"data":"NTg5","error":null}
{"data":"NTg5","error":null}
{"data":"NTg5","error":null}
```
