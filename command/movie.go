package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/igungor/ilber/bot"
	"github.com/igungor/telegram"
)

func init() {
	register(cmdMovie)
}

var cmdMovie = &Command{
	Name:      "imdb",
	ShortLine: "ay-em-dii-bii",
	Run:       runMovie,
}

func runMovie(ctx context.Context, b *bot.Bot, msg *telegram.Message) {
	args := msg.Args()
	opts := &telegram.SendOptions{ParseMode: telegram.ModeMarkdown}
	if len(args) == 0 {
		term := randChoice(movieExamples)
		txt := fmt.Sprintf("hangi filmi arıyorsun? örneğin: */imdb %s*", term)
		_, err := b.SendMessage(msg.Chat.ID, txt, opts)
		if err != nil {
			log.Printf("Error while sending message: %v\n", err)
		}
		return
	}

	// the best search engine is still google.
	// i've tried imdb, themoviedb, rottentomatoes, omdbapi.
	// themoviedb search engine was the most accurate yet still can't find any
	// result if any release date is given in query terms.
	urls, err := search(b.Config.GoogleAPIKey, b.Config.GoogleSearchEngineID, "", args...)
	if err != nil {
		log.Printf("Error searching %v: %v\n", args, err)
		if err == errSearchQuotaExceeded {
			_, _ = b.SendMessage(msg.Chat.ID, `¯\_(ツ)_/¯`, opts)
		}
		return
	}

	for _, url := range urls {
		if strings.Contains(url, "imdb.com/title/tt") {
			_, err := b.SendMessage(msg.Chat.ID, url, opts)
			if err != nil {
				log.Printf("Error while sending message. Err: %v\n", err)
			}
			return
		}
	}

	_, err = b.SendMessage(msg.Chat.ID, "aradığın filmi bulamadım 🙈", opts)
	if err != nil {
		log.Printf("Error while sending message. Err: %v\n", err)
		return
	}
}

var movieExamples = []string{
	"Spirited Away",
	"Mulholland Dr",
	"Oldboy",
	"Interstellar",
	"12 Angry Men",
	"Cidade de Deus",
	"The Big Lebowski",
	"There Will Be Blood",
	"Ghost in the Shell",
	"The Grey",
	"Seven Samurai",
}
