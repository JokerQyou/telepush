package checkchan

import (
	"context"
	"fmt"
	"github.com/muety/telepush/config"
	"github.com/muety/telepush/inlets"
	"github.com/muety/telepush/model"
	"github.com/muety/telepush/resolvers"
	"net/http"
	"net/url"
)

type CheckChanWebhookInlet struct{}

func New() inlets.Inlet {
	return &CheckChanWebhookInlet{}
}

func parseForm(form url.Values) (*Payload, error) {
	var payload Payload
	payload.ID = form.Get("id")
	payload.URL = form.Get("url")
	payload.Value = form.Get("value")
	payload.Link = form.Get("link")
	if html := form.Get("html"); html != "" {
		payload.HTML = &html
	}
	return &payload, nil
}

func (i *CheckChanWebhookInlet) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(8192); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		payload, err := parseForm(r.Form)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		message := &model.DefaultMessage{
			Origin: "Check酱",
			Text:   fmt.Sprintf("[%s](%s)", payload.Value, payload.Link),
			Type:   resolvers.TextType,
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, config.KeyMessage, message)
		ctx = context.WithValue(ctx, config.KeyParams, &model.MessageParams{DisableLinkPreviews: true})

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
