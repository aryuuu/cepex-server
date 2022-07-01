# cepex-server

Multiplayer online 100 card game

### Prerequisite

- go ^1.15
- docker
- docker-compose

### How to run

- Copy `env.example` to `.env`

```shell
cp env.example env
```

- run the service

```shell
make run
```

- or run the service in development mode with live-reloading

```shell
make dev-air
```

### Run the test

```shell
make test
```
