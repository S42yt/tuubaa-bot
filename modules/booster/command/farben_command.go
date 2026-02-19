package command

import (
	"fmt"
	"os"

	bembed "github.com/S42yt/tuubaa-bot/modules/booster/embed"
	ulog "github.com/S42yt/tuubaa-bot/utils/logger"
	"github.com/bwmarrin/discordgo"
)

func FarbenHandler() func(*discordgo.Session, *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		ulog.Debug("FarbenHandler invoked user=%s", i.Member.User.ID)

		data := i.ApplicationCommandData()
		if len(data.Options) == 0 {
			resp := &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Content: "Bitte eine Auswahl treffen."}}
			if err := s.InteractionRespond(i.Interaction, resp); err != nil {
				ulog.Warn("FarbenHandler: InteractionRespond failed: %v", err)
				return err
			}
			return nil
		}

		choice := data.Options[0].StringValue()

		selectable := []string{"Unschuldiges Kind", "Verdächtiges Kind", "Schuldiges Kind", "Mit Entführer", "Meisterentführer", "Beifahrer", "Van Upgrader"}

		envMap := map[string]string{
			"Unschuldiges Kind": os.Getenv("ROLE_UNSCHULDIGES_KIND"),
			"Verdächtiges Kind": os.Getenv("ROLE_VERDAECHTIGES_KIND"),
			"Schuldiges Kind":   os.Getenv("ROLE_SCHULDIGES_KIND"),
			"Mit Entführer":     os.Getenv("ROLE_MIT_ENTFUEHRER"),
			"Meisterentführer":  os.Getenv("ROLE_MEISTERENTFUEHRER"),
			"Beifahrer":         os.Getenv("ROLE_BEIFAHRER"),
			"Van Upgrader":      os.Getenv("ROLE_VAN_UPGRADER"),
		}

		boosterRoleID := os.Getenv("ROLE_VAN_UPGRADER")
		if boosterRoleID == "" {
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
			data := bembed.BuildResponse("Ungültige Auswahl", "Die gewählte Option ist ungültig. wie hast du das gemacht o.o", 0xe74c3c, "", true)
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: data})
		}

		removeIDs := []string{}
		for _, n := range selectable {
			id := envMap[n]
			if id == "" {
				continue
			}
			// Never include the booster/Van Upgrader role in the generic removal list
			// so selecting a normal color does not strip the Van Upgrader role.
			if id == boosterRoleID {
				continue
			}
			removeIDs = append(removeIDs, id)
		}

		selRoleID := envMap[choice]
		if selRoleID == "" {
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
		return nil
	}
}
