# little-sso | gRPC

little-sso - little auth service (gRPC).


## Installation

```bash
$ cd ./sso && go mod download
```

## Build

You need to fill the env-file with valid variables. Then, run to build the project:
```bash
$ cd ./sso && \
    go build -o=./tmp/bin/little_sso ./sso/cmd/sso
```
The build artifacts will be stored in the `./sso/tmp/bin/` directory.


## Development server

Run for a dev server. :
```bash
$ cd ./sso && \
    go build -o=./tmp/bin/little_sso ./sso/cmd/sso && \
    ./tmp/bin/little_sso
```
Navigate to `http://localhost:44044/`.


## Running functional tests

Run app and then run to execute the functional tests:
```bash
$ cd ./sso && \
    go test ./tests
```

## Quality control

```bash
$ cd ./sso && \
	go fmt ./... && \
	go mod tidy -v
```
