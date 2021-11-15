package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joaom00/botgo/pkg/handlers"
	"github.com/joaom00/botgo/pkg/slashcommands"
	"github.com/joho/godotenv"
)

var (
	commands = []*discordgo.ApplicationCommand{
		slashcommands.ThanksCmd,
		slashcommands.WalletCmd,
		slashcommands.SalaryCmd,
		slashcommands.InvestCmd,
	}

	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"thanks":   slashcommands.ThanksHandler,
		"carteira": slashcommands.WalletHandler,
		"salario":  slashcommands.SalaryHandler,
		"investir": slashcommands.InvestHandler,
	}
)

var s *discordgo.Session

func init() {
	var err error
	godotenv.Load()
	s, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
		return
	}

}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is now running")
	})

	s.AddHandler(handlers.MessageCreate)

	if err := s.Open(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
		return
	}

	defer s.Close()

	for _, v := range commands {
		_, err := s.ApplicationCommandCreate(os.Getenv("APP_ID"), os.Getenv("GUILD_ID"), v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	sc := make(chan os.Signal)
	signal.Notify(sc, os.Interrupt)
	<-sc
	log.Println("Graceful shutdown")
}
