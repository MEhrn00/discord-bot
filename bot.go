package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/MEhrn00/discord-bot/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token  string
	Prefix string
)

func init() {
	godotenv.Load()
	Token = strings.TrimSuffix(os.Getenv("DISCORD_TOKEN"), "\n")
	Prefix = strings.TrimSuffix(os.Getenv("PREFIX_CHARACTER"), "\n")

	if Prefix == "" {
		Prefix = "!"
	}
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Panicln("Error creating discrod session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(ready)
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		log.Panicln("Error opening connection,", err)
		return
	}

	log.Println("Bot is running. Press Ctrl-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	err := s.UpdateStatus(0, "IDK what to do with this bot....")
	if err != nil {
		log.Panicln("Failed updating status,", err)
		return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, Prefix) {
		command := strings.TrimPrefix(m.Content, Prefix)
		command = strings.Split(command, " ")[0]

		switch command {
		case "ping":
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Pong! %d ms", s.HeartbeatLatency()/1000000))
		case "help":
			commands.HelpHandler(s, m.Message, Prefix)
		case "kanye":
			commands.KanyeHandler(s, m.Message)
		case "trump":
			commands.TrumpHandler(s, m.Message)
		case "joke":
			commands.JokeHandler(s, m.Message)
		case "status":
			status := strings.TrimPrefix(m.Content, Prefix+"status ")
			commands.StatusHandler(s, m.Message, status)
		default:
			commands.ErrorHandler(s, m.Message, command, Prefix)
		}
	}

}
