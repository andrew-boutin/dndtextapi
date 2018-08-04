# Design Plans

Some of the stuff that still has to be done and notes about how to do some of it.

## Features

- **User logout: quick attempt ran into issue where session would just get re-created on next api request
- *Application authn - send messages on behalf of a user (ex: Slack bot) - new use cases (won't need every route)
- Summary on get all vs full on get single
- Message story "sub-types" (dice roll, emote, talk, DM, adventure goal/topic change, etc..)
- User/channel characters
- Channel notes, inventory, etc.
- Transactions per route for rollbacks
- Swagger spec
- DB Migrations - for CD - separate db user
- Resource links: Self links. Collection links. (HAL) maybe https://github.com/pmoule/go2hal
- DMID still exists in channel. -1 anyone can talk as DM. DM still needs to have a character in the channel.
- Owner doesn't have to have a character - needs a character to send story messages though.
- Users can have multiple Characters in a Channel.

## Fix Ups

- Delete Channel has dependencies on Messages and Characters
- Delete User has dependencies on Messages and Characters`
- QueryParamExtractor no error. QueryParamExtractorRequired error.
- Need an err msg somewhere when a container fails so example int tests in travis can easily tell why the app didn't start
- app takes a while to fully come up now may be related to govendor cmd change - may be able to add another step to dockerfile - https://github.com/kardianos/govendor/blob/master/doc/faq.md
- Get around having to rebuild docker images (map volume on startup or something etc.)
- 403 instead of 401 in many places
- Context err body in JSON format. One liner?
- Missing unit tests. A single test file in each package should get code coverage to report accurately.
- Add a lot more logging (especially in all middleware right before aborting request w/ only status)
- 404 consistent handling. Remove duplicate checks. Make sure everywhere supports it.
- Notes about connecting and inspecting the database.
- Restrict permissions for db api user.
- Improve error handling (500s probably should only log the error and not expose the underlying problem.)
- Improve packages docs - only list ones that have info that should be shared
- Hide User.IsAdmin from non /admin routes
- `wait-for-it.sh` in single place.
- db column names use `_` between words
- Separate file for sample data
- Probably shouldn't need GOVENDOR_PATH and GOLINT_PATH
- Update docs that say the headers are required

### Admin routes session issue

TODO: Authn and then go to /admin route in browser and observe server response.
TODO: Potentially related to non variable cookie initializer
TODO: Potentially related to copy/pasting cookie without path info etc.

Set up with `RegisterAdminRoutes(authorized)`
and admin route `g.GET("/test", RequireAdminHandler, RequiredHeadersMiddleware(acceptHeader), AdminGetChannels)` works fine
while admin route `g.GET("/admin/users", RequireAdminHandler, RequiredHeadersMiddleware(acceptHeader), AdminGetUsers)` gets 401

```bash
[sessions] ERROR! securecookie: the value is not valid
time="2018-07-20T01:31:02Z" level=error msg="No session data found denying access."
[GIN] 2018/07/20 - 01:31:02 | 401 |    6.550323ms |      172.19.0.1 | GET      /admin/users
```

and you end up with 3 cookies in the response. One for `/`, `/admin`, and `/admin/users`.

Set up with

```go
admin := authorized.Group("/")
admin.Use(RequireAdminHandler)
RegisterAdminRoutes(admin)
```

while admin route `g.GET("/test", RequiredHeadersMiddleware(acceptHeader), AdminGetChannels)` works fine
and admin route `g.GET("/admin/messages", RequiredHeadersMiddleware(acceptHeader), AdminGetMessages)` gets 401

```bash
[sessions] ERROR! securecookie: the value is not valid
time="2018-07-20T01:43:01Z" level=error msg="No session data found denying access."
[GIN] 2018/07/20 - 01:43:01 | 401 |       790.1µs |      172.19.0.1 | GET      /admin/users
```

and you end up with 3 cookies in the response. One for `/`, `/admin`, and `/admin/users`.

Problem appears to be related to routes that have multiple slashes (aside from for /:id). So routes that are
`/...` work while `/.../...` have session issues.

### Bot OBO User

REST API would allow 3rd party apps (Slack/HipChat/Mupchat/etc.) to send Messages on behalf of Users.
Ex: Could have a SlackBot that can listen in Slack and send out requests to this API so
the final output could still be viewable on the site, but "gameplay" is in Slack.