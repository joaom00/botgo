package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/internal/utils"
)

func Thanks(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")

	if len(args) == 1 {
		s.ChannelMessageSend(m.ChannelID, "Comando: !thanks <@user> <motivo>")
		return
	}

	reason := strings.Join(args[2:], " ")
	msg := fmt.Sprintf(`Ei %s!, você foi agradecido :tada:

%s agradeceu à você por:

> %s

Link: https://discordapp.com/channels/%s/%s/%s
        `,
		m.Mentions[0].Mention(),
		m.Author.Mention(),
		reason,
		m.GuildID,
		m.ChannelID,
		m.Message.ID,
	)

	channelID, err := utils.GetChannelByName(s, "thanks")
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Algo deu errado :(\nPor favor, reporte a um moderador")
		return
	}

	_, err = s.ChannelMessageSend(channelID, msg)
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Algo deu errado :(\nPor favor, reporte a um moderador")
		return
	}
}
