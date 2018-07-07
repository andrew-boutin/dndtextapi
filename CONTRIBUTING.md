# Contributing

Check out the notable [packages](#PACKAGES.md).

## Running

Utilize `Makefile`.

    make up

    make down

## Testing

Can `curl` the available endpoints.

    curl -d '{"name":"my name", "description":"my description", "ownerid":1, "isprivate":false, "dmid":1}' -H "Content-Type: application/json" -X POST localhost:8080/channels

    curl -d -H "Accept:application/json" -X GET localhost:8080/channels

    curl -d -H "Accept:application/json" -X GET localhost:8080/channels
