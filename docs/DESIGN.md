# Design

## Objects

### Users

Users represent an end User of the system. They're tied to a login in an authentication system.

Users have an *isAdmin* flag. Admins can reach the `/admin/...` routes to administer almost every object except they can't administer other admin Users or create new admins. Adding or removing admins should be done directly on the database using the database user with elevated permissions.

### Channels

Channels are where stories and communication happen. When a User creates a Channel they become the owner of that Channel.

A User is considered to be "in Channel" if they own the Channel or have a Character in the Channel.

Channels have a *visibility* flag which is either *public* or *private*. Public means anyone can view the story Messages in the Channel and also see the Channel details. The meta Messages aren't available even though the Channel is public. Private means only Users who are members of the Channel can view any of the Messages in the Channel and the Channel details.

### Characters

Characters represent Users inside of Channels. A User can have multiple Characters in a single Channel if they want to. Characters allow Users to store information about who they are in that particular Channel so everyone else can easily reference that information.

Only Channel owners can create new Characters in their Channel - this is how they invite Users to join their Channels. They identify the User the Character is intended for and aren't allowed to set the Character's name. Then the User who now owns that Character can decide to either delete the Character (reject the invitation) or update the Character - here they're required to provide a name. A Character that has a name filled out shows that the User decided to join the Channel. Channel owners can also delete Characters in their Channel so they can remove Users if necessary. However, only the Character owner can update the Character.

### Messages

Messages are how everyone communicates with each other. They're tied to a specific Channel and Character. This means they're also tied to specific Users since the Character the Message is from is tied to a User.

Messages have a *msgType* which is either *meta* or *story*. Meta Messages are anything that isn't considered part of the final output of the adventure such as Channel members talking Meta, Channel notifications, etc. Story Messages are all of the "in game" Messages such as character actions, characters speaking, dice rolls, DM output, etc. that make up the actual adventure story. Meta messages are only ever available to Users who are members of the Channel. Story Messages can be visible to other Users depending on the visbility.

## Authentication

Authentication is integrated with Google using Oauth2. A User can navigate to /login where they will be redirected to a Google login page for this application. If they successfully authenticate with Google they'll be redirected back to the app at /callback. Here either a new User will be created in the database or their existing User will be loaded up (if they've logged in before). A session will be created when a User logs in. Subsequent requests can be made using the cookie created from the login process.

All routes, except for the /public endpoints, will first verify that their is an active session for the User that is attempting to access the routes. If there is then the User will be looked up and loaded into the context. If not, then access gets denied.

## Endpoints

TODO: Audit these

Anonymous Routes

- Get public Channels GET /public/channels
- Get public Channel GET /public/channels/:channelID
- Get story Messages from public Channel GET /public/channels/:channelID/messages

Authentication Routes

- Login GET /login
- Google authentication callback GET /callback

Channel Routes

- Get Channels GET /channels
  - Optional query param level=owner|member
- Get Channel GET /channels/id
- Create Channel POST /channels
- Delete Channel DELETE /channels/id
- Update Channel PUT /channels/id

Message Routes

- Get Messages for Channel GET /channels/:channelID/messages
  - Optional query param msgType=meta|story
- Get Message GET /channels/:channelID/messages/id
- Create Message POST /channels/:channelID/messages
- Delete Message DELETE /channels/:channelID/messages/id
- Update Message PUT /channels/:channelID/messages/id

User Routes

- Get Users for Channel GET /channels/:channelID/users
- Update User PUT /users/id
- Delete User DELETE /users/id

Character Routes

TODO:

Admin Routes TODO:

- Get all Channels GET /channels
- Get a Channel GET /channels/id
- Update a Channel PUT /channels/id
- Delete a Channel DELETE /channels/id

- Get all Messages in a Channel GET /channels/:channelID/messages
- Get a Message GET /messages/id
- Update a Message PUT /messages/id
- Delete a Message DELETE /messages/id

- Get all Users GET /users
- Get a User GET /users/id
- Update a User PUT /users/id
- Delete a User DELETE /users/id

- Get all Characters GET /channels/:channelID/characters
- Get a Character GET /characters/id
- Update a Character PUT /characters/id
- Delete a Character DELETE /characters/id

## Usecases

There is a list of [`use cases`](docs/USECASES.md) that describe who might want to do what and how they would do it.