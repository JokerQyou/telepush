// Introduced by jejoivanic in 5f48bd2bc154f63f7cf1afacdbde5f1d0a493fbb

package internal

import (
	"github.com/n1try/telegram-middleman-bot/config"
	"github.com/n1try/telegram-middleman-bot/model"
)

const (
	TextType = "TEXT"
	FileType = "FILE"
)

var (
	botConfig    *config.BotConfig
	textResolver = &MessageResolver{
		IsValid: validateText,
		Resolve: resolveText,
		Value:   logText,
	}
	fileResolver = &MessageResolver{
		Resolve: resolveFile,
		IsValid: validateFile,
		Value:   logFile,
	}
)

type MessageResolver struct {
	IsValid func(*model.InMessage) error
	Resolve func(string, *model.InMessage) error
	Value   func(*model.InMessage) string
}

func init() {
	botConfig = config.Get()
}

func GetResolver(ttype string) *MessageResolver {
	switch ttype {
	case FileType:
		return fileResolver
	}
	return textResolver
}