package commands

import (
	"fmt"
	"regexp"
	"strings"

	vembed "github.com/S42yt/tuubaa-bot/modules/misc/embed"
	"github.com/bwmarrin/discordgo"
	ulog "github.com/S42yt/tuubaa-bot/utils/logger"
)

var ruleChoices = []*discordgo.ApplicationCommandOptionChoice{
	{Name: "§1 Begegne allen Nutzern jederzeit freundlich und respektvoll.", Value: "1"},
	{Name: "§2 Beachte die Nutzungsbedingungen (Terms of Service) von Discord.", Value: "2"},
	{Name: "§3 Werbung für eigene oder fremde Inhalte ist nicht erlaubt.", Value: "3"},
	{Name: "§4 Diskussionen über sensible Themen wie Politik oder Religion sind untersagt.", Value: "4"},
	{Name: "§5 Beleidigungen, Provokationen sowie rassistische, sexistische oder radikale...", Value: "5"},
	{Name: "§6 Namen, Profilbilder und Status dürfen keine Beleidigungen, Provokationen...", Value: "6"},
	{Name: "§7 Das Vortäuschen einer fremden Identität ist verboten.", Value: "7"},
	{Name: "§8 Die Nutzung mehrerer Discord-Accounts ist untersagt.", Value: "8"},
	{Name: "§9 Störungen in Sprachkanälen durch laute Geräusche, Stimmverzerrer...", Value: "9"},
	{Name: "§10 Das Teilen von NSFW-Inhalten oder ähnlichem ist strengstens untersagt.", Value: "10"},
	{Name: "§11 Der Support darf nicht missbraucht werden. Wendet euch nur bei...", Value: "11"},
	{Name: "§12 Betteln oder Nachfragen nach Rängen ist nicht gestattet.", Value: "12"},
	{Name: "§13 Den Anweisungen des Teams ist Folge zu leisten. In Zweifelsfällen hat...", Value: "13"},
	{Name: "§14 Kein \"Backseat Arting\": Wenn jemand eine Zeichnung postet und nicht...", Value: "14"},
	{Name: "§15 Dating, Flirten oder unangemessenes Verhalten sind auf dem Server nicht...", Value: "15"},
	{Name: "§16 Bitte fragt tuubaa nicht, ob ich eure Freundschaftsanfragen...", Value: "16"},
}

// ruleTexts is the canonical list used both for /rule and as the fallback for /setup rules
var ruleTexts = map[string]string{
	"1":  "**§1** Begegne allen Nutzern jederzeit freundlich und respektvoll.",
	"2":  "**§2** Beachte die Nutzungsbedingungen (Terms of Service) von Discord.",
	"3":  "**§3** Werbung für eigene oder fremde Inhalte ist nicht erlaubt.",
	"4":  "**§4** Diskussionen über sensible Themen wie Politik oder Religion sind untersagt.",
	"5":  "**§5** Beleidigungen, Provokationen sowie rassistische, sexistische oder radikale Aussagen werden nicht toleriert.",
	"6":  "**§6** Namen, Profilbilder und Status dürfen keine Beleidigungen, Provokationen oder extremen Aussagen enthalten. Bei einem Hinweis durch das Team sind diese unverzüglich zu ändern.",
	"7":  "**§7** Das Vortäuschen einer fremden Identität ist verboten.",
	"8":  "**§8** Die Nutzung mehrerer Discord-Accounts ist untersagt.",
	"9":  "**§9** Störungen in Sprachkanälen durch laute Geräusche, Stimmverzerrer, Soundboards o. Ä. sind verboten.",
	"10": "**§10** Das Teilen von NSFW-Inhalten oder ähnlichem ist strengstens untersagt.",
	"11": "**§11** Der Support darf nicht missbraucht werden. Wendet euch nur bei ernsthaften Anliegen an das Team.",
	"12": "**§12** Betteln oder Nachfragen nach Rängen ist nicht gestattet.",
	"13": "**§13** Den Anweisungen des Teams ist Folge zu leisten. In Zweifelsfällen hat das Team Entscheidungsrecht, auch über das Regelwerk hinaus.",
	"14": "**§14** Kein \"Backseat Arting\": Wenn jemand eine Zeichnung postet und nicht explizit um Feedback bittet, ist jegliche Form von Kritik zu unterlassen (z. B. „Ich hätte das anders gemacht.).",
	"15": "**§15** Dating, Flirten oder unangemessenes Verhalten sind auf dem Server nicht gestattet.",
	"16": "**§16** Bitte fragt **tuubaa** nicht, ob ich eure Freundschaftsanfragen annehmen euch in einem Video malen kann.",
}

func defaultRuleText() string {
	lines := make([]string, 0, len(ruleTexts))
	for i := 1; i <= len(ruleTexts); i++ {
		if text, ok := ruleTexts[fmt.Sprintf("%d", i)]; ok {
			lines = append(lines, text)
		}
	}
	return strings.Join(lines, "\n\n")
}

const SetRuleModalID = "set_rule_modal"

func RuleHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			return nil
		}

		ruleValue := options[0].StringValue()
		ruleText, ok := ruleTexts[ruleValue]
		if !ok {
			return fmt.Errorf("unknown rule value: %s", ruleValue)
		}

		_, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Components: vembed.BuildRuleEmbed(ruleText),
			Flags:      discordgo.MessageFlagsIsComponentsV2,
		})
		if err != nil {
			return fmt.Errorf("failed to send rule message: %w", err)
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: vembed.BuildRuleSuccessResponse(),
		})
		if err != nil {
			return fmt.Errorf("failed to respond to interaction: %w", err)
		}
		return nil
	}
}

func SetRuleHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: SetRuleModalID,
				Title:    "Regel Admin Interface",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "rule",
								Label:       "Setzt die Regel (leer = Standard)",
								Style:       discordgo.TextInputParagraph,
								Placeholder: "Leer lassen um die Standard-Regeln zu verwenden.",
								Required:    false,
							},
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to open set_rule modal: %w", err)
		}
		return nil
	}
}

func SetRuleModalHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		data := i.ModalSubmitData()

		ruleContent := ""
		for _, row := range data.Components {
			if ar, ok := row.(*discordgo.ActionsRow); ok {
				for _, comp := range ar.Components {
					if ti, ok := comp.(*discordgo.TextInput); ok && ti.CustomID == "rule" {
						ruleContent = strings.TrimSpace(ti.Value)
					}
				}
			}
		}

		if ruleContent == "" {
			ruleContent = defaultRuleText()
		} else {
			re := regexp.MustCompile(`(§\d+)`)
			ruleContent = re.ReplaceAllString(ruleContent, "**$1**")
		}

		msg, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Components: vembed.BuildSetRuleEmbed(ruleContent, s, i.Interaction),
			Flags:      discordgo.MessageFlagsIsComponentsV2,
		})
		if err != nil {
			return fmt.Errorf("failed to send rulebook: %w", err)
		}

		ulog.Debug("Rulebook set in channel %s, message %s by %s", i.ChannelID, msg.ID, i.Member.User.Username)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: vembed.BuildRuleSuccessResponse(),
		})
		if err != nil {
			return fmt.Errorf("failed to respond after modal submit: %w", err)
		}

		return nil
	}
}

func GetRuleChoices() []*discordgo.ApplicationCommandOptionChoice {
	return ruleChoices
}