package slashcommands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/helper"
	"github.com/joaom00/botgo/internal/services/database"
)

var InvestCmd = &discordgo.ApplicationCommand{
	Name:        "investir",
	Description: "Invista suas JCoins em criptomoedas",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "moeda",
			Description: "Criptomoeda que você deseja investir",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "quantidade",
			Description: "Quantidade de JCoins que você deseja investir",
			Required:    true,
		},
	},
}

func InvestHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cryptoCoin := i.ApplicationCommandData().Options[0].StringValue()
	quantity := i.ApplicationCommandData().Options[1].FloatValue()

	if !helper.Find(helper.Coins, cryptoCoin) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("A moeda `%s` não existe. Use `/moedas` para conhecer as moedas válidas", cryptoCoin),
			},
		})
		return
	}

	wallet, err := database.GetWallet(i.Member.User.ID)
	if err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "voce nao tem carteira",
			},
		})
		return
	}

	if wallet.Amount < quantity {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Você não possui essa quantidade para investir",
			},
		})
		return
	}

	if err = database.Invest(wallet, cryptoCoin, quantity); err != nil {
		log.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "deu erro",
			},
		})
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "A moeda existe",
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
