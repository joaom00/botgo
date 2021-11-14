package slashcommands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

const errMessage = "Algo deu errado :(.\nPor favor, reporte o erro a um moderador"

var SalaryCmd = &discordgo.ApplicationCommand{
	Name:        "salario",
	Description: "Receba as suas 30 JCoins do dia",
}

func SalaryHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var content string

	wallet, err := database.GetWallet(i.Member.User.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			content = "Você não possui uma carteira, utilize o comando `/carteira criar` para criar uma"
		}
		log.Printf("Error in get a wallet: %v", err)
		content = errMessage
	}

	if err = wallet.AddSalary(); err != nil {
		log.Printf("Error in add salary: %v", err)
		content = errMessage
	}

	content = fmt.Sprintf("Você ganhou 30JC! Agora você possui %.2f JCoins", wallet.Amount)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
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
