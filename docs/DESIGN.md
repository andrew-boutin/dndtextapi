# Design

## Objects

### Users

Users have an *isAdmin* flag. Admins can reach the `/admin/...` routes to perform admin operations on Channels, Messages, and Users. Admins can't administer other admin Users or create new admins.

Adding or removing admins should be done directly on the database using the database user with elevated permissions.

### Channels

Channels have a *visibility* which is either *public* or *private*. Public means anyone can view the story Messages in the Channel and also see the Channel details. The meta Messages aren't available even though the Channel is public. Private means only Users who are members of the Channel can view any of the Messages in the Channel and the Channel details.

### Messages

Messages have a *msgType* which is either *meta* or *story*. Meta Messages are anything that isn't considered part of the final output of the adventure such as Channel members talking Meta, Channel notifications, etc. Story Messages are all of the "in game" Messages such as character actions, characters speaking, dice rolls, DM output, etc. that make up the actual adventure story. Meta messages are only ever available to Users who are members of the Channel. Story Messages can be visible to other Users depending on the visbility.

## Authentication

... TODO: fill in [golang oauth2](https://github.com/golang/oauth2/) ...

> Tldr; set up Google cloud project (free) to get client creds, set up /login to redirect user to Google login w/ client creds, Google redirects user to /callback with an access code, I take access code and call Google API to get Google User profile, I then either create or load User mapped to unique Google email, create session for User, require other end points be able to look up active session before proceeding

## Endpoints

Anonymous Routes

- Get public Channels GET /public/channels
- Get public Channel GET /public/channels/id
- Get story Messages from public Channel GET /public/messages?channelID=id

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

Admin Routes TODO:

- Get all Channels GET /channels
- Get a Channel GET /channels/id
- Update a Channel PUT /channels/id
- Delete a Channel DELETE /channels/id

- Get all Messages in a Channel GET /messages?channelID=id
- Get a Message GET /messages/id
- Update a Message PUT /messages/id
- Delete a Message DELETE /messages/id

- Get all Users GET /users
- Get a User GET /users/id
- Update a User PUT /users/id
- Delete a User DELETE /users/id

## Usecases

There is a list of [`use cases`](docs/USECASES.md) that describe who might want to do what and how they would do it.