package disc

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Karitham/WaifuBot/database"
	"github.com/Karitham/WaifuBot/query"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
)

// Roll drops a random character and adds it to the database
func (b *Bot) Roll(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	userData, err := database.ViewUserData(m.Author.ID)
	if err != nil {
		return nil, err
	}

	if nextRollTime := time.Until(userData.Date.Add(c.TimeBetweenRolls.Duration)); nextRollTime > 0 {
		return nil, fmt.Errorf("cannot roll now, roll in %s", nextRollTime.Truncate(time.Second))
	}

	char, err := query.CharSearchByPopularity(
		rand.New(
			rand.NewSource(
				time.Now().UnixNano(),
			),
		).Intn(c.MaxCharacterRoll),
	)
	if err != nil {
		return nil, err
	}

	database.InputChar{
		UserID: m.Author.ID,
		CharList: database.CharLayout{
			ID:    char.Page.Characters[0].ID,
			Image: char.Page.Characters[0].Image.Large,
			Name:  char.Page.Characters[0].Name.Full,
		},
		Date: time.Now(),
	}.AddChar()

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
				"You can roll again at %s",
				m.ID.Time().Add(c.TimeBetweenRolls.Duration).Format("15:04"),
			),
		},
	}, nil
}
