package rules

import (
	_ "github.com/S42yt/tuubaa-bot/modules/misc/events"
	"github.com/S42yt/tuubaa-bot/core"
	"github.com/S42yt/tuubaa-bot/modules/misc/commands"
	"github.com/bwmarrin/discordgo"
)

func init() {
	ruleCmd := &core.Command{
		Name:        "rule",
		Description: "Schicke eine bestimmte Regel in den Chat (Es ist Anonym!)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "rule",
				Description: "Wähle eine Regel",
				Required:    true,
				Choices:     commands.GetRuleChoices(),
			},
		},
		AllowEveryone: true,
		Handler:       commands.RuleHandler(),
	}

	setupCmd := &core.Command{
		Name:        "setup",
		Description: "Server-Setup Befehle",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "rules",
				Description: "Erstelle das Regelwerk",
			},
		},
		AllowAdmin: true,
		Handler:    commands.SetRuleHandler(),
	}

	modalHandler := &core.ModalHandler{
		CustomID:   commands.SetRuleModalID,
		AllowAdmin: true,
		Handler:    commands.SetRuleModalHandler(),
	}

	_ = core.Register(ruleCmd)
	_ = core.Register(setupCmd)
	_ = core.RegisterModal(modalHandler)
}