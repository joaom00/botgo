package slashcommands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/internal/services/database"
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
	content := ""

	switch i.ApplicationCommandData().Options[0].Name {
	case "criar":
		cont, err := createWalletHandler(i.Member.User.ID)
		if err != nil {
			log.Println(err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   1 << 6,
					Content: cont,
				},
			})
			return
		}

		content = cont
	case "saldo":
		wallet, err := database.GetWallet(i.Member.User.ID)
		if err != nil {
			log.Println(err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   1 << 6,
					Content: "Algo deu errado ao tentar acessar o seu saldo.\nPor favor, reporte o erro a um moderador",
				},
			})
			return
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Você possui %dJC", wallet.Amount),
			},
		})
	default:
		content = "Algo deu errado :("
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func createWalletHandler(userID string) (string, error) {
	_, err := database.GetWallet(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if _, err = database.CreateWallet(userID); err != nil {
				return "Algo deu errado ao tentar criar sua carteira.\nPor favor, reporte o erro a um moderador", err
			}
			return "Como é a sua primeira vez, você ganhará `100 JCoins`!", err
		}
		return "Algo deu errado ao tentar criar sua carteira.\nPor favor, reporte o erro a um moderador", err
	}

	return "Você já possui uma carteira", nil
}
