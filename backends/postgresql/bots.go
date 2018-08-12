// Copyright (C) 2018, Baking Bits Studios - All Rights Reserved

package postgresql

import (
	sqlP "database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/dchest/uniuri"

	"github.com/andrew-boutin/dndtextapi/bots"
	log "github.com/sirupsen/logrus"
)

const (
	botsTable     = "bots"
	botsReturning = "RETURNING id, workspace, owner_id, created_on, last_updated"

	botsCredsTable    = "bot_client_credentials"
	botCredsReturning = "RETURNING id, bot_id, client_id, client_secret, created_on, last_updated"
)

var botColumns = []string{
	"id",
	"workspace",
	"owner_id",
	"created_on",
	"last_updated",
}

var botCredsColumns = []string{
	"id",
	"bot_id",
	"client_id",
	"client_secret",
	"created_on",
	"last_updated",
}

func init() {
	// Add the bot table name in front of the columms to avoid ambigious references.
	for i, col := range botColumns {
		botColumns[i] = fmt.Sprintf("%s.%s", botsTable, col)
	}
}

// GetAllBots retrieves all of the bots in the database.
func (b Backend) GetAllBots() (bots.BotCollection, error) {
	sql, args, err := PSQLBuilder().
		Select(botColumns...).
		From(botsTable).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build get all bots query.")
		return nil, err
	}

	rows, err := b.db.Queryx(sql, args...)
	if err != nil {
		log.WithError(err).Error("Failed to execute get all bots query.")
		return nil, err
	}

	outBots := make(bots.BotCollection, 0)
	for rows.Next() {
		var bot bots.Bot
		err = rows.StructScan(&bot)
		if err != nil {
			log.WithError(err).Error("Failed to load bot from get all bots query.")
			return nil, err
		}

		outBots = append(outBots, &bot)
	}

	return outBots, nil
}

// GetBot retrieves a single bot from the database using the input ID.
func (b Backend) GetBot(id int) (*bots.Bot, error) {
	bot := &bots.Bot{}
	wasFound, err := b.getSingle(id, botsTable, botColumns, bot)
	if err != nil {
		log.WithError(err).Error("Failed to execute get bot query.")
		return nil, err
	} else if !wasFound {
		return nil, bots.ErrBotNotFound
	}

	return bot, nil
}

// GetBotCreds retrieves the credentials for the bot matching the input ID.
func (b Backend) GetBotCreds(botID int) (*bots.BotClientCredentials, error) {
	sql, args, err := PSQLBuilder().
		Select(botCredsColumns...).
		From(botsCredsTable).
		Where(sq.Eq{"bot_id": botID}).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Issue building query for get bot creds.")
		return nil, err
	}

	botCreds := &bots.BotClientCredentials{}
	err = b.db.Get(botCreds, sql, args...)
	if err != nil {
		if err == sqlP.ErrNoRows {
			return nil, bots.ErrBotCredsNotFound
		}
		return nil, err
	}
	return botCreds, err
}

// DeleteBot deletes the bot matching the input ID.
func (b Backend) DeleteBot(id int) error {
	wasFound, err := b.deleteSingle(id, botsTable)
	if err != nil {
		log.WithError(err).Error("Failed to execute delete bot query.")
	} else if !wasFound {
		return bots.ErrBotNotFound
	}
	return err
}

// CreateBot creates a new bot using the input data.
func (b Backend) CreateBot(botData *bots.Bot) (*bots.Bot, error) {
	kvs := map[string]interface{}{
		"workspace": botData.Workspace,
		"owner_id":  botData.OwnerID,
	}

	newBot := &bots.Bot{}
	err := b.createSingle(botsTable, botsReturning, kvs, newBot)
	if err != nil {
		log.WithError(err).Error("Issue with create bot query.")
		return nil, err
	}

	return newBot, nil
}

// CreateBotCreds creates a new pair of credentials for the bot
// matching the input ID. The credentials are strored in a separate
// table than the bots.
func (b Backend) CreateBotCreds(bot *bots.Bot) (*bots.BotClientCredentials, error) {
	// TODO: Move client creds generation out into middle layer
	// TODO: Any requirements on these for length, formatting, etc.?
	clientID := uniuri.NewLen(8)
	clientSecret := uniuri.NewLen(24)

	kvs := map[string]interface{}{
		"bot_id":        bot.ID,
		"client_id":     clientID,
		"client_secret": clientSecret,
	}

	newBotCreds := &bots.BotClientCredentials{}
	err := b.createSingle(botsCredsTable, botCredsReturning, kvs, newBotCreds)
	if err != nil {
		log.WithError(err).Error("Issue with create bot credentials query.")
		return nil, err
	}

	return newBotCreds, nil
}

// UpdateBot updates the bot matching the input ID with the input data.
func (b Backend) UpdateBot(id int, botData *bots.Bot) (*bots.Bot, error) {
	setMap := map[string]interface{}{
		"workspace": botData.Workspace,
	}

	updatedBot := &bots.Bot{}
	wasFound, err := b.updateSingle(id, botsTable, botsReturning, setMap, updatedBot)
	if err != nil {
		log.WithError(err).Error("Issue with query for update bot.")
		return nil, err
	} else if !wasFound {
		return nil, bots.ErrBotNotFound
	}

	return updatedBot, nil
}

// DeleteBotCreds deletes the bot credentials that match the input bot ID.
func (b Backend) DeleteBotCreds(botID int) error {
	// TODO: Handle not found
	sql, args, err := PSQLBuilder().
		Delete(botsCredsTable).
		Where(sq.Eq{"bot_id": botID}).
		ToSql()
	if err != nil {
		log.WithError(err).Error("Failed to build query for delete bot credentials.")
		return err
	}

	_, err = b.db.Exec(sql, args...)
	if err != nil {
		log.WithError(err).Error("Failed to execute delete bot credentials query.")
	}

	return err
}
