# Design

## Objects

### Users

### Channels

Channels have a *visibility* which is either *public* or *private*. Public means anyone can view the story Messages in the Channel and also see the Channel details. The meta Messages aren't available even though the Channel is public. Private means only Users who are members of the Channel can view any of the Messages in the Channel and the Channel details.

### Messages

Messages have a *msgType* which is either *meta* or *story*. Meta Messages are anything that isn't considered part of the final output of the adventure such as Channel members talking Meta, Channel notifications, etc. Story Messages are all of the "in game" Messages such as character actions, characters speaking, dice rolls, DM output, etc. that make up the actual adventure story. Meta messages are only ever available to Users who are members of the Channel. Story Messages can be visible to other Users depending on the visbility.

## Authentication

... TODO: fill in [golang oauth2](https://github.com/golang/oauth2/) ...

## Endpoints

Channel Routes

- Get Channels GET /channels
  - Optional query param level=owner|member
- Get Channel GET /channels/id
- Create Channel POST /channels
- Delete Channel DELETE /channels/id
- Update Channel PUT /channels/id

Message Routes

- Get Messages for Channel GET /messages?channelID=id
  - Optional query param msgType=meta|story
- Get Message GET /messages/id
- Create Message POST /messages
- Delete Message DELETE /messages/id
- Update Message PUT /messages/id

User Routes

- Get Users for Channel GET /users?channelID=id
- Update User PUT /users/id
- Delete User DELETE /users/id

## Use Cases

### Anonymous Users

Anonymous User wants to get all public Channels. TODO: 4

- GET /channels TODO: 4

Anonymous User wants to get a single public Channel. TODO: 4

- GET /channel/id TODO: 4

Anonymous User wants to get the story Messages from a public Channel. TODO: 4

- GET /messages?channelID=id TODO: 4

Anonymous User wants to create an account.

- GET /login

Anonymous User wants to sign in.

- GET /login

### Authenticated Users

User wants to get their account info.

- GET /users/id

User wants to update their account info.

- PUT /users/id

User wants to delete their account.

- DELETE /users/id

User wants to get all public Channels.

- GET /channels

User wants to get a single public Channel.

- GET /channels/id

User wants to get story Messages from public Channel.

- GET /messages?channelID=id?msgType=story

User wants to sign out. TODO: 1

- GET /logout TODO: 1

## Channel Members

User wants to get all of the Channels they have access to.

- GET /channels

User wants to get all of the Channels they are a member of.

- GET /channels?level=member

User wants to get a single Channel they are a member of.

- GET /channels/id

User wants to get story Messages from a channel they're a member of.

- GET /messages?channelID=id&msgType=story

User wants to see all of the messages of a channel they're a member of.

- GET /messages?channelID=id

User wants to create a Messsage in the Channel.

- POST /messages

User wants to edit their Message in the Channel.

- PUT /messages/id

User wants to delete their Message in the Channel they're a member of.

- DELETE /messages/id

User wants to accept an invitation to join a Channel. TODO: 3

User wants to leave a Channel they're a member of. TODO: 3

User wants to list all other Users in a Channel they're a member of.

- GET /users?channelID=id

### Channel Owners

User wants to get all Channels that they're the owner of.

- GET /channels?level=owner

User wants to create a channel.

- POST /channels

User wants to update their channel data.

- PUT /channels/id

User wants to add someone to their channel. TODO: 2

User wants to remove someone from their channel. TODO: 2 (do not let them remove themselves)

User wants to delete a Message in their Channel.

- DELETE /messages/id

User wants to delete their Channel.

- DELETE /channels/id

### Admin Users

TODO: 5