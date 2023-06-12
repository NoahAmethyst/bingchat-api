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
	question := "how is weather today in Seattle?"

	message, err := chat.SendMessage(question)
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
	t.Logf("suggest:%+v\nanswer:%+v", message.Suggest, respBuilder.String())

	t.Logf("%+v", chat.CheckAlive())

}

func Test_CheckAlive(t *testing.T) {
	chat, err := bingchat_api.NewBingChat(os.Getenv("COOKIE"), bingchat_api.ConversationBalanceStyle, 2*time.Minute)
	if err != nil {
		panic(err)
	}
	t.Logf("%+v", chat.CheckAlive())
}
