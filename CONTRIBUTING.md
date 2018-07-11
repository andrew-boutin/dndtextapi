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

Can `curl` the available endpoints for manual testing:

```bash
# -- Channels --
# Create Channel POST /channels
curl -d '{"name":"my name", "description":"my description", "ownerid":1, "isprivate":false, "dmid":1, "users":[{"id":1}]}' -H "Content-Type: application/json" -H "Accept: application/json" -X POST localhost:8080/channels

# Get Channels GET /channels
curl -H "Accept: application/json" -X GET localhost:8080/channels

# Get Channel GET /channels/id
curl -H "Accept: application/json" -X GET localhost:8080/channels/1

# Delete Channel DELETE /channels/id
curl -X DELETE localhost:8080/channels/1

# Update Channel PUT /channels/id
curl -d '{"name":"my name updated", "description":"my description updated", "ownerid":1, "isprivate":true, "dmid":1, "users":[{"id":1}]}' -H "Content-Type: application/json" -H "Accept: application/json" -X PUT localhost:8080/channels/1

# -- Messages --
# Create Message POST /messages
curl -d '{"userid":1, "channelid":1,"content":"some new content in here"}' -H "Content-Type: application/json" -H "Accept: application/json" -X POSTlocalhost:8080/messages

# Get Messages for Channel GET /messages?channelID=id
curl -H "Accept: application/json" -X GET localhost:8080/messages?channelID=1

# Get Message GET /messages/id
curl -H "Accept: application/json" -X GET localhost:8080/messages/1

# Delete Message DELETE /messages/id
curl -X DELETE localhost:8080/messages/1

# Update Message PUT /messages/id
curl -d '{"content":"some updated content"}' -H "Content-Type: application/json" -H "Accept: application/json" -X PUT localhost:8080/messages/1
```