# DnD Text REST API

DnD Text API is a Golang REST API for a DnD style text adventure site.

## Endpoints

```bash
# Create Channel POST /channels
curl -d '{"name":"my name", "description":"my description", "ownerid":1, "isprivate":false, "dmid":1, "users":[{"id":1}]}' -H "Content-Type: application/json" -H "Accept: application/json" -X POST localhost:8080/channels

# Get Channels GET /channels
curl -H "Accept: application/json" -X GET localhost:8080/channels

# Get Channel GET /channel/id
curl -H "Accept: application/json" -X GET localhost:8080/channels/1

# Delete Channel DELETE /channel/id
curl -X DELETE localhost:8080/channels/1

# Update Channel PUT /channel/id
curl -d '{"name":"my name updated", "description":"my description updated", "ownerid":1, "isprivate":true, "dmid":1, "users":[{"id":1}]}' -H "Content-Type: application/json" -H "Accept: application/json" -X PUT localhost:8080/channels/1
```

## Development

Repo primarily uses Golang and Postgresql.

Start the database and server:

    make

See the [CONTRIBUTING](CONTRIBUTING.md) guidelines for more info.

## License

Dnd Text API is under [LICENSE](LICENSE).

## Copyright

This repository is under [COPYRIGHT](COPYRIGHT.md).