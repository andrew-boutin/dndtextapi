# Design Plans

Some of the stuff that still has to be done and notes about how to do some of it.

## Wanted Features

1. Authn
2. Transactions per route for rollbacks
3. Functional tests
4. Glide
5. Authz maybe (maybe just admin user flag w/ admin only routes)
6. Swagger spec
7. DB Migrations - for CD
8. Resource links: Self links. Collection links.

## Fix Ups

- Context err body in JSON format. One liner?
- Missing unit tests. A single test file in each package should get code coverage to report accurately.
- Membership tables unique ids to enforce single pairs
- Add a lot more logging
- 404 consistent handling. Remove duplicate checks. Make sure everywhere supports it.

## Notes

Story message sub-types: DM, roll, player...

REST API would allow people to use Slack/HipChat/Mupchat/etc. instead of using the website UI.
Ex: Could have a SlackBot that can listen in Slack and send out requests to this API so
the final output could still be viewable on the site, but "gameplay" is in Slack.