# Integration Tests

The integration tests test out the integration between the API and database.

This is done by compose starting up two services: `mockserver` and `inttests`. `mockserver` is a [containerized mock server](https://github.com/jamesdbloom/mockserver) that gets configured to respond to specific requests with specific responses along with verifying that specific calls were made. This allows only the `app` and `db` to be tested by mocking out all external dependencies such as Google oauth2 authentication. `inttests` is a Python3 container that executes all of the python unit tests found in the `/inttests` directory. The python tests primarily utilize the `requests` package to excersise the `app` routes. However, it also has a [python mock-server client](https://github.com/internap/python-mockserver-client) installed that allows the integration tests to easily set up how the external dependencies should react. The `/inttests/requirements.txt` file determines what gets installed in the python container. The config file `config-inttest.yml` is used for the tests. To run the integration tests:

Running the tests requires that the API and database already be running locally (`make up`). The command to start the tests won't automatically start up the app and database to prevent accidentally blowing away data.

To start the tests:

    make inttests

Test output will look something like:

```bash
Successfully built 4d7398d76ca4
Successfully tagged dndtextapi_inttest:latest
Creating dndtextapi_inttest_1 ... done
Attaching to dndtextapi_inttest_1
inttest_1  | ....
inttest_1  | ----------------------------------------------------------------------
inttest_1  | Ran 4 tests in 0.024s
inttest_1  |
inttest_1  | OK
dndtextapi_inttest_1 exited with code 0
```

The mock-server allows you to see what was configured, requests received, what configuration the requests matched up to, and responses sent back. Simply view the container logs:

    docker logs -f dndtextapi_mockserver_1

The logs should look something like:

```bash
2018-07-24 23:25:02,332 INFO o.m.m.HttpStateHandler request:

        {
          "method" : "GET",
          "path" : "/oauth2/v2/userinfo",
          "queryStringParameters" : {
            "access_token" : [ "myfakeaccesstoken" ]
          },
          "headers" : {
            "content-length" : [ "0" ],
            "User-Agent" : [ "Go-http-client/1.1" ],
            "Host" : [ "mockserver:1080" ],
            "Accept-Encoding" : [ "gzip" ]
          },
          "keepAlive" : true,
          "secure" : false
        }

 matched expectation:

        {
          "method" : "GET",
          "path" : "/oauth2/v2/userinfo"
        }

2018-07-24 23:25:02,342 INFO o.m.m.HttpStateHandler returning response:

        {
          "statusCode" : 200,
          "headers" : {
            "connection" : [ "keep-alive" ]
          },
          "body" : "{\"email\": \"regularuser@fake.com\"}"
        }
```