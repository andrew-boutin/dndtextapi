# Design Plans

Some of the stuff that still has to be done and notes about how to do some of it.

- unit tests
- functional tests
- DB Migrations
- Authn
- Authz

Only accept application/json & look for header.
Return content type application/json header.
Links. Self links. Collection links.
Swagger spec.

Channel types. Single channel or multiple?
Message types. DM, roll, player...

REST API would allow people to use Slack/HipChat/Mupchat/etc. instead of using the website UI.
Ex: Could have a SlackBot that can listen in Slack and send out requests to this API so
the final output could still be viewable on the site, but "gameplay" is in Slack.