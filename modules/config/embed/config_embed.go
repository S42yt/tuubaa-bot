package embed

import (
	"fmt"

	v2 "github.com/S42yt/tuubaa-bot/utils/embed"
	"github.com/bwmarrin/discordgo"
)

func BuildRoleSetResponse(roleKey, roleID, executor string) *discordgo.InteractionResponseData {
	main := v2.NewTextDisplayBuilder().SetContent(fmt.Sprintf("## Configuration updated: %s", roleKey)).Build()
	footer := v2.NewTextDisplayBuilder().SetContent(fmt.Sprintf("- Set to <@&%s> by %s", roleID, executor)).Build()

	comp := v2.NewContainerBuilder().SetAccentColor(0x99EE99).
		AddComponent(main).
		AddComponent(footer).
		Build()

	return &discordgo.InteractionResponseData{
		Components: []discordgo.MessageComponent{comp},
		Flags:      discordgo.MessageFlagsIsComponentsV2 | discordgo.MessageFlagsEphemeral,
	}
}
