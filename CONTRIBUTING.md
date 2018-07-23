# Contributing

Set up the [prerequisites](docs/PREREQUISITES.md).

Check out [design](docs/DESIGN.md) to learn how the existing stuff works.

Check out [design plans](docs/DESIGNPLANS.md) to see what needs work.

Add yourself to [contributors](docs/CONTRIBUTORS.md) if you work on something.

## Tech Stack

Golang used for the REST API server. Check out the notable [packages](docs/PACKAGES.md).

Postgresql used for the DB.

Docker and docker-compose used for containerization.

Make for commands. See [`Makefile`](Makefile).

## Running

Running anything outside of Docker, such as unit tests, will require running `govendor install +local` first.

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