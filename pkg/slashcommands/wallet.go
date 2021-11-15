package slashcommands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/internal/database"
	"github.com/joaom00/botgo/pkg/embedbuilder"
	"go.mongodb.org/mongo-driver/mongo"
)

var WalletCmd = &discordgo.ApplicationCommand{
	Name:        "carteira",
	Description: "carteira subcommands",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "criar",
			Description: "Crie uma carteira para começar a investir",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
		{
			Name:        "saldo",
			Description: "Vejo o saldo da sua carteira",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		},
	},
}

func WalletHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var content string

	switch i.ApplicationCommandData().Options[0].Name {
	case "criar":
		_, err := database.GetWallet(i.Member.User.ID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				if _, err = database.CreateWallet(i.Member.User.ID); err != nil {
					content = errMessage
					break
				}
				content = "Como é a sua primeira vez, você ganhará `100 JCoins`!"
				break
			}
			log.Printf("Error in get a wallet: %v", err)
			content = errMessage
			break
		}
		content = "Você já possui uma carteira"
	case "saldo":
		wallet, err := database.GetWallet(i.Member.User.ID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				content = "Você não possui uma carteira, utilize o comando `/carteira criar` para criar uma"
				break
			}
			log.Printf("Error in get a wallet: %v", err)
			content = errMessage
			break
		}

		embed := embedbuilder.New().
			SetAuthor(i.Member.User.Username+"#"+i.Member.User.Discriminator, i.Member.User.AvatarURL("")).
			SetTitle(fmt.Sprintf("Você possui %.2fJC", wallet.Amount))

		for _, coin := range wallet.Coins {
			embed.AddField(coin.Symbol, fmt.Sprintf("%.8f", coin.Quantity))
		}

		embed.AddField("Total", fmt.Sprintf("%.2f", wallet.Total()))

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
			},
		})
	default:
		content = errMessage
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
