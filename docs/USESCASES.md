# Usecases

Who might want to do what and how they would do it.

## Anonymous Users

Anonymous User wants to get all public Channels.

- GET /public/channels

Anonymous User wants to get a single public Channel.

- GET /public/channel/id

Anonymous User wants to get the story Messages from a public Channel.

- GET /public/channels/:id/messages

Anonymous User wants to create an account.

- GET /login

Anonymous User wants to sign in.

- GET /login

## Authenticated Users

User wants to get their account info.

- GET /users/:id

User wants to update their account info.

- PUT /users/:id

User wants to delete their account.

- DELETE /users/:id

User wants to get all public Channels.

- GET /channels

User wants to get a single public Channel.

- GET /channels/:id

User wants to get story Messages from public Channel.

- GET /channels/:id/messages?msgType=story

User wants to sign out. TODO: 1

- GET /logout TODO: 1

User wants to get all of their Characters. TODO:

User wants to get a single Character of theirs. TODO:

## Channel Members

User wants to get all of the Channels they have access to.

- GET /channels

User wants to get all of the Channels they are a member of.

- GET /channels?level=member

User wants to get a single Channel they are a member of.

- GET /channels/:id

User wants to get story Messages from a channel they're a member of.

- GET /channels/:id/messages?msgType=story

User wants to see all of the messages of a channel they're a member of.

- GET /channels/:id/messages

User wants to create a Messsage in the Channel.

- POST /channels/:id/messages

User wants to edit their Message in the Channel.

- PUT /channels/:id/messages/:id

User wants to delete their Message in the Channel they're a member of.

- DELETE /channels/:id/messages/:id

User wants to accept an invitation to join a Channel.

- PUT /channels/:id/characters/:id

User wants to leave a Channel they're a member of.

- DELETE /channels/:id/characters/:id

User wants to list all other Characters in a Channel they're a member of.

- GET /channels/:id/characters

## Channel Owners

User wants to get all Channels that they're the owner of.

- GET /channels?level=owner

User wants to create a channel.

- POST /channels

Owner wants to update their channel data.

- PUT /channels/:id

Owner wants to add someone to their channel.

- POST /channels/:id/characters

Owner wants to remove someone from their channel.

- DELETE /channels/:id/characters/:id

Owner wants to delete a Message in their Channel.

- DELETE /channels/:id/messages/:id

Owner wants to delete their Channel.

- DELETE /channels/:id

## Admin Users

- Admin wants to get all Users. TODO:

- Admin wants to get User. TODO:

- Admin wants to update User. TODO:

- Admin wants to delete User. TODO:

- Admin wants to get all Channels. TODO:

- Admin wants to get Channel. TODO:

- Admin wants to update Channel. TODO:

- Admin wants to delete Channel. TODO:

- Admin wants to remove a User from a Channel. TODO:

- Admin wants to get all Messages for Channel. TODO:

- Admin wants to get Message. TODO:

- Admin wants to update Message. TODO:

- Admin wants to delete Message. TODO: