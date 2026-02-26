package core

import (
	"errors"
	"fmt"
	"sync"

	ulog "github.com/S42yt/tuubaa-bot/utils/logger"
	"github.com/bwmarrin/discordgo"
)

type ModalHandler struct {
	CustomID    string
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate) error
	AllowAdmin  bool
	AllowDev    bool
	AllowStaff  bool
	AllowEveryone bool
}

var (
	modalMu              sync.RWMutex
	modalHandlers        = map[string]*ModalHandler{}
	modalHandlerAttached bool
)

func RegisterModal(m *ModalHandler) error {
	if m == nil || m.CustomID == "" {
		return errors.New("invalid modal handler: missing CustomID")
	}
	modalMu.Lock()
	defer modalMu.Unlock()
	if _, ok := modalHandlers[m.CustomID]; ok {
		return fmt.Errorf("modal handler %s already registered", m.CustomID)
	}
	modalHandlers[m.CustomID] = m
	ulog.Debug("registered modal handler: %s", m.CustomID)
	return nil
}

func attachModalHandler(s *discordgo.Session) {
	attachMu.Lock()
	defer attachMu.Unlock()
	if !modalHandlerAttached {
		s.AddHandler(modalInteractionHandler)
		modalHandlerAttached = true
		ulog.Debug("modal interaction handler attached")
	}
}

func modalInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	customID := i.ModalSubmitData().CustomID

	modalMu.RLock()
	m, ok := modalHandlers[customID]
	modalMu.RUnlock()

	if !ok {
		ulog.Warn("no modal handler registered for customID: %s", customID)
		_ = respondEphemeral(s, i, "This modal is no longer handled")
		return
	}

	if i.Member == nil {
		_ = respondEphemeral(s, i, "Could not validate member")
		return
	}

	isAdmin := i.Member.Permissions&discordgo.PermissionAdministrator != 0

	if !(m.AllowEveryone || m.AllowAdmin && isAdmin || m.AllowStaff && isAdmin) {
		_ = respondEphemeral(s, i, "You do not have permission to submit this modal")
		return
	}

	if err := m.Handler(s, i); err != nil {
		ulog.Error("modal handler %s error: %v", customID, err)
		_ = respondEphemeral(s, i, "An internal error occurred")
	}
}