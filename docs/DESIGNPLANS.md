# Design Plans

Some of the stuff that still has to be done and notes about how to do some of it.

## Wanted Features

- User banned flag (admin update)
- User last login
- User logout
- Beef up design doc (especially authentication)
- Owner invite user to/remove user from channel
- User accept invitation to/leave channel
- Admin routes (User flag)
- Summary on get all vs full on get single
- Message story "sub-types" (dice roll, emote, talk, DM, adventure goal/topic change, etc..)
- User/channel characters
- Channel goal/topic, notes, inventory, etc.
- Application authn - send messages on behalf of a user (ex: Slack bot) - new use cases?
- Transactions per route for rollbacks
- Functional tests
- Resource links: Self links. Collection links.
- Swagger spec
- DB Migrations - for CD - separate db user
- CI
- Authz (should be fine for a while with just admin flag)

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

## Notes

REST API would allow 3rd party apps (Slack/HipChat/Mupchat/etc.) to send Messages on behalf of Users.
Ex: Could have a SlackBot that can listen in Slack and send out requests to this API so
the final output could still be viewable on the site, but "gameplay" is in Slack.