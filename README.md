# Project prisoner-fencing

Simple game inspired by a Reddit post on r/gameideas.
Made with golang and svelte
Game is under development.

## Getting Started

### MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Live reload both front-end and back-end:

```bash
make watch-all
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
