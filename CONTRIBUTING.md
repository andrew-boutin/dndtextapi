# Contributing

Check out the notable [packages](docs/PACKAGES.md).

Add yourself to [contributors][docs/CONTRIBUTORS.md] if you work on something.

[Design plans][DESIGNPLANS.md] may have some info on what needs work.

## Running

Utilize `Makefile`.

    make up

    make down

## Testing

Can `curl` the available endpoints.

    curl -d '{"name":"my name", "description":"my description", "ownerid":1, "isprivate":false, "dmid":1}' -H "Content-Type: application/json" -X POST localhost:8080/channels

    curl -d -H "Accept:application/json" -X GET localhost:8080/channels

    curl -d -H "Accept:application/json" -X GET localhost:8080/channels/1

    curl -d -H "Accept:application/json" -X DELETE localhost:8080/channels/1

    curl -d '{"name":"my name updated", "description":"my description updated", "ownerid":1, "isprivate":true, "dmid":1}' -H "Content-Type: application/json" -X PUT localhost:8080/channels/1