package command

import (
	"fmt"

	bembed "github.com/S42yt/tuubaa-bot/modules/booster/embed"
	"github.com/S42yt/tuubaa-bot/modules/config"
	ulog "github.com/S42yt/tuubaa-bot/utils/logger"
	"github.com/bwmarrin/discordgo"
)

func FarbenHandler() func(*discordgo.Session, *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		ulog.Debug("FarbenHandler invoked: guild=%s", i.GuildID)
		ulog.Debug("FarbenHandler invoked user=%s", i.Member.User.ID)

		data := i.ApplicationCommandData()
		ulog.Debug("FarbenHandler: interaction data name=%s options=%d", data.Name, len(data.Options))
		if len(data.Options) == 0 {
			ulog.Warn("FarbenHandler: no options provided")
			resp := &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Content: "Bitte eine Auswahl treffen."}}
			if err := s.InteractionRespond(i.Interaction, resp); err != nil {
				ulog.Warn("FarbenHandler: InteractionRespond failed: %v", err)
				return err
			}
			return nil
		}

		choice := data.Options[0].StringValue()
		ulog.Debug("FarbenHandler: user %s selected=%s", i.Member.User.ID, choice)

		
		selectable := []string{"Unschuldiges Kind", "Verdächtiges Kind", "Schuldiges Kind", "Mit Entführer", "Meisterentführer", "Beifahrer", "Van Upgrader"}
		choiceKey := map[string]string{
			"Unschuldiges Kind": "ROLE_UNSCHULDIGES_KIND",
			"Verdächtiges Kind": "ROLE_VERDAECHTIGES_KIND",
			"Schuldiges Kind":   "ROLE_SCHULDIGES_KIND",
			"Mit Entführer":     "ROLE_MIT_ENTFUEHRER",
			"Meisterentführer":  "ROLE_MEISTERENTFUEHRER",
			"Beifahrer":         "ROLE_BEIFAHRER",
			"Van Upgrader":      "ROLE_VAN_UPGRADER",
		}

		rolesMap, err := config.GetRoles(i.GuildID)
		if err != nil {
			ulog.Warn("FarbenHandler: GetRoles error for guild %s: %v", i.GuildID, err)
			data := bembed.BuildResponse("Fehler", "Fehler beim Lesen der Konfiguration.", 0xe74c3c, "", true)
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: data})
		}

		envMap := map[string]string{}
		for human, key := range choiceKey {
			if v, ok := rolesMap[key]; ok {
				envMap[human] = v
			} else {
				envMap[human] = ""
			}
		}
		ulog.Debug("FarbenHandler: envMap(human->roleID)=%v rolesMap=%v", envMap, rolesMap)

		boosterRoleID := envMap["Van Upgrader"]
		ulog.Debug("FarbenHandler: boosterRoleID=%s (from envMap)", boosterRoleID)
		if boosterRoleID == "" {
			ulog.Warn("FarbenHandler: missing booster role config for guild=%s", i.GuildID)
			data := bembed.BuildResponse("Van Upgrader fehlt", "Die Van Upgrader-Rolle ist auf diesem Server nicht konfiguriert. lol", 0xe74c3c, "", true)
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: data})
		}

		hasBooster := false
		for _, rid := range i.Member.Roles {
			if rid == boosterRoleID {
				hasBooster = true
				break
			}
		}
		if !hasBooster {
			ulog.Debug("FarbenHandler: user %s does not have booster role", i.Member.User.ID)
			data := bembed.BuildResponse("Nicht special :((((", "Du benötigst die Van Upgrader Rolle, um diesen Befehl zu verwenden :(", 0xe74c3c, "", true)
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: data})
		}

		found := false
		for _, n := range selectable {
			if n == choice {
				found = true
				break
			}
		}
		if !found {
			ulog.Warn("FarbenHandler: invalid selection '%s' from user %s", choice, i.Member.User.ID)
			data := bembed.BuildResponse("Ungültige Auswahl", "Die gewählte Option ist ungültig. wie hast du das gemacht o.o", 0xe74c3c, "", true)
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: data})
		}

		removeIDs := []string{}
		for _, n := range selectable {
			id := envMap[n]
			if id == "" {
				continue
			}
			if id == boosterRoleID {
				continue
			}
			removeIDs = append(removeIDs, id)
		}
		ulog.Debug("FarbenHandler: removeIDs=%v", removeIDs)

		selRoleID := envMap[choice]
		if selRoleID == "" {
			ulog.Warn("FarbenHandler: selected role not configured for choice=%s guild=%s", choice, i.GuildID)
			data := bembed.BuildResponse("Rolle nicht konfiguriert", fmt.Sprintf("Die Rolle '%s' ist nicht konfiguriert LOL meld dich an musa", choice), 0xe74c3c, "", true)
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: data})
		}

		if choice == "Van Upgrader" {
			for _, rid := range removeIDs {
				if rid == selRoleID {
					continue
				}
				for _, ur := range i.Member.Roles {
					if ur == rid {
						if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, rid); err != nil {
							ulog.Warn("FarbenHandler: failed to remove role %s from %s: %v", rid, i.Member.User.ID, err)
						}
					}
				}
			}
			if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, selRoleID); err != nil {
				ulog.Error("FarbenHandler: failed to add role %s to %s: %v", selRoleID, i.Member.User.ID, err)
				data := bembed.BuildResponse("Fehler", "Fehler beim Hinzufügen der Rolle.", 0xe74c3c, "", true)
				if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: data}); err != nil {
					ulog.Warn("FarbenHandler: InteractionRespond failed: %v", err)
					return err
				}
				ulog.Debug("FarbenHandler: aborting after failed add")
				return nil
			}

			thumb := ""
			if i.Member.User.Avatar != "" {
				thumb = i.Member.User.AvatarURL("1024")
			}
			resp := bembed.BuildResponse(fmt.Sprintf("Farbe gesetzt: %s", choice), "Van Upgrader gesetzt. Jetzt bist du ein normaler Upgrader hehe", 0x2ecc71, thumb, true)
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: resp}); err != nil {
				ulog.Warn("FarbenHandler: InteractionRespond failed: %v", err)
				return err
			}
			ulog.Debug("FarbenHandler: Van Upgrader set for user %s", i.Member.User.ID)
			return nil
		}

		for _, rid := range removeIDs {
			if rid == selRoleID {
				continue
			}
			for _, ur := range i.Member.Roles {
				if ur == rid {
					if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, rid); err != nil {
						ulog.Warn("FarbenHandler: failed to remove role %s from %s: %v", rid, i.Member.User.ID, err)
					}
				}
			}
		}

		if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, selRoleID); err != nil {
			ulog.Error("FarbenHandler: failed to add role %s to %s: %v", selRoleID, i.Member.User.ID, err)
			data := bembed.BuildResponse("Fehler", "Fehler beim Hinzufügen der Rolle.", 0xe74c3c, "", true)
			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: data}); err != nil {
				ulog.Warn("FarbenHandler: InteractionRespond failed: %v", err)
				return err
			}
			ulog.Debug("FarbenHandler: aborting after failed add (non-van)")
			return nil
		}

		thumb := ""
		if i.Member.User.Avatar != "" {
			thumb = i.Member.User.AvatarURL("1024")
		}
		resp := bembed.BuildResponse(fmt.Sprintf("Farbe bekommen: %s", choice), fmt.Sprintf("Du hast die Farbe '%s' erhalten. YAY", choice), 0x2ecc71, thumb, true)
		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: resp}); err != nil {
			ulog.Warn("FarbenHandler: InteractionRespond failed: %v", err)
			return err
		}
		ulog.Debug("FarbenHandler: role %s assigned to user %s", selRoleID, i.Member.User.ID)
		return nil
	}
}
