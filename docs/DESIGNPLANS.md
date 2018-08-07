# Design Plans

Some of the stuff that still has to be done and notes about how to do some of it.

## Features

- Add user_id to message? Would simplify a lot of logic. Don't allow update to this field.
- Characters managed under `/characters` and use query params to specify either a user or channel
- **User logout: quick attempt ran into issue where session would just get re-created on next api request
- *Application authn - send messages on behalf of a user (ex: Slack bot) - new use cases (won't need every route)
- Summary on get all vs full on get single
- Message story "sub-types" (dice roll, emote, talk, DM, adventure goal/topic change, etc..)
- Channel notes, inventory, etc.
- Transactions per route for rollbacks
- Swagger spec
- DB Migrations - for CD - separate db user
- Resource links: Self links. Collection links. (HAL) maybe https://github.com/pmoule/go2hal
- DMID still exists in channel. -1 anyone can talk as DM. DM still needs to have a character in the channel.
- Owner doesn't have to have a character - needs a character to send story messages though.
- Users can have multiple Characters in a Channel.

## Fix Ups

- Int tests for delete character, channel, and user. Verify dependent objects actually get deleted.
- QueryParamExtractor no error. QueryParamExtractorRequired error.
- Need an err msg somewhere when a container fails so example int tests in travis can easily tell why the app didn't start
- app takes a while to fully come up now may be related to govendor cmd change - may be able to add another step to dockerfile - https://github.com/kardianos/govendor/blob/master/doc/faq.md
- Get around having to rebuild docker images (map volume on startup or something etc.)
- Context err body in JSON format. One liner?
- Missing unit tests. A single test file in each package should get code coverage to report accurately.
- Add more middleware logging
- 404 consistent handling. Remove duplicate checks. Make sure everywhere supports it.
- Notes about connecting and inspecting the database.
- Restrict permissions for db api user.
- Improve error handling (500s probably should only log the error and not expose the underlying problem.)
- Improve packages docs - only list ones that have info that should be shared
- Hide User.IsAdmin from non /admin routes
- `wait-for-it.sh` in single place.
- Separate file for sample data
- Probably shouldn't need GOVENDOR_PATH and GOLINT_PATH
- Update docs that say the headers are required

### Bot OBO User

REST API would allow 3rd party apps (Slack/HipChat/Mupchat/etc.) to send Messages on behalf of Users.
Ex: Could have a SlackBot that can listen in Slack and send out requests to this API so
the final output could still be viewable on the site, but "gameplay" is in Slack.