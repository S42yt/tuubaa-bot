package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CookieHandler() func(*discordgo.Session, *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		opts := i.ApplicationCommandData().Options
		var user *discordgo.User
		if len(opts) > 0 {
			if u := opts[0].UserValue(s); u != nil {
				user = u
			}
		}

		if user == nil {
			return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: "Bitte gib einen Benutzer an."},
			})
		}

		author := i.Member.User

		messages := []string{
			"Ohaa %s, du hast einen Cookie 🍪 von %s bekommen :D. Wie toll!",
			"Einen Moment... %s, du hast einen Cookie 🍪 von %s bekommen?!?!?",
			"Eyyy %s! %s hat dir einen unglaublich leckeren Cookie 🍪 geschenkt!",
			"Ohh sieh mal %s, %s schenkt dir einen Cookie 🍪!",
			"Yooo %s!!! %s hat einen Cookie 🍪 aus der Dose geklaut für dich!!!",
			"%s wirft %s einen Cookie 🍪 an den Kopf. Treffer!",
		}

		text := messages[rand.Intn(len(messages))]
		content := fmt.Sprintf(text, user.Mention(), author.Mention())

		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: content},
		})
	}
}
