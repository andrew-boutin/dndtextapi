# Integration Tests

The integration tests test out the integration between the API and database.

Running the tests require that the API and database be running locally (`make up`). The following command to start the tests won't automatically start up the app and database to prevent accidentally blowing away data. To start the tests:

    make inttests

Compose spins up a container that runs all of the tests in `dndtextapi/inttests`. This is a suite of Python3 tests primarily utilizing the `requests` package to excersise the API routes.

Test output will look something like

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