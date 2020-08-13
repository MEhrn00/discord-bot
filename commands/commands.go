package commands

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type HelpCmd struct {
	cmd  string
	help string
}

var (
	helpmsg = []HelpCmd{
		{
			cmd:  "help",
			help: "Display this help message",
		},
		{
			cmd:  "ping",
			help: "Pong",
		},
		{
			cmd:  "kanye",
			help: "Print a quote by Kanye West",
		},
		{
			cmd:  "trump",
			help: "Print a quote by Tronald Dump",
		},
		{
			cmd:  "joke",
			help: "Print a random joke",
		},
		{
			cmd:  "status [text]",
			help: "Update the bot's status to the specified text",
		},
	}
)

func HelpHandler(s *discordgo.Session, m *discordgo.Message, Prefix string) {

	var msgEmbed discordgo.MessageEmbed

	for _, element := range helpmsg {
		msgEmbed.Description += Prefix + element.cmd + ": " + element.help + "\n\n"
	}

	msgEmbed.Title = "Commands List"
	msgEmbed.Color = 0

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &msgEmbed)
	if err != nil {
		log.Println("Error help command list,", err)
	}
}

func ErrorHandler(s *discordgo.Session, m *discordgo.Message, cmd string, Prefix string) {
	s.ChannelMessageSend(m.ChannelID, "Error command `"+cmd+"` is not a valid command use `"+Prefix+"help`")
}

func KanyeHandler(s *discordgo.Session, m *discordgo.Message) {
	type jsonQuote struct {
		Quote string `json:"quote"`
	}

	r, err := http.Get("https://api.kanye.rest/")

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error making request")
		return
	}

	defer r.Body.Close()
	contents, _ := ioutil.ReadAll(r.Body)
	var parsed jsonQuote

	err = json.Unmarshal(contents, &parsed)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error parsing Json")
		return
	}

	var msgEmbed discordgo.MessageEmbed
	msgEmbed.Title = "Quote"
	msgEmbed.Description = "\"" + parsed.Quote + "\""
	msgEmbed.Color = 0
	msgEmbed.Author = &discordgo.MessageEmbedAuthor{Name: "Kanye West", IconURL: "https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Ftse1.mm.bing.net%2Fth%3Fid%3DOIP.WdtzTAoKtJvsKdNXuRWY9QHaH6%26pid%3DApi&f=1"}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &msgEmbed)
	if err != nil {
		log.Println("Error sending kanye quote,", err)
	}

}

func TrumpHandler(s *discordgo.Session, m *discordgo.Message) {
	type jsonQuote struct {
		AppearedAt time.Time `json:"appeared_at"`
		CreatedAt  time.Time `json:"created_at"`
		QuoteID    string    `json:"quote_id"`
		Tags       []string  `json:"tags"`
		UpdatedAt  time.Time `json:"updated_at"`
		Value      string    `json:"value"`
		Embedded   struct {
			Author []struct {
				AuthorID  string      `json:"author_id"`
				Bio       interface{} `json:"bio"`
				CreatedAt time.Time   `json:"created_at"`
				Name      string      `json:"name"`
				Slug      string      `json:"slug"`
				UpdatedAt time.Time   `json:"updated_at"`
				Links     struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
			} `json:"author"`
			Source []struct {
				CreatedAt     time.Time   `json:"created_at"`
				Filename      interface{} `json:"filename"`
				QuoteSourceID string      `json:"quote_source_id"`
				Remarks       interface{} `json:"remarks"`
				UpdatedAt     time.Time   `json:"updated_at"`
				URL           string      `json:"url"`
				Links         struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
			} `json:"source"`
		} `json:"_embedded"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	}

	r, err := http.NewRequest("GET", "https://api.tronalddump.io/random/quote", nil)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error making request")
		return
	}

	r.Header.Set("Accept", "application/hal+json")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(r)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error reading request response")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error reading body")
		return
	}
	var parsed jsonQuote

	err = json.Unmarshal(body, &parsed)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error parsing json")
		return
	}

	var msgEmbed discordgo.MessageEmbed
	msgEmbed.Title = "Tweet"
	msgEmbed.Description = "\"" + parsed.Value + "\""
	msgEmbed.URL = parsed.Embedded.Source[0].URL
	msgEmbed.Color = 0
	msgEmbed.Author = &discordgo.MessageEmbedAuthor{Name: "Mr. President", IconURL: "https://s.abcnews.com/images/Politics/donald-trump-ap-jpo-181222_hpMain_4x3_992.jpg"}
	msgEmbed.Footer = &discordgo.MessageEmbedFooter{Text: "Posted"}
	msgEmbed.Timestamp = parsed.AppearedAt.Format("2006-01-02 15:04:05")
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &msgEmbed)
	if err != nil {
		log.Panicln("Error sending embed, ", err)
	}
}

func JokeHandler(s *discordgo.Session, m *discordgo.Message) {
	type jokeData struct {
		Error    bool   `json:"error"`
		Category string `json:"category"`
		Type     string `json:"type"`
		Setup    string `json:"setup"`
		Delivery string `json:"delivery"`
		Flags    struct {
			Nsfw      bool `json:"nsfw"`
			Religious bool `json:"religious"`
			Political bool `json:"political"`
			Racist    bool `json:"racist"`
			Sexist    bool `json:"sexist"`
		} `json:"flags"`
		ID   int    `json:"id"`
		Lang string `json:"lang"`
	}

	apiUrl := "https://sv443.net/jokeapi/v2/joke/Any?type=twopart"

	r, err := http.Get(apiUrl)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error making request")
		return
	}

	defer r.Body.Close()
	contents, _ := ioutil.ReadAll(r.Body)
	var parsedData jokeData

	err = json.Unmarshal(contents, &parsedData)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error parsing Json")
		return
	}

	var msgEmbed discordgo.MessageEmbed
	msgEmbed.Title = parsedData.Setup
	msgEmbed.Description = parsedData.Delivery
	msgEmbed.Color = 16711680
	msgEmbed.Author = &discordgo.MessageEmbedAuthor{Name: "Joke"}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &msgEmbed)
	if err != nil {
		log.Panicln("Error sending embed, ", err)
	}

}

func StatusHandler(s *discordgo.Session, m *discordgo.Message, status string) {
	err := s.UpdateStatus(0, status)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed updating the bot's status")
		log.Panicln("Failed updating status,", err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Updated the bot's status to `"+status+"`")
}
