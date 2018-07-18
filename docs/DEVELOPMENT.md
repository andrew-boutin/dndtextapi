# Development

## Manual Endpoint Testing

Most of the endpoints require headers like `CONTENT-TYPE` and/or `ACCEPT`. Also, many of the endpoints are for operations other than `GET`.This makes hitting the endpoints directly mostly unusable directly through the browser. On top of this, most endpoints require authentication which requires sending a cookie along with your requests so things like `curl` and standard `Postman` can't be used out of the box either.

However, you can a combintation of `Chrome`, `Postman`, and `Interceptor` to be able to easily make REST API requests.

### Installation and Setup

Install the [`Interceptor Chrome Extension`](https://chrome.google.com/webstore/detail/postman-interceptor/aicmkgpgakddgnaphhhpliifpcfhicfo?hl=en-US).

Open the extension and turn `Request Capture` on.

Install the [`Postman Chrome App`](https://chrome.google.com/webstore/detail/postman/fhbjgbiflinjbdggehcddcbncdddomop/related) (not the Postman standalone application).

Open the `Postman Chrome App`.

Enable `Postman Interceptor` in the `Postman Chrome App` by clicking on the icon in the tool bar (the icon looks like a satellite and if you hover it will say "Postman Interceptor") and then clicking on the slide button to enable it.

### Login & Requests

In `Chrome`, browse to `localhost:8080/login`. If you haven't logged into Google for the `dndtextapi` application lately you should be redirected to a Google login page. Log in. After you log in (or if you were still authenticated) you'll be redirected to a blank page with a URL that looks something like `http://localhost:8080/callback?state=<some_state>&code=<some_code>`.

At this point you have authenticated with Google, the app successfully loaded your User, and a session was initiated for you. Now you need to get your cookie to make new requests. In `Chrome` browse to `localhost:8080/channels`. This will force your browser to send a GET request to the REST API using the cookie from your session.

In the `Postman Chrome App` look in the History sidebar and find the `localhost:8080/login` request and click on it. This should open up that particular request in the main part of the screen. Click on the "Headers" tab and find the "Cookie" header. Copy the Value which should look something like `dndtextapisession=<some_stuff_here>`.

Create a new request in the `Postman Chrome App` (either use the `New` button or click the `+` icon next to the request tabs). Set up the request to be a `GET` on `localhost:8080/channels` and create a header of "Key" `Accept` with "Value" `application/json`. Next, add another header with "Key" `Cookie` and "Value" using the copied cookie value from the previous step. Issue the request (click `Send`) and you should get a valid response.

Using the same cookie header, you'll be able to make other requests to different paths with different operations and various other headers.

### Workarounds

If you can successfully go through the `/login` process, but still get 401 with logs similiar to

```bash
time="2018-07-18T17:50:57Z" level=error msg="No session data found denying access."
[sessions] ERROR! securecookie: the value is not valid
```

then you can try to clear cookies from your browser and re login. I've ran into this after restarting the server.

## Containers

Execute `docker ps -a` to see a list of Docker containers. If you've already ran `make up` you should see something similar to:

```bash
18:53 $ docker ps -a
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                   PORTS                    NAMES
6ce9399cdfee        dndtextapi_app      "./wait-for-it.sh db…"   6 minutes ago       Up 6 minutes             0.0.0.0:8080->8080/tcp   dndtextapi_app_1
540550a797fe        postgres            "docker-entrypoint.s…"   6 minutes ago       Up 6 minutes             0.0.0.0:5432->5432/tcp   dndtextapi_db_1
```

This shows the Go REST API `dndtextapi_app_1`, available on port 8080, and the Postgresql database `dndtextapi_db_1`.

## Logs

Get the container name or id from `docker ps -a` and then run `docker logs -f dndtextapi_app_1`. In this example `dndtextapi_app_1` is the name of the Go REST API Docker container that's running. The output will look like:

```bash
19:02 $ docker logs -f dndtextapi_app_1
wait-for-it.sh: waiting 15 seconds for db:5432
wait-for-it.sh: db:5432 is available after 2 seconds
time="2018-07-17T22:53:06Z" level=info msg="backend config type is postgres"
...
```

The `-f` in the command "tails" the logs so you will be able to see new output as it comes in. This is useful to tail the logs while testing out the API so you can see log messages from the requests you make.