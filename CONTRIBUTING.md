# Contributing

Set up the [prerequisites](docs/PREREQUISITES.md).

Check out [design](docs/DESIGN.md) to learn how the existing stuff works.

Check out [design plans](docs/DESIGNPLANS.md) to see what needs work.

Add yourself to [contributors](docs/CONTRIBUTORS.md) if you work on something.

## Tech Stack

Golang used for the REST API server. Check out the notable [packages](docs/PACKAGES.md).

Postgresql used for the DB.

Docker and docker-compose used for containerization.

Google for Oauth2 authentication.

Make for commands. See [`Makefile`](Makefile).

Python for containerized integration tests.

Mock-server for containerized mock authentication server.

Python mock-server client for setup between integration tests and mock authentication server.

## Running

Running anything outside of Docker, such as unit tests, will require running `make fetchdeps` first.

Start or stop the db and server:

    make up

    make down

Restart only the server (useful for code changes):

    make reapp

## Testing

Unit tests:

    make test

Integration tests (detailed info at [integration tests](docs/INTEGRATION_TESTS.md)):

    make inttests

Others:

    make fmt lint vet