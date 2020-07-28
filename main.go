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
	key := os.Getenv("OSU_TOKEN") //insert your key here, has to have quotes around it.
	return key
}

func BotInit() *discordgo.Session {
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
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
			Mode: osuapi.ModeOsu,
			Limit: 5,
		})
		if err != nil {
			fmt.Println("what happened?")
		}
		if len(scores) == 0 {
			session.ChannelMessageSend(message.ChannelID, "Player " + player + " has not submitted scores in a while!")
		} else {
			msg := ""
			for _, score := range scores {
				msg = msg + fmt.Sprintf("**Map:** %d\n**Score**: %d | **PP**: %.2f\n\n", score.BeatmapID, score.Score.Score, score.Score.PP)
			}

			embed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{},
				Color: 0xff69b4,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name: "recent scoret lol",
						Value: msg,
					},
				},
			}
			session.ChannelMessageSendEmbed(message.ChannelID, embed)
		}
	}

	if strings.HasPrefix(content, "go!status") {
		if len(args) < 2 {
			session.ChannelMessageSend(message.ChannelID, "Response Status: " + strings.ToUpper(GetStatus("https://osu.ppy.sh")))
		} else {
			website := args[1]
			session.ChannelMessageSend(message.ChannelID, "Response Status: " + strings.ToUpper(GetStatus(website)))
		}
	}

	if strings.HasPrefix(content, "go!osu") {
		player := ""
		if len(args) < 2 { player = "jeesusmies" } else { player = args[1] }
		api := osuapi.NewClient(GetKey())
		stats, err := api.GetUser(osuapi.GetUserOpts{
			Username: player,
			Mode: osuapi.ModeOsu,
		})

		if err != nil { 
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("err: %s", err))
			return;
		}
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Color: 0xff69b4,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name: fmt.Sprintf("%s\nCountry: %s", stats.Username, stats.Country),
					Value: fmt.Sprintf("**pp**: %.2f\n**rank**: #%d (Country: #%d)\n**level**: %.2f\n**accuracy**: %.2f ", stats.PP, stats.Rank, stats.CountryRank, stats.Level, stats.Accuracy),
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "go!osu {user}",
				IconURL: message.Author.AvatarURL(""),

			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://a.ppy.sh/" + strconv.Itoa(stats.UserID),
				Width: 128,
				Height: 128,
			},
		}
		session.ChannelMessageSendEmbed(message.ChannelID, embed)
	}

	if strings.HasPrefix(content, "go!about") {
		session.ChannelMessageSend(message.ChannelID, "**What is this?**\nThis is a Discord bot for a game called *osu!*, with features showing player stats, recent played maps and best scores.\n**Who made this?**\nMostly `Byte#0101`, with small contributions by `jeesusmies#0500.`")
	}
}
