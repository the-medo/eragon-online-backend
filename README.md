
### Go backend masterclass
Golang, Postgres, Redis, Gin, gRPC, Docker, Kubernetes, AWS, CI/CD

https://github.com/techschool/simplebank/blob/master/README.md


## Tool setup on windows
### Using make commands on windows

https://gnuwin32.sourceforge.net/packages/make.htm

1. https://sourceforge.net/projects/gnuwin32/
2. install
3. add to path system variables

### Install scoop to use "migrate":
https://scoop.sh/
```
> Set-ExecutionPolicy RemoteSigned -Scope CurrentUser # Optional: Needed to run a remote script the first time
> irm get.scoop.sh | iex
```

https://github.com/ScoopInstaller/Scoop#readme


### Install "migrate" through scoop
` scoop install migrate `
- add migrate to path system variables (in ~\scoop\apps\migrate\ [version])

### Running sqlc generate on "windows"
`docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate`

---

### postgres driver
`go get github.com/lib/pq`
`go get github.com/stretchr/testify`

---
### Create migration

`migrate create -ext sql -dir db/migration -seq migration_name`

---
### Install jq

https://stedolan.github.io/jq/download/

---
### Install protobuf
https://grpc.io/docs/languages/go/quickstart/

Latest release:
https://github.com/protocolbuffers/protobuf/releases

```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```


---
## Workflow
1. **create new migration** `migrate create -ext sql -dir db/migration -seq migration_name`
2. **edit migration up/down sql files**
3. **run migrations** - `make migrateup`
4. **create new query files in** - `db/query`
5. **SQLC - generate sql.go file** `make sqlc-generate`
6. **run mockgen** - `make mock`
7. **run tests** - `make test`