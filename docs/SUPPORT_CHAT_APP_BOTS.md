# Support Chat App Bots

This API will support chat app bots to enable other messaging platforms (such as Slack, HipChat, Mupchat, etc) to send messages on behalf of users. For example, you would be able to have a Slack bot in your Slack workspace listening to messages in your Slack channels. When a user sends a message in the Slack channel the Slack bot would be able to make a request to this API in order to add that user's message to the channel in the API. This would allow users to turn a Slack channel into a place for them to have a DnD text adventure. The users would still be able to use the API and/or [`dndtextui`](https://github.com/mupchrch/dndtextui) as well.

## Steps

### 1. Create bot in API

Need to be able to manage bot data in the API. This will allow users to create bots. Will add `type Bot Struct`, matching database table `bots`, and new endpoints for bot management CRUDL `/bots` & `/bots/:botID`. Anyone can retrieve a single or list all bots so they can find ones to use. Only the user who created the bot, and admins, will be able to modify the bot data. A user would be able to issue the following to create a new bot for their Slack workspace

    POST /bots data={"workspace":"moneyinthebank","ownerid":1}

```go
type Bot Struct {
    ID          int       `json:"ID" db:"id"`
    Workspace   string    `json:"Workspace" db:"workspace"`
    OwnerID     int       `json:"OwnerID" db:"owner_id"`
    LastUpdated time.Time `json:"LastUpdated" db:"last_updated"`
    CreatedOn   time.Time `json:"CreatedOn" db:"created_on"`
}
```

```sql
CREATE TABLE bots (
    id bigserial primary key,
    owner_id bigserial references users(id),
    workspace varchar(200) NOT NULL,
    created_on timestamp default current_timestamp,
    last_updated timestamp default current_timestamp
);
```

For authentication purposes, described later, each bot will need a pair of client credentials generated. Will add `type BotClientCredentials Struct`, matching database table `bot_client_credentials`, and a new endpoint to retrieve a single set of credentials `/bots/:botID/creds`. Only the `Bot.OwnerID`, and admins, can retrieve the bot's credentials. A bot owner would be able to issue the following to get the credentials for their bot

    GET /bots/:id/creds

```go
type BotClientCredentials struct {
    BotID        int       `json:"BotID" db:"bot_id"`
    ClientID     string    `json:"ClientID" db:"client_id"`
    ClientSecret string    `json:"ClientSecret" db:"client_secret"`
    LastUpdated  time.Time `json:"LastUpdated" db:"last_updated"`
    CreatedOn    time.Time `json:"CreatedOn" db:"created_on"`
}
```

```sql
CREATE TABLE bot_client_credentials (
    id bigserial primary key,
    bot_id bigserial references bots(id),
    client_id varchar(200) NOT NULL,
    client_secret varchar(200) NOT NULL,
    created_on timestamp default current_timestamp,
    last_updated timestamp default current_timestamp
);
```

The client credentials for a bot should be automatically created when the bot is created so will need to add that logic into the bot create handler. The client credentials are kept as a separate object to make them more secure.

Will also need to add postgresql triggers for both bots and client credentials `last_updated` fields.

### 2. Retrieve bot client credentials

The user will need the client credentials for their bot in order to properly configure it. All they need to do is issue the following

    GET /bots/:id/creds

### 3. Add bot to Slack org

Follow [Slack Bot Getting Started](https://api.slack.com/bot-users#getting-started) to create the bot in your Slack workspace.

### 4. Set up bot and run

Add the client credentials config along with the config from Slack for the bot. Run the bot and verify the health check passes.

### 5. Add bot to channel in API

Channel owners in the API need a way to authorize a bot to send messages in their channel. Will do this by adding fields to the channel object: `Channel.BotID` and `Channel.BotChannel`. These will be optional and default to their empty values. The channel owner, and admins, can update the channel to add or remove a bot from it.

    PUT /channels/:id data={..., botid=1, botchannel=somechannelname}

Will add in validation logic on create and update to require that either both of these optional fields are blank or both are filled out. The `BotChannel` would be the name of the Slack channel. This will enable a single bot to send messages in a channel.

### 6. Add bot to corresponding Slack channel

Follow [Slack Bot Setting Up The Events API](https://api.slack.com/bot-users#setup-events-api) to add your Slack bot to channels you want and to also make sure that it has access to the messagse in that channel.

At this point the bot will start listening to all messages that get sent in the channel that it was added to. It will essentially ignore all traffic since no users will have registered themselves yet.

### 7. User allows Bot to send messages for them in API

Need to allow channel members to allow/disallow a bot from sending messages on behalf of their characters. Can do this by adding a field to the character object `Character.BotUsername`. This will be an optional field that defaults to blank. The user, and admins, can update their character to fill in `BotUsername` with their corresponding Slack username. This will enable the bot, that was added to the channel, to send messages from the Slack user identified from `BotUsername` for this character. The field can be updated back to blank to take away the bots authorization.

    PUT /channels/:id/characters/:id data{..., botusername=slackusername}

### 8. User in Slack channel uses cmd to register with the Bot

A Slack user can send a command to the bot to get themselves registered on that end. They will need to be signed into Slack using the username that was used on the `Character.BotUsername` in the API. Also, they should issue the command in the Slack channel that was used on the `Channel.BotChannel` in the API.

    @botboi register characterid=1

A similiar command will exist to remove themselves as well.

### 9. Bot parses register cmd and verifies data with API

The bot is already listening for messages in the channel at this time. It will grab the register command and form a request. It will authenticate against the API first. Then it will attempt to retrieve the character matching the input character ID using its Slack bot id, the Slack channel that the request is from, and the character ID. If the `GET` on the character is successful then the Slack user has wired up everything correctly. At this point the bot will add the Slack user as one that has registered and being listening to all messages they send in the channel. The bot would just remove the user from the registered mapping if they send the command to remove themselves.

### 10. Registered user sends message in Slack channel

At this point a user that has registered with the Slack bot can type any message into the Slack channel and the Slack bot will attempt to send the message to the API similarly to how it retrieved the user's character info. The Slack bot will have certain formats it expects the messages to be in so it knows if messages are for meta/story, actions, dice rolls, etc. If the formatting is wrong the Slack bot will respond to the user in Slack with an error message. If the formatting looks good the message will be sent to the API. The Slack bot will also inform the Slack channel if there are issues connecting to the API or with the message itself.

## Preventing User Hijacking

TODO: Still need to re-write this section.

- The Channel configures both `BotID` and `BotChannelName`. Also, `ChannelID` and `UserID` are already part of the Character.
- Users can only add/remove their own data. Channel owner and admins can also remove other users data.
- User would have to use the correct data on both sides where authentication is required for each system.
- User only registering on Slack side wouldn't matter since the API will check to make sure they actually have access to what they claim.
- When bot sends message to API the API can check to verify that the registration info the bot has matches what it has as well. Allow/deny requests with this.
- User only registering Slack/bot info in the API wouldn't do them any good since they would need to register with the bot through Slack to have messages sent so don't need to worry about restrictions on that end either.
- Assuming any authenticated Slack User in the Slack channel should be allowed to sign themselves up with the bot. Slack admins could restrict them from joining the channel or organization if they didn't want them to have access.

Tieing Slack channel name to API channel name prevents Slack users from signing up for different API channels that are intended for Slack channels that they may not have access to. This is because a Slack Bot could be in multiple Slack Channels and the User could only have access to one of them and sign up for any Bot linked Channel from that one.

The command only needs the corresponding character ID since the Slack bot id, channel name, and username are already known. The character ID will lead to the corresponding user and channel data as well since the user ID and channel ID are part of the character associated to the character ID. This means the Slack username on the character, Slack channel name on the channel, and the Slack bot on the channel can all be cross referenced.

## Bot Design

Use [Simple-Slack-Bot](https://github.com/GregHilston/Simple-Slack-Bot) to create a Slack bot. Create a new GitHub repo for this. API side won't care what type of bot it is so other frameworks could be used in the future.

Will add in config files and config loading to handle the client credentials and Slack authentication variables.

The bot on startup, and periodically while running, will authenticate with the API and retrieve its own bot data. This will verify that the bot config is correct and that there is connectivity.

Messages the bot received from registerd users will be checked to see if they meet formatting requirements to be a valid message for the API. If there are issues with the message then the Slack bot will send an error message to the Slack channel.

The Slack bot will let the Slack channel know about any connectivity issues sending messages along with any error messages from the API.

## Authentication

### Bot Authentication

Utilize [oauth2 client credentials flow](https://developer.foresee.com/docs/oauth2-client-credentials-flow) to allow bots to authenticate with the api. The earlier steps handle the creation, access, and config of the client credentials per bot. Can utilize the [golang oauth2 client credentials package](https://godoc.org/golang.org/x/oauth2/clientcredentials) when implementing this.

The bot will need to issue a POST to `/oauth/token?grant_type=client_credentials` with header `Authorization: Basic <Base 64 encoded value formed from client_id:client_secret>`. This would return an access token. The bot would then issue requests with header `Authorization: Bearer <access token>`.

The `AuthenticationMiddleware` will need to be updated to check and see if the `Authorization` header is filled out to determine if it's a bot attempting to authenticate. If not then process the request like before and check for cookies/sessions. If it's a bot then verify that the access token is valid. If valid then need to determine who the request is from, load up their user, and put that user in the `gin.Context`. From there the other middleware can operate the same as before.

TODO: Details about how the authentication middleware will verify all of the details. Slack bot is in that channel. Slack channels match up / bot can send messages in that channel. Slack usernames match up / bot can send OBO that character.

TODO: More token research:

- Can the access token be a JWT that we can include information such as bot id and user id that it's on behalf of? This would only work if the bot got a new token for every user it was sending messages for. Could use a query parameter potentially and utilize that in the `AuthenticationMiddleware`.
- Need some more specifics on issuing, managing, and authenticating tokens. Refresh tokens? Expired?

### User Authentication Changes

Can optionally change how users authenticate too. Stop using sessions/cookies and instead issue access tokens after the Google profile data has been retrieved. This would allow bots and users to authenticate in a similar manner and hopefully simplify things.

## Additional Info

There is the potential for bot utility commands in Slack that wouldn't get sent to the API as messages.

- @botboi list characters // List characters in channel
- @botboi whoami          // Show character/user info

These would let the users in Slack get information about their Dnd text adventure without having to nagivate to the site.
