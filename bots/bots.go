// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package bots

import (
	"fmt"
	"time"
)

var (
	// ErrBotNotFound is the error to use when the bot is not found.
	ErrBotNotFound = fmt.Errorf("bot not found")

	// ErrBotCredsNotFound is the error to use when the bot creds are not found.
	ErrBotCredsNotFound = fmt.Errorf("bot client credentials not found")
)

// Bot holds all of the bot data
type Bot struct {
	ID          int       `json:"ID" db:"id"`
	Workspace   string    `json:"Workspace" db:"workspace"`
	OwnerID     int       `json:"OwnerID" db:"owner_id"`
	LastUpdated time.Time `json:"LastUpdated" db:"last_updated"`
	CreatedOn   time.Time `json:"CreatedOn" db:"created_on"`
}

// BotCollection is a collection of bots.
type BotCollection []*Bot

// BotClientCredentials holds all of the client credentials data for a single bot.
type BotClientCredentials struct {
	ID           int       `json:"ID" db:"id"`
	BotID        int       `json:"BotID" db:"bot_id"`
	ClientID     string    `json:"ClientID" db:"client_id"`
	ClientSecret string    `json:"ClientSecret" db:"client_secret"`
	LastUpdated  time.Time `json:"LastUpdated" db:"last_updated"`
	CreatedOn    time.Time `json:"CreatedOn" db:"created_on"`
}
