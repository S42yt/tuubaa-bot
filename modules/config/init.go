package config

import (
	"github.com/S42yt/tuubaa-bot/core"
	"github.com/S42yt/tuubaa-bot/modules/config/commands"
	"github.com/bwmarrin/discordgo"
)

func init() {
	setRole := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "setrole",
		Description: "Set a configured role for this guild",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "role",
				Description: "Which configurable role to set",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Unschuldiges Kind", Value: "ROLE_UNSCHULDIGES_KIND"},
					{Name: "Verdächtiges Kind", Value: "ROLE_VERDAECHTIGES_KIND"},
					{Name: "Schuldiges Kind", Value: "ROLE_SCHULDIGES_KIND"},
					{Name: "Mit Entführer", Value: "ROLE_MIT_ENTFUEHRER"},
					{Name: "Meisterentführer", Value: "ROLE_MEISTERENTFUEHRER"},
					{Name: "Beifahrer", Value: "ROLE_BEIFAHRER"},
					{Name: "Van Upgrader", Value: "ROLE_VAN_UPGRADER"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "target",
				Description: "The Discord role to assign for this key",
				Required:    true,
			},
		},
	}

	cfgCmd := &core.Command{
		Name:        "config",
		Description: "Guild-specific configuration",
		Options:     []*discordgo.ApplicationCommandOption{setRole},
		AllowAdmin:  true,
		Handler:     commands.ConfigRoleHandler(),
	}

	_ = core.Register(cfgCmd)
}
