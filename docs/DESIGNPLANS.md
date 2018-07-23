# Design Plans

Some of the stuff that still has to be done and notes about how to do some of it.

## Wanted Features

- User last login
- User logout
- Owner invite user to/remove user from channel
- User accept invitation to/leave channel
- Summary on get all vs full on get single
- Message story "sub-types" (dice roll, emote, talk, DM, adventure goal/topic change, etc..)
- User/channel characters
- Channel notes, inventory, etc.
- Application authn - send messages on behalf of a user (ex: Slack bot) - new use cases?
- Transactions per route for rollbacks
- Integration tests that require authentication
- Swagger spec
- DB Migrations - for CD - separate db user
- CI
- Resource links: Self links. Collection links. (HAL) maybe https://github.com/pmoule/go2hal

## Fix Ups

- Context err body in JSON format. One liner?
- Missing unit tests. A single test file in each package should get code coverage to report accurately.
- Add a lot more logging (especially in all middleware right before aborting request)
- 404 consistent handling. Remove duplicate checks. Make sure everywhere supports it.
- Notes about connecting and inspecting the database.
- Restrict permissions for db api user.
- Improve error handling (500s probably should only log the error and not expose the underlying problem.)
- Make cmds don't go into vendor dir
- Improve packages docs - only list ones that have info that should be shared
- User bio - should it also have NOT NULL... could a User go in and modify bio to NULL?
- Hide User.IsAdmin from non /admin routes
- Beef up design doc (especially authentication)

### Admin routes session issue

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

## Notes

REST API would allow 3rd party apps (Slack/HipChat/Mupchat/etc.) to send Messages on behalf of Users.
Ex: Could have a SlackBot that can listen in Slack and send out requests to this API so
the final output could still be viewable on the site, but "gameplay" is in Slack.