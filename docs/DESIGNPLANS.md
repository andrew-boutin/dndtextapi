# Design Plans

Some of the stuff that still has to be done and notes about how to do some of it.

## Features

- Functional tests
- Transactions per route for rollbacks
- DB Migrations - for CD
- Authn definitely
- Authz maybe (admin user)
- Resource links: Self links. Collection links.
- Swagger spec
- Glide

## Fix Ups

- One liner to abort context with new err and any previously defined errs. Current impl keeps rolling through the middleware stack.
- Channel ops work w/ 0 users defined. Currently gives SQL error "no values defined".
- Missing unit tests. A single test file in each package should get code coverage to report accurately.

## Notes

Channel types. Single channel or multiple?
Message types. DM, roll, player...

REST API would allow people to use Slack/HipChat/Mupchat/etc. instead of using the website UI.
Ex: Could have a SlackBot that can listen in Slack and send out requests to this API so
the final output could still be viewable on the site, but "gameplay" is in Slack.