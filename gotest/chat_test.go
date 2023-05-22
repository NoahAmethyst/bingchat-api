package gotest

import (
	"fmt"
	bingchat_api "github.com/NoahAmethyst/bingchat-api"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_Conversation(t *testing.T) {
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
			fmt.Println()
			break
		}
		respBuilder.WriteString(msg)
	}

	t.Logf("%s", respBuilder.String())

}
