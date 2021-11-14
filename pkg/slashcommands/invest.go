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
	var content string

	cryptoCoin := i.ApplicationCommandData().Options[0].StringValue()
	quantity := i.ApplicationCommandData().Options[1].FloatValue()

	if !utils.Find(utils.Coins, cryptoCoin) {
		content = fmt.Sprintf("A moeda `%s` não existe. Use `/moedas` para conhecer as moedas válidas", cryptoCoin)
	}

	wallet, err := database.GetWallet(i.Member.User.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			content = "Você não possui uma carteira, utilize o comando `/carteira criar` para criar uma."
		}
		log.Printf("Error in get a wallet: %v", err)
		content = errMessage
	}

	if wallet.Amount < quantity {
		content = "Você não possui essa quantidade para investir"
	}

	var total float64
	total, err = wallet.Invest(strings.ToUpper(cryptoCoin), quantity)
	if err != nil {
		log.Printf("Error in investing: %v", err)
		content = errMessage
	}

	content = fmt.Sprintf("O valor investido foi %.2fJC. Agora você tem %.8f%s imaginários", quantity, total, cryptoCoin)

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
