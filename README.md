Permify Tester
==============

This is a simple test setup to validate the edges of permify and to ensure that I understand how it responds to certain pressure.


## **Building**
```
go mod tidy
go mod download
go build -o tester cmd/*
```

## Pressure Testing

1. First start permify in a console (watching the logs)
```
$ docker compose stop && docker compose rm -f && docker compose up -d && docker logs -f permify-repro-permify-1
```
2. Next run the tester playing with rate limits (getting them at or above the server), and count of relationship sets and number of times we create or delete them.
```
$ ./tester -iterations 10 --count 100 --rate-limit 100
```

## Dev Notes
1. The Client is a pure HTTP client, which has a built in rate limiter 
2. The schema is one from the permify examples - we create one of each, then delete them in reverse order
3. Modify the `server-rate-limit` in the [docker-compose.yml](./docker-compose.yml) to adjust where things should fail from the tester side.
4. Note the postgres has a default connection max for users of 100. This is adjustable

