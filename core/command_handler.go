package core

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	ulog "github.com/S42yt/tuubaa-bot/utils/logger"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate) error

	Cooldown        int
	OnlyBotChannel  bool
	DisabledUsers   map[string]bool
	DisabledChannel map[string]bool
	AllowAdmin      bool
	AllowDev        bool
	AllowStaff      bool
	AllowEveryone   bool
}

var (
	mu                         sync.RWMutex
	commands                   = map[string]*Command{}
	attachMu                   sync.Mutex
	interactionHandlerAttached bool
	eventsAttached             bool
	messageHandlerAttached     bool
)

func Register(cmd *Command) error {
	if cmd == nil || cmd.Name == "" {
		return errors.New("invalid command")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := commands[cmd.Name]; ok {
		return fmt.Errorf("command %s already registered", cmd.Name)
	}
	commands[cmd.Name] = cmd
	return nil
}

func Init(s *discordgo.Session) error {
	if s.State == nil || s.State.User == nil {
		return errors.New("session state not ready; call Init from Ready handler")
	}
	for _, g := range s.State.Guilds {
		if err := InitWithGuild(s, g.ID); err != nil {
			ulog.Warn("failed to init commands for guild %s: %v", g.ID, err)
		}
	}
	return nil
}

func InitWithGuild(s *discordgo.Session, guildID string) error {
	appID := s.State.User.ID

	mu.RLock()
	defer mu.RUnlock()

	if globals, err := s.ApplicationCommands(appID, ""); err == nil {
		ulog.Debug("found %d global commands", len(globals))
	}
	if guilds, err := s.ApplicationCommands(appID, guildID); err == nil {
		ulog.Debug("found %d guild commands for %s", len(guilds), guildID)
	}

	if err := removeGlobalConflicts(s); err != nil {
		ulog.Warn("failed to remove global conflicting commands: %v", err)
	}

	if !modalHandlerAttached {
		s.AddHandler(modalInteractionHandler)
		modalHandlerAttached = true
	}

	if err := Clear(s, guildID); err != nil {
		ulog.Warn("failed to clear existing guild commands: %v", err)
	}

	if s.State == nil || s.State.User == nil {
		return errors.New("session state not ready; call InitWithGuild from Ready handler")
	}

	var appCommands []*discordgo.ApplicationCommand
	for _, c := range commands {
		appCommands = append(appCommands, &discordgo.ApplicationCommand{
			Name:        c.Name,
			Description: c.Description,
			Options:     c.Options,
		})
	}

	if _, err := s.ApplicationCommandBulkOverwrite(appID, guildID, appCommands); err != nil {
		return err
	}

	if cmds, err := s.ApplicationCommands(appID, guildID); err == nil {
		ulog.Info("published %d guild commands for guild %s", len(cmds), guildID)
	} else {
		ulog.Warn("could not fetch guild commands after publish: %v", err)
	}

	attachMu.Lock()
	defer attachMu.Unlock()

	if !interactionHandlerAttached {
		s.AddHandler(interactionHandler)
		interactionHandlerAttached = true
	}

	if !messageHandlerAttached {
		s.AddHandler(messageHandler)
		messageHandlerAttached = true
	}

	if !eventsAttached {
		for _, h := range eventHandlers {
			s.AddHandler(h)
		}
		eventsAttached = true
	}

	return nil
}

func Clear(s *discordgo.Session, guildID string) error {
	if s.State == nil || s.State.User == nil {
		return errors.New("session state not ready; call Clear after Ready")
	}
	appID := s.State.User.ID
	if _, err := s.ApplicationCommandBulkOverwrite(appID, guildID, []*discordgo.ApplicationCommand{}); err != nil {
		return err
	}
	return nil
}

func removeGlobalConflicts(s *discordgo.Session) error {
	if s.State == nil || s.State.User == nil {
		return errors.New("session state not ready; call removeGlobalConflicts after Ready")
	}
	appID := s.State.User.ID

	globals, err := s.ApplicationCommands(appID, "")
	if err != nil {
		return fmt.Errorf("fetch global commands: %w", err)
	}

	mu.RLock()
	defer mu.RUnlock()

	for _, g := range globals {
		if _, ok := commands[g.Name]; ok {
			if err := s.ApplicationCommandDelete(appID, "", g.ID); err != nil {
				ulog.Warn("failed to delete global command %s (%s): %v", g.Name, g.ID, err)
			} else {
				ulog.Info("deleted global command %s (%s)", g.Name, g.ID)
			}
		}
	}

	return nil
}

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	name := i.ApplicationCommandData().Name

	mu.RLock()
	cmd, ok := commands[name]
	mu.RUnlock()
	if !ok {
		_ = respondEphemeral(s, i, "This command no longer exists")
		return
	}

	if i.Member == nil {
		_ = respondEphemeral(s, i, "Could not validate member")
		return
	}

	member := i.Member.User

	if cmd.Cooldown > 0 {
		if remaining := checkCooldown(member.ID, cmd.Name); remaining > 0 {
			_ = respondEphemeral(s, i, fmt.Sprintf("You have a %d seconds cooldown", int(remaining)))
			return
		}
		setCooldown(member.ID, cmd.Name, cmd.Cooldown)
	}

	if cmd.OnlyBotChannel {
		if os.Getenv("BOT_CHANNEL") != "" && i.ChannelID != os.Getenv("BOT_CHANNEL") {
			_ = respondEphemeral(s, i, fmt.Sprintf("This command can only be executed in <#%s>", os.Getenv("BOT_CHANNEL")))
			return
		}
	}

	if cmd.DisabledUsers != nil {
		if cmd.DisabledUsers[member.ID] {
			_ = respondEphemeral(s, i, "You are excluded from this command")
			return
		}
	}
	if cmd.DisabledChannel != nil {
		if cmd.DisabledChannel[i.ChannelID] {
			_ = respondEphemeral(s, i, "This channel is excluded from this command")
			return
		}
	}

	isAdmin := false
	if i.Member.Permissions&discordgo.PermissionAdministrator != 0 {
		isAdmin = true
	}

	if !(cmd.AllowEveryone || cmd.AllowAdmin && isAdmin || cmd.AllowDev && os.Getenv("DEV_ID") == member.ID || cmd.AllowStaff && isAdmin) {
		_ = respondEphemeral(s, i, "You do not have permission for this command")
		return
	}

	if err := cmd.Handler(s, i); err != nil {
		_ = respondEphemeral(s, i, "An internal error occurred")
	}
}

func respondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	if i == nil || s == nil {
		return errors.New("nil session or interaction")
	}
	data := &discordgo.InteractionResponseData{Content: content, Flags: discordgo.MessageFlagsEphemeral}
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: data,
	})
}

var (
	cdMu    sync.Mutex
	cdStore = map[string]int64{}
)

func setCooldown(userID, cmdName string, seconds int) {
	cdMu.Lock()
	defer cdMu.Unlock()
	key := userID + ":" + cmdName
	cdStore[key] = time.Now().Unix() + int64(seconds)
}

func checkCooldown(userID, cmdName string) int64 {
	cdMu.Lock()
	defer cdMu.Unlock()
	key := userID + ":" + cmdName
	exp, ok := cdStore[key]
	if !ok {
		return 0
	}
	now := time.Now().Unix()
	if exp <= now {
		delete(cdStore, key)
		return 0
	}
	return exp - now
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) == 0 || m.Content[0] != '!' {
		return
	}

	command := m.Content[1:]

	switch command {
	case "registerCommands":
		handleRegisterCommands(s, m)
	}
}

func handleRegisterCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildID := m.GuildID
	if guildID == "" {
		ch, err := s.Channel(m.ChannelID)
		if err != nil || ch.GuildID == "" {
			ulog.Warn("handleRegisterCommands: could not resolve guild ID: %v", err)
			s.ChannelMessageSend(m.ChannelID, "Could not resolve Guild ID")
			return
		}
		guildID = ch.GuildID
	}

	ulog.Debug("handleRegisterCommands: resolved GUILD_ID=%s", guildID)
	s.ChannelMessageSend(m.ChannelID, "Registering commands...")

	if err := InitWithGuild(s, guildID); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to register commands: %v", err))
		ulog.Error("Failed to register commands: %v", err)
		return
	}

	mu.RLock()
	count := len(commands)
	mu.RUnlock()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully registered %d commands!", count))
	ulog.Info("Commands registered via !registerCommands by user %s", m.Author.ID)
}