package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	osuapi "github.com/thehowl/go-osuapi"
)

var (
	scoreAmount string
	BeatmapID   string
	fullCombo   string
	PP          string
)

func GetKey() string {
	key := os.Args[1] //insert your key here, has to have quotes around it.
	if key == "" {
		fmt.Println("Enter a key as first argument!")
	}
	return key
}

func BotInit() *discordgo.Session {
	discord, err := discordgo.New("Bot " + os.Args[2])
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

func OutputAll(fullcombo string, scoreAmount string, beatmapid string, pp string) string {
	output := "Full Combo: " + fullcombo + "\nScore: " + scoreAmount + "\nBeatmap: https://osu.ppy.sh/b/" + beatmapid + "\nPP: " + pp + "\n"
	return output
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
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
			Username: player,
			Mode:     osuapi.ModeOsu,
			Limit:    10,
		})
		if err != nil {
			fmt.Println("what happened?")
		}
		if len(scores) == 0 {
			session.ChannelMessageSend(message.ChannelID, "Player " + player + " has not submitted scores in a while!")
		} else {
			var pog string = ""
			for _, score := range scores {
				fullCombo = strconv.FormatBool(bool(score.Score.FullCombo))
				BeatmapID = strconv.Itoa(score.BeatmapID)
				scoreAmount = strconv.Itoa(int(score.Score.Score))
				PP = strconv.Itoa(int(score.Score.PP))
				pog = pog + OutputAll(fullCombo, scoreAmount, BeatmapID, PP)
			}
			session.ChannelMessageSend(message.ChannelID, player + "\n" + pog)
		}
	}

	if strings.HasPrefix(content, "go!best") {
		player := args[1]
		api := osuapi.NewClient(GetKey())
		scores, err := api.GetUserBest(osuapi.GetUserScoresOpts{
			Username: player,
			Mode:     osuapi.ModeOsu,
			Limit: 5,
		})
		if err != nil {
			fmt.Println("what happened?")
		}
		var pog string = ""
		for _, score := range scores {
			fullCombo = strconv.FormatBool(bool(score.Score.FullCombo))
			BeatmapID = strconv.Itoa(score.BeatmapID)
			scoreAmount = strconv.Itoa(int(score.Score.Score))
			PP = strconv.Itoa(int(score.Score.PP))
			pog = pog + OutputAll(fullCombo, scoreAmount, BeatmapID, PP)
		}
		session.ChannelMessageSend(message.ChannelID, player + "\n" + pog)
		
	}

	if strings.HasPrefix(content, "go!status") {
		if len(args) < 2 {
			session.ChannelMessageSend(message.ChannelID, "Response Status For osu!: "+strings.ToUpper(GetStatus("https://osu.ppy.sh")))
		} else {
			website := args[1]
			session.ChannelMessageSend(message.ChannelID, "Response Status: "+strings.ToUpper(GetStatus(website)))
		}
	}
}
