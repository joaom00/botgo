package slashcommands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/internal/services/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var SalaryCmd = &discordgo.ApplicationCommand{
	Name:        "salario",
	Description: "Receba as suas 30 JCoins do dia",
}

func SalaryHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	wallet, err := database.AddSalary(i.Member.User.ID)
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
		log.Println(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Algo deu errado :(.\nPor favor, reporte o erro a um moderador",
			},
		})
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Você ganhou 30JC! Agora você possui %f JCoins", wallet.Amount),
		},
	})
	if err != nil {
		log.Panicln(err)
		s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Flags:   1 << 6,
			Content: "Algo deu errado :(\nPor favor, reporte a um moderador",
		})
		return
	}
}
