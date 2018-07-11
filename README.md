# DnD Text REST API

DnD Text API is a Golang REST API for a DnD style text adventure site.

## Endpoints

Channel Routes

- Create Channel POST /channels
- Get Channels GET /channels
- Get Channel GET /channels/id
- Delete Channel DELETE /channels/id
- Update Channel PUT /channels/id

Message Routes

- Create Message POST /messages
- Get Messages for Channel GET /messages?channelID=id
- Get Message GET /messages/id
- Delete Message DELETE /messages/id
- Update Message PUT /messages/id

## Development

Repo primarily uses Golang and Postgresql.

Start the database and server:

    make

See the [CONTRIBUTING](CONTRIBUTING.md) guidelines for more info.

## License

Dnd Text API is under [LICENSE](LICENSE).

## Copyright

This repository is under [COPYRIGHT](COPYRIGHT.md).