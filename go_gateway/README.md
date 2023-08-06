

```bash
etcdctl put "/gateway/httpRouter/" "{\"id\":1,\"host\":[\"127.0.0.1\"],\"method\":[\"GET\"],\"location\":\"/hello\"}"

curl http://127.0.0.1:8080/hello
```

