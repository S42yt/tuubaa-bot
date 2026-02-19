package embed

import (
	"fmt"

	v2 "github.com/S42yt/tuubaa-bot/utils/embed"
	"github.com/bwmarrin/discordgo"
)

func BuildComponents(text, imageURL string, accentColor int, executorName string) discordgo.MessageComponent {
	mainText := v2.NewTextDisplayBuilder().
		SetContent(fmt.Sprintf("## %s", text)).
		Build()

	mg := v2.NewMediaGalleryBuilder().
		AddImageURL(imageURL).
		Build()

	footer := v2.NewTextDisplayBuilder().
		SetContent(fmt.Sprintf("-# von %s", executorName)).
		Build()

	return v2.NewContainerBuilder().
		SetAccentColor(accentColor).
		AddComponent(mainText).
		AddComponent(mg).
		AddComponent(footer).
		Build()
}

func BuildResponse(text, imageURL string, accentColor int, executorName string) *discordgo.InteractionResponseData {
	comp := BuildComponents(text, imageURL, accentColor, executorName)

	return &discordgo.InteractionResponseData{
		Components: []discordgo.MessageComponent{comp},
		Flags:      discordgo.MessageFlagsIsComponentsV2,
	}
}