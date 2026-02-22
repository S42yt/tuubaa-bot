package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/S42yt/tuubaa-bot/core"
	vembed "github.com/S42yt/tuubaa-bot/modules/config/embed"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func handleSetRole(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	data := i.ApplicationCommandData().Options[0]
	var roleKey string
	var targetRoleID string

	for _, opt := range data.Options {
		switch opt.Name {
		case "role":
			roleKey = opt.StringValue()
		case "target":
			if r := opt.RoleValue(s, i.GuildID); r != nil {
				targetRoleID = r.ID
			} else {
				targetRoleID = opt.StringValue()
			}
		}
	}

	if roleKey == "" || targetRoleID == "" {
		return respond(s, i, "Invalid arguments")
	}

	db := core.NewMongoHandler()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Connect(ctx); err != nil {
		return respond(s, i, fmt.Sprintf("Failed to connect to DB: %v", err))
	}
	defer db.Disconnect(ctx)

	coll := db.Collection("guild_configs")
	filter := bson.M{"guild_id": i.GuildID}
	update := bson.M{"$set": bson.M{fmt.Sprintf("roles.%s", roleKey): targetRoleID}}
	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return respond(s, i, fmt.Sprintf("Failed to save config: %v", err))
	}
	if res.MatchedCount == 0 {
		doc := bson.M{"guild_id": i.GuildID, "roles": bson.M{roleKey: targetRoleID}}
		if _, err := coll.InsertOne(ctx, doc); err != nil {
			return respond(s, i, fmt.Sprintf("Failed to create config: %v", err))
		}
	}

	resp := vembed.BuildRoleSetResponse(roleKey, targetRoleID, i.Member.User.Username)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: resp,
	})
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	data := &discordgo.InteractionResponseData{Content: content, Flags: discordgo.MessageFlagsEphemeral}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}
