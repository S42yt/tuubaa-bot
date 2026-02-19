package embed

import (
	"fmt"

	v2 "github.com/S42yt/tuubaa-bot/utils/embed"
	"github.com/bwmarrin/discordgo"
)

func buildContainerMessage(accent int, lines ...string) *discordgo.InteractionResponseData {
	c := v2.NewContainerBuilder().SetAccentColor(accent)
	for _, l := range lines {
		c.AddComponent(v2.NewTextDisplayBuilder().SetContent(l).Build())
	}
	comp := c.Build()
	return &discordgo.InteractionResponseData{Components: []discordgo.MessageComponent{comp}, Flags: discordgo.MessageFlagsIsComponentsV2}
}

func MissingBoosterConfig() *discordgo.InteractionResponseData {
	return buildContainerMessage(0xe74c3c, "## Van Upgrader fehlt", "Die Van Upgrader ist auf diesem Server nicht konfiguriert.")
}

func AccessDenied() *discordgo.InteractionResponseData {
	return buildContainerMessage(0xe74c3c, "## Zugriff verweigert :(((", "Du benötigst die Van Upgrader Rolle, um diesen Befehl zu verwenden :(.")
}

func InvalidSelection() *discordgo.InteractionResponseData {
	return buildContainerMessage(0xe74c3c, "## Ungültige Auswahl", "Die gewählte Option ist ungültig lol.")
}

func RoleNotConfigured(choice string) *discordgo.InteractionResponseData {
	return buildContainerMessage(0xe74c3c, "## Rolle nicht konfiguriert", fmt.Sprintf("Die Rolle '%s' ist nicht konfiguriert EHRE", choice))
}

func Error(msg string) *discordgo.InteractionResponseData {
	return buildContainerMessage(0xe74c3c, "## Fehler", msg)
}

func Success(choice, thumb string) *discordgo.InteractionResponseData {
	c := v2.NewContainerBuilder().SetAccentColor(0x2ecc71)
	c.AddComponent(v2.NewTextDisplayBuilder().SetContent(fmt.Sprintf("## Rolle gesetzt: %s", choice)).Build())
	c.AddComponent(v2.NewTextDisplayBuilder().SetContent(fmt.Sprintf("Rolle '%s' gesetzt.", choice)).Build())
	comp := c.Build()
	if thumb != "" {
		sb := v2.NewSectionBuilder()
		for _, inner := range c.Components {
			sb.AddComponent(inner)
		}
		sb.SetAccessory(v2.NewThumbnailBuilder().SetURL(thumb).Build())
		sec := sb.Build()

		outer := v2.NewContainerBuilder().SetAccentColor(0x2ecc71)
		outer.AddComponent(sec)
		outerComp := outer.Build()
		return &discordgo.InteractionResponseData{Components: []discordgo.MessageComponent{outerComp}, Flags: discordgo.MessageFlagsIsComponentsV2}
	}

	return &discordgo.InteractionResponseData{Components: []discordgo.MessageComponent{comp}, Flags: discordgo.MessageFlagsIsComponentsV2}
}

func VanUpgraderSuccess(thumb string) *discordgo.InteractionResponseData {
	c := v2.NewContainerBuilder().SetAccentColor(0x2ecc71)
	c.AddComponent(v2.NewTextDisplayBuilder().SetContent("## Van Upgrader gesetzt :)").Build())
	c.AddComponent(v2.NewTextDisplayBuilder().SetContent("Van Upgrader gesetzt. JETZT BIST DU STANDART HEHE").Build())
	comp := c.Build()
	if thumb != "" {
		sb := v2.NewSectionBuilder()
		for _, inner := range c.Components {
			sb.AddComponent(inner)
		}
		sb.SetAccessory(v2.NewThumbnailBuilder().SetURL(thumb).Build())
		sec := sb.Build()

		outer := v2.NewContainerBuilder().SetAccentColor(0x2ecc71)
		outer.AddComponent(sec)
		outerComp := outer.Build()
		return &discordgo.InteractionResponseData{Components: []discordgo.MessageComponent{outerComp}, Flags: discordgo.MessageFlagsIsComponentsV2}
	}

	return &discordgo.InteractionResponseData{Components: []discordgo.MessageComponent{comp}, Flags: discordgo.MessageFlagsIsComponentsV2}
}

func BuildResponse(title, body string, accent int, thumb string, ephemeral bool) *discordgo.InteractionResponseData {
	c := v2.NewContainerBuilder().SetAccentColor(accent)
	c.AddComponent(v2.NewTextDisplayBuilder().SetContent(fmt.Sprintf("## %s", title)).Build())
	c.AddComponent(v2.NewTextDisplayBuilder().SetContent(body).Build())
	comp := c.Build()

	flags := discordgo.MessageFlagsIsComponentsV2
	if ephemeral {
		flags |= discordgo.MessageFlagsEphemeral
	}

	if thumb != "" {
		sb := v2.NewSectionBuilder()
		for _, inner := range c.Components {
			sb.AddComponent(inner)
		}
		sb.SetAccessory(v2.NewThumbnailBuilder().SetURL(thumb).Build())
		sec := sb.Build()

		outer := v2.NewContainerBuilder().SetAccentColor(accent)
		outer.AddComponent(sec)
		outerComp := outer.Build()
		return &discordgo.InteractionResponseData{Components: []discordgo.MessageComponent{outerComp}, Flags: flags}
	}

	return &discordgo.InteractionResponseData{Components: []discordgo.MessageComponent{comp}, Flags: flags}
}
