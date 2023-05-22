package gotest

import (
	bingchat_api "github.com/NoahAmethyst/bingchat-api"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_ConversationStream(t *testing.T) {
	chat, err := bingchat_api.NewBingChat(os.Getenv("COOKIE"), bingchat_api.ConversationBalanceStyle, 2*time.Minute)
	if err != nil {
		panic(err)
	}
	message, err := chat.SendMessage("how is the weather today in Seattle")
	if err != nil {
		panic(err)
	}
	var respBuilder strings.Builder
	for {
		msg, ok := <-message.Notify
		if !ok {
			break
		}
		respBuilder.WriteString(msg)
	}

	t.Logf("%+v", message.Suggest)

	t.Logf("%s", respBuilder.String())

}
