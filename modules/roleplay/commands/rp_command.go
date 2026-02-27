package commands

import (
	"math/rand"
	"strings"
	"time"

	"github.com/S42yt/tuubaa-bot/modules/roleplay/api"
	"github.com/S42yt/tuubaa-bot/modules/roleplay/embed"
	ulog "github.com/S42yt/tuubaa-bot/utils/logger"
	"github.com/bwmarrin/discordgo"
)

type rolePlayData struct {
	Text    string
	TextAll string
}

var rolePlayMap map[string]rolePlayData

func init() {
	rand.Seed(time.Now().UnixNano())
	rolePlayMap = map[string]rolePlayData{
		"blush":    {"user2, user1 errötet wegen dir!", "user1 errötet"},
		"cheers":   {"user2, prosit von user1!", "user1 prost!"},
		"cool":     {"user2, user1 findet dich cool!", "user1 ist cool"},
		"cry":      {"user2, user1 weint wegen dir!", "user1 weint"},
		"cuddle":   {"user2, user1 kuschelt mit dir!", "user1 kuschelt mit allen!"},
		"facepalm": {"user2, user1 facepalmt wegen dir!", "user1 facepalm"},
		"happy":    {"user2, user1 ist glücklich wegen dir!", "user1 ist glücklich"},
		"hug":      {"user2, user1 umarmt dich!", "user1 umarmt alle!"},
		"laugh":    {"user2, user1 lacht wegen dir!", "user1 lacht"},
		"love":     {"user2, user1 liebt dich! <:LoveTuba:1090372406897561791>", "user1 liebt alle! <:LoveTuba:1090372406897561791>"},
		"mad":      {"user2, user1 ist wütend auf dich!", "user1 ist wütend"},
		"nervous":  {"user2, user1 ist nervös wegen dir!", "user1 ist nervös"},
		"no":       {"user1 sagt nein zu user2!", "user1 NEIN!"},
		"pat":      {"user2, user1 streichelt dich! <:PatpatTuba:1120695389973123253>", "user1 streichelt alle! <:PatpatTuba:1120695389973123253>"},
		"sad":      {"user2, user1 ist traurig wegen dir! <:DepressedEMOTE:1226978508736172123>", "user1 ist traurig! <:DepressedEMOTE:1226978508736172123>"},
		"scared":   {"user2, user1 hat Angst vor dir! <:tuubaa_w:1235347591709982813>", "user1 hat Angst! <:tuubaa_w:1235347591709982813>"},
		"shy":      {"user2, user1 ist schüchtern wegen dir! <:tuubaa_verlegen_Emote:1236346476649513043>", "user1 ist schüchtern! <:tuubaa_verlegen_Emote:1236346476649513043>"},
		"slap":     {"user2, user1 schlägt dich!", "user1 schlägt"},
		"sleep":    {"user2, user1 schläft wegen dir! <:SleepTuba:1123745924720627722>", "user1 schläft! <:SleepTuba:1123745924720627722>"},
		"smug":     {"user2, user1 ist zufrieden wegen dir!", "user1 ist zufrieden"},
		"yay":      {"user2, user1 freut sich wegen dir!", "user1 freut sich"},
	}
}

func RolePlayHandler(kind string) func(*discordgo.Session, *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		ulog.Debug("RolePlayHandler invoked kind=%s user=%s", kind, i.Member.User.ID)
		data, ok := rolePlayMap[kind]
		if !ok {
			ulog.Warn("RolePlayHandler: unknown kind=%s", kind)
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Den Befehl gibt es irgendwie nicht.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return nil
		}

		var target *discordgo.User
		opts := i.ApplicationCommandData().Options
		if len(opts) > 0 && opts[0].Type == discordgo.ApplicationCommandOptionSubCommand {
			if len(opts[0].Options) > 0 {
				opts = opts[0].Options
			} else {
				opts = []*discordgo.ApplicationCommandInteractionDataOption{}
			}
		}
		if len(opts) > 0 {
			if u := opts[0].UserValue(s); u != nil {
				target = u
			}
		}

		gif, err := api.GetGifURL(kind)
		if err != nil || gif == "" {
			ulog.Error("RolePlayHandler: failed to fetch gif kind=%s err=%v", kind, err)
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Fehler beim Laden des GIFs.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return nil
		}

		user1 := i.Member.User

		var text string
		if target != nil {
			text = replaceUsers(data.Text, user1.Mention(), target.Mention())
		} else {
			text = replaceUsers(data.TextAll, user1.Mention(), "")
		}

		colors := []int{0x3498db, 0xe74c3c, 0x9b59b6, 0x2ecc71}
		accent := colors[rand.Intn(len(colors))]

		ulog.Debug("RolePlayHandler: building embed kind=%s gif=%s", kind, gif)

		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: embed.BuildResponse(text, gif, accent, i.Member.DisplayName()),
		}); err != nil {
			ulog.Error("RolePlayHandler: InteractionRespond failed kind=%s err=%v", kind, err)
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Fehler beim Senden der Antwort.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		return nil
	}
}

func replaceUsers(s, u1, u2 string) string {
	res := s
	if u1 != "" {
		res = replaceOnce(res, "user1", u1)
	}
	if u2 != "" {
		res = replaceOnce(res, "user2", u2)
	}
	return res
}

func replaceOnce(s, old, new string) string {
	return strings.Replace(s, old, new, 1)
}
