# Bot

REST API would allow 3rd party apps (Slack/HipChat/Mupchat/etc.) to send Messages on behalf of Users.
Ex: Could have a SlackBot that can listen in Slack and send out requests to this API so
the final output could still be viewable on the site, but "gameplay" is in Slack.

Use [Simple-Slack-Bot](https://github.com/GregHilston/Simple-Slack-Bot) to create a Slack bot. API side won't care what type of bot it is so other frameworks could be used in the future.

## Steps

Create Bot in API.

    POST /bots data={...}

Retrieve bot client credentials.

    GET /bots/:id/creds

Set up bot and run.

Add bot to Slack org (Slack docs for this).

Add Bot to Channel in API.

    PUT /channels/:id data={..., botid=1, botchannel=somechannelname}

Bot gets added to corresponding Slack channel (Slack docs for this).

User allows Bot to send Messages for them in the API.

    PUT /channels/:id/characters/:id data{..., botusername=slackusername}

User in Slack channel sends message to register themselves with the Bot.

    @botboi register characterid=1

Bot parses register command and verifies data with API. If all good then adds mapping data to start listening for Messages from the Slack User.

Bot listens to all traffic in the channel. Ignores messages from Slack users that aren't registered.

When a message is received from a registered user - parses the message (formatting specific to Slack) to determine what kind of message to send (story vs. meta, dice roll, DM, talking, etc.). Send the message to the API.

API receives Message from Bot. Authenticates Bot. Verifies User registration data. Adds Message if all good.

## Authentication

Need a way for the bot to authenticate with the api.

Client ID/Secret per bot created when Bot is created. Kept in a separate table `bot_client_credentials`. Accessed on different route only accessible to the Bot Owner and Admins. GET `/bots/:id/creds`. User adds these credentials to their Bot config.

https://godoc.org/golang.org/x/oauth2/clientcredentials API requires client id/secret for communication from the Bot. Client credentials potentially just POST to /oauth/token with form data `grant_type=client_credentials` (would also need to send over client id and client secret). Users would have to first go to `oauth/auth` like they do now (by login redirect).

Potentially stop using sessions and issue tokens to users as well (JWT). This way both users and bots would use tokens for requests which may simplify things.

Could have authentication middleware load up the user the bot is OBO into the context so the api calls could be handled the same as if the User initiated the request.

TODO: Need some more specifics on issuing, managing, and authenticating tokens. Need to use JWT to differentiate between Bots and Users, have things such as IDs, etc.?

## Registering Bot

Need to let a user register a Slack bot/Slack org. This also requires giving the Bot access to a Channel and removing the Bot from a Channel.

Add new endpoints for bot management CRUDL `/bots` & `/bots/:botID`

Add new struct and database table `bots` with schema that matches the `Bot` struct.

```go
type Bot Struct {
    ID int                // 1
    OrgName string        // moneyinthebank
    OwnerID int           // ID of User that created the bot
    LastUpdated time.Time
    CreatedOn time.Time
}
```

Add `BotID` and `BotChannelName` to the Channel struct. Both of these fields will be optional. This allows the Channel owner a way to add/remove a Bot from their Channel. Doesn't seem to make sense to need to allow multiple Slack bots to send messages to a single API Channel so the id on the Channel should be fine. Either both fields are empty or both fields are filled out.

Create struct for client credentials.

```go
type ClientCredentials struct {
    ClientID     string
    ClientSecret string
}
```

Create database table `bot_client_credentials` that matches the struct - with additional field for `BotID`.

Add endpoint `GET /bots/:id/creds` that retrieves the credentials for the given Bot.

Have new client credentials automatically created when a new Bot is created.

## Let Bot Message OBO User

User needs to be able to allow the bot to send messages OBO them.

- First register through API. Add string `BotUsername` to Character and make it optional. Blank means the User has not set up that Character to let a Bot send Messages for them. The Channel configures both `BotID` and `BotChannelName`. Also, `ChannelID` and `UserID` are already part of the Character.
- Register with a command in a Slack channel that the bot is present in. Command will require character id. The bot and channel name are already known. The Character in the API already has the Slack name. So if that Character defined the matching Slack username then everything matches up. Bot can then begin listening for any Messages they send.
- User would have to use the correct data on both sides where authentication is required for each system.
- Also need to handle User removing access. Commands on each side again.
- Users can only add/remove their own data.
- Channel owner and admins can also remove other users data.
- Assuming any authenticated Slack User in the Slack channel should be allowed to sign themselves up with the bot. Slack admins could restrict them from joining the channel or organization if they didn't want them to have access.
- User only registering Slack/bot info in the API wouldn't do them any good since they would need to register with the bot through Slack to have messages sent so don't need to worry about restrictions on that end either.
- When bot sends message to API the API can check to verify that the registration info the bot has matches what it has as well. Allow/deny requests with this.
- User only registering on Slack side wouldn't matter since the API will check to make sure they actually have access to what they claim.

Tieing Slack channel name to API channel name prevents Slack users from signing up for different API channels that are intended for Slack channels that they may not have access to. This is because a Slack Bot could be in multiple Slack Channels and the User could only have access to one of them and sign up for any Bot linked Channel from that one.

## Other

Potential for bot utility commands in Slack that wouldn't get sent to the API as messages.

- @botboi list characters // List characters in channel
- @botboi whoami          // Show character/user info
