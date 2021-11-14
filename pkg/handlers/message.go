package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/pkg/commands"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	switch strings.Split(m.Content, " ")[0] {
	case "!obrigado", "!obrigada":
		commands.Thanks(s, m)
	}

}
