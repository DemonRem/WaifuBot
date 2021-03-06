package disc

import (
	"fmt"
	"log"
	"time"

	"github.com/Karitham/WaifuBot/database"
	"github.com/Karitham/WaifuBot/query"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"go.mongodb.org/mongo-driver/mongo"
)

// Roll drops a random character and adds it to the database
func (b *Bot) Roll(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	userData, err := database.ViewUserData(m.Author.ID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if nextRollTime := time.Until(userData.Date.Add(c.TimeBetweenRolls.Duration)); nextRollTime > 0 {
		return nil, fmt.Errorf("**illegal roll**,\nroll in %s", nextRollTime.Truncate(time.Second))
	}

	char, err := query.CharSearchByPopularity(b.seed.Uint64() % c.MaxCharacterRoll)

	if err != nil {
		return nil, err
	}
	if ok, _ := database.CharID(char.Page.Characters[0].ID).VerifyWaifu(m.Author.ID); ok {
		return b.Roll(m)
	}

	err = database.CharStruct(char).AddRolled(m.Author.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &discord.Embed{
		Title: char.Page.Characters[0].Name.Full,
		URL:   char.Page.Characters[0].SiteURL,
		Description: fmt.Sprintf(
			"You rolled character %d\nIt appears in :\n- %s",
			char.Page.Characters[0].ID, char.Page.Characters[0].Media.Nodes[0].Title.Romaji,
		),
		Thumbnail: &discord.EmbedThumbnail{
			URL: char.Page.Characters[0].Image.Large,
		},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf(
				"You can roll again in %s",
				c.TimeBetweenRolls.Truncate(time.Second),
			),
		},
	}, nil
}
