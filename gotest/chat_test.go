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
	questions := [3]string{"how is weather today in Nanjing?", "How to fry a steak", "Who was the winner of the last World Cup?"}

	for _, _question := range questions {

		message, err := chat.SendMessage(_question)
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
	}

}
