# CONTRIBUTING

Contributions are always welcome, no matter how large or small. Before contributing,
please read the [code of conduct](CODE_OF_CONDUCT.md).

## Setup

> Install Go 1.11.1
> Install Docker to run tests

Gothic uses the Go Modules support built into Go 1.11 to build. The easiest is to clone Gothic in a directory outside of GOPATH, as in the following example:

```sh
$ git clone https://github.com/jrapoport/gothic
$ cd gothic
$ make deps
```

## Building

```sh
$ make build
```

## Running database migrations for development and testing

- Make sure your database can be accessed with user `root` without a password.
- Alternatively, you can start MySQL inside a docker container running `hack/mysqld.sh`.

### Migrations for development

```sh
$ make migrate_dev
```

### Migrations for testing

```sh
$ make migrate_test
```

## Testing

```sh
$ make test
```

## Pull Requests

We actively welcome your pull requests.

1. Fork the repo and create your branch from `master`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.

## License

By contributing to Gothic CMS, you agree that your contributions will be licensed
under its [MIT license](LICENSE).
