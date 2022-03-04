package alertmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/muety/telepush/config"
	"github.com/muety/telepush/inlets"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/resolvers"
	"github.com/muety/telepush/util"
)

var (
	tokenRegex = regexp.MustCompile("^Bearer (.+)$")
)

type AlertmanagerInlet struct{}

func New() inlets.Inlet {
	return &AlertmanagerInlet{}
}

func (i *AlertmanagerInlet) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m Message

		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&m); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		message := transformMessage(&m)

		ctx := r.Context()
		ctx = context.WithValue(ctx, config.KeyMessage, message)
		ctx = context.WithValue(ctx, config.KeyParams, &model.MessageParams{DisableLinkPreviews: true})

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func transformMessage(in *Message) *model.DefaultMessage {
	var sb strings.Builder
	sb.WriteString("*Alertmanager* wrote:\n\n")

	for i, a := range in.Alerts {
		// Status
		var statusEmoji string
		switch a.Status {
		case "firing":
			statusEmoji = "❗️"
			break
		case "resolved":
			statusEmoji = "✅"
		}
		sb.WriteString(fmt.Sprintf("*⌛️ Status:* %s %s\n", a.Status, statusEmoji))

		// Source URL
		sb.WriteString(fmt.Sprintf("*🔗 Source*: [Link](%s)\n", a.Url))

		// Labels
		if len(a.Labels) > 0 {
			sb.WriteString(fmt.Sprintf("*🏷 Labels:*\n"))
			for k, v := range a.Labels {
				k = util.EscapeMarkdown(k)
				v = util.EscapeMarkdown(v)
				sb.WriteString(fmt.Sprintf("– `%s` = `%s`\n", k, v))
			}
		}

		// Annotations
		if len(a.Annotations) > 0 {
			sb.WriteString(fmt.Sprintf("*📝 Annotations:*\n"))
			for k, v := range a.Annotations {
				k = util.EscapeMarkdown(k)
				v = util.EscapeMarkdown(v)
				sb.WriteString(fmt.Sprintf("– `%s` = `%s`\n", k, v))
			}
		}

		if i < len(in.Alerts)-1 {
			sb.WriteString("---\n\n")
		}
	}

	return &model.DefaultMessage{
		Text: sb.String(),
		Type: resolvers.TextType,
	}
}
