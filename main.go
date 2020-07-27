package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"net/http"
	"strconv"
	osuapi "github.com/thehowl/go-osuapi"

)
func GetKey() string {
	key := "" //insert your key here, has to have quotes around it.
	return key
}

func BotInit() *discordgo.Session {
	discord, err := discordgo.New("Bot " + ""//discord token here)
	if err != nil {
		fmt.Println("Error initialising the bot!")
	}
	return discord
}
func main() {
	bot := BotInit()
	// Register the messageCreate func as a callback for MessageCreate events.
	bot.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	bot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err := bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}

func GetStatus(s string) string {
	resp, err := http.Get(s)
    if err != nil {
        panic(err)
	}
	
	return resp.Status
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	content := message.Content
	args := strings.Fields(content)

	if message.Author.ID == session.State.User.ID {
		return
	}
	
	if message.Content == "osu!" {
		session.ChannelMessageSend(message.ChannelID, "Bad game B)")
	}

	if strings.HasPrefix(content, "go!recent") {
		player := args[1]
		api := osuapi.NewClient(GetKey())
		scores, err := api.GetUserRecent(osuapi.GetUserScoresOpts{
			Username: "WhiteCat",
			Mode: osuapi.ModeOsu,
		})
		if err != nil {
			fmt.Println("what happened?")
		}
		session.ChannelMessageSend(message.ChannelID, "Player : " + player)
		for _, score := range scores {
			session.ChannelMessageSend(message.ChannelID, "Full Combo : " + strconv.FormatBool(bool(score.Score.FullCombo)))
			session.ChannelMessageSend(message.ChannelID, "PP : " + strconv.Itoa(int(score.Score.Score)))
			session.ChannelMessageSend(message.ChannelID, "PP : " + strconv.Itoa(int(score.Score.PP)) + "\n")

		}
	}

	if strings.HasPrefix(content, "go!status") {
		website := args[1]
		session.ChannelMessageSend(message.ChannelID, "Response Status: " + strings.ToUpper(GetStatus(website)))
	}
}
