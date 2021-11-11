package utils

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func GetChannelByName(s *discordgo.Session, channelName string) (channelID string, err error) {
	st, err := s.GuildChannels(os.Getenv("GUILD_ID"))
	if err != nil {
		return channelID, err
	}

	for _, c := range st {
		if c.Name == channelName {
			channelID = c.ID
			return
		}
	}

	return channelID, nil

}
