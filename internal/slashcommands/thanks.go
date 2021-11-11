package slashcommands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/internal/embedbuilder"
	"github.com/joaom00/botgo/internal/utils"
)

var ThanksCmd = &discordgo.ApplicationCommand{
	Name:        "thanks",
	Description: "Agradeça um usuário por ter te ajudado(a)",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "usuário",
			Description: "Usuário o qual você deseja agradecer",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "motivo",
			Description: "De que forma o usuário lhe ajudou",
			Required:    true,
		},
	},
}

func ThanksHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	authorName := fmt.Sprintf("%s#%s", i.ApplicationCommandData().Options[0].UserValue(s).Username, i.ApplicationCommandData().Options[0].UserValue(s).Discriminator)
	authorURL := i.ApplicationCommandData().Options[0].UserValue(s).AvatarURL("")
	title := fmt.Sprintf("%s foi agradecido(a)! :tada:", i.ApplicationCommandData().Options[0].UserValue(s).Username)
	reason := fmt.Sprintf("%s", i.ApplicationCommandData().Options[1].StringValue())
	link := fmt.Sprintf("https://discordapp.com/channels/%s/%s/%s", i.GuildID, i.ChannelID, i.Interaction.ID)

	embed := embedbuilder.New().
		SetAuthor(authorName, authorURL).
		SetThumbnail(i.ApplicationCommandData().Options[0].UserValue(s).AvatarURL("")).
		SetTitle(title).
		AddField("Motivo", reason).
		AddField("Link", link).
		SetColor(0xF83B91).
		MessageEmbed

	channelID, err := utils.GetChannelByName(s, "thanks")
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "Algo deu errado :(\nPor favor, reporte a um moderador",
			},
		})
		return
	}

	if _, err = s.ChannelMessageSendEmbed(channelID, embed); err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "Algo deu errado :(\nPor favor, reporte a um moderador",
			},
		})
		return
	}

	args := []interface{}{
		i.ApplicationCommandData().Options[0].UserValue(s).Mention(),
		i.Member.Mention(),
		i.ApplicationCommandData().Options[1].StringValue(),
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Ei %s!, você foi agradecido(a) :tada:\n\n%s agradeceu você por:\n\n> %s", args...),
		},
	})
	if err != nil {
		log.Println(err)
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Flags:   1 << 6,
			Content: "Algo deu errado :(\nPor favor, reporte a um moderador",
		})
		return
	}
}
