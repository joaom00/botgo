package slashcommands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/internal/database"
	"github.com/joaom00/botgo/utils"
	"go.mongodb.org/mongo-driver/mongo"
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

	if !utils.Find(utils.Coins, cryptoCoin) {
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
		if err == mongo.ErrNoDocuments {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Você não possui uma carteira, utilize o comando `/carteira criar` para criar uma.",
				},
			})
			return
		}
		log.Printf("Error in get a wallet: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: errMessage,
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

	var total float64
	total, err = wallet.Invest(strings.ToUpper(cryptoCoin), quantity)
	if err != nil {
		log.Printf("Error in investing: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: errMessage,
			},
		})
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("O valor investido foi %.2fJC. Agora você tem %.8f%s imaginários", quantity, total, strings.ToUpper(cryptoCoin)),
		},
	})
	if err != nil {
		log.Printf("Error in interaction response: %v", err)
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Flags:   1 << 6,
			Content: errMessage,
		})
		return
	}

}
