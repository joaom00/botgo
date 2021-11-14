package utils

import (
	"os"
	"strings"

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

func Find(arr []string, el string) bool {
	for _, v := range arr {
		if v == strings.ToUpper(el) {
			return true
		}
	}

	return false
}
