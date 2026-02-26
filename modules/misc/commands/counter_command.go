package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/S42yt/tuubaa-bot/core"
	v2 "github.com/S42yt/tuubaa-bot/utils/embed"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CounterHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {

		var count int
		if g, err := s.GuildWithCounts(i.GuildID); err == nil && g != nil {
			count = g.ApproximateMemberCount
		}

		db := core.NewMongoHandler()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.Connect(ctx); err != nil {
			return respondCounter(s, i, count, false, fmt.Sprintf("DB error: %v", err))
		}
		defer db.Disconnect(ctx)

		coll := db.Collection("guild_configs")
		filter := bson.M{"guild_id": i.GuildID}
		var doc bson.M
		if err := coll.FindOne(ctx, filter).Decode(&doc); err != nil {
			return respondCounter(s, i, count, false, "counter channel not configured")
		}

		chID, _ := doc["counter_channel"].(string)
		newName := fmt.Sprintf("「👥」Kinder✩%d", count)
		if chID == "" {
			return respondCounter(s, i, count, false, "counter channel not configured")
		}

		if _, err := s.ChannelEdit(chID, &discordgo.ChannelEdit{
			Name: newName,
		}); err != nil {
			return respondCounter(s, i, count, false, fmt.Sprintf("failed to update channel: %v", err))
		}

		return respondCounter(s, i, count, true, fmt.Sprintf("updated <#%s>", chID))
	}
}

func respondCounter(s *discordgo.Session, i *discordgo.InteractionCreate, count int, updated bool, note string) error {
	title := v2.NewTextDisplayBuilder().SetContent("### Server Counter").Build()
	body := v2.NewTextDisplayBuilder().SetContent(fmt.Sprintf("Members: %d\n%s", count, note)).Build()
	accent := 0x2ecc71
	if !updated {
		accent = 0x992222
	}
	comp := v2.NewContainerBuilder().SetAccentColor(accent).
		AddComponent(title).
		AddComponent(body).
		Build()

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{comp},
			Flags:      discordgo.MessageFlagsIsComponentsV2 | discordgo.MessageFlagsEphemeral,
		},
	})
}
