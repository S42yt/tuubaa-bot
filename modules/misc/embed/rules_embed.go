package embed

import (
	v2 "github.com/S42yt/tuubaa-bot/utils/embed"
	ulog "github.com/S42yt/tuubaa-bot/utils/logger"
	"github.com/bwmarrin/discordgo"
)

func BuildRuleEmbed(ruleText string) []discordgo.MessageComponent {
	title := v2.NewTextDisplayBuilder().SetContent("### Regeln")
	content := v2.NewTextDisplayBuilder().SetContent(ruleText).Build()

	comp := v2.NewContainerBuilder().
		SetAccentColor(0x992222).
		AddComponent(title).
		AddComponent(content).
		Build()

	return []discordgo.MessageComponent{comp}
}

func BuildSetRuleEmbed(ruleContent string, s *discordgo.Session, i *discordgo.Interaction) []discordgo.MessageComponent {
	guild, err := s.State.Guild(i.GuildID)
	if err != nil {
		ulog.Warn("Failed to get guild from state for rules embed: %v", err)
		guild = &discordgo.Guild{ID: i.GuildID}
	}

	titleText := v2.NewTextDisplayBuilder().SetContent("### [Regelwerk](https://discord.com/terms)").Build()
	thumbnail := v2.NewThumbnailBuilder().SetURL(guild.IconURL("512")).Build()
	body := v2.NewTextDisplayBuilder().SetContent(ruleContent).Build()

	headerSection := v2.NewSectionBuilder().
		AddComponent(titleText).
		AddComponent(body).
		SetAccessory(thumbnail).
		Build()

	footer := v2.NewTextDisplayBuilder().SetContent("-# tuubaa").Build()

	comp := v2.NewContainerBuilder().
		SetAccentColor(0x992222).
		AddComponent(headerSection).
		AddComponent(footer).
		Build()

	return []discordgo.MessageComponent{comp}
}

func BuildRuleSuccessResponse() *discordgo.InteractionResponseData {
	content := v2.NewTextDisplayBuilder().SetContent("Regel wurde gesendet!").Build()

	comp := v2.NewContainerBuilder().
		SetAccentColor(0x99EE99).
		AddComponent(content).
		Build()

	return &discordgo.InteractionResponseData{
		Components: []discordgo.MessageComponent{comp},
		Flags:      discordgo.MessageFlagsIsComponentsV2 | discordgo.MessageFlagsEphemeral,
	}
}