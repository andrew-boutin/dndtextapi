# Contributing

Check out [Design][docs/DESIGN.md] to learn how the existing stuff works.

Check out [Design plans][docs/DESIGNPLANS.md] to see what needs work.

Add yourself to [contributors][docs/CONTRIBUTORS.md] if you work on something.

Set up the [prerequisites][docs/PREREQUISITES.md].

## Tech Stack

Golang used for the REST API server. Check out the notable [packages](docs/PACKAGES.md).

Postgresql used for the DB.

Docker and docker-compose used for containerization.

Make for commands. See [`Makefile`](Makefile).

## Running

Start or stop the db and server:

    make up

    make down

Restart only the server (useful for code changes):

    make reapp

## Testing

Unit tests:

    make test

Style:

    make lint

Can `curl` the available endpoints for manual testing.
