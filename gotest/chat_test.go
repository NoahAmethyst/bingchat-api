package gotest

import (
	bingchat_api "github.com/NoahAmethyst/bingchat-api"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_ConversationStream(t *testing.T) {
	chat, err := bingchat_api.NewBingChat(os.Getenv("COOKIE"), bingchat_api.ConversationBalanceStyle, 2*time.Minute)
	if err != nil {
		panic(err)
	}

	defer func() {
		chat.Close()
		t.Logf("%+v", chat.CheckAlive())
	}()
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

	alive := chat.CheckAlive()
	t.Logf("%+v", alive)

	if !alive {
		return
	}

	question = "how about Washington"
	message, err = chat.SendMessage(question)
	if err != nil {
		panic(err)
	}
	respBuilder.Reset()
	for {
		msg, ok := <-message.Notify
		if !ok {
			break
		}
		respBuilder.WriteString(msg)
	}
	t.Logf("suggest:%+v\nanswer:%+v", message.Suggest, respBuilder.String())

}

func Test_MultiConversation(t *testing.T) {
	var wait sync.WaitGroup
	for i := 0; i < 3; i++ {
		wait.Add(1)
		go func() {
			chat, err := bingchat_api.NewBingChat(os.Getenv("COOKIE"), bingchat_api.ConversationBalanceStyle, 2*time.Minute)
			if err != nil {
				panic(err)
			}

			defer func() {
				wait.Done()
				chat.Close()
				t.Logf("%+v", chat.CheckAlive())
			}()
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
		}()
	}
	wait.Wait()
}

func Test_CheckAlive(t *testing.T) {
	chat, err := bingchat_api.NewBingChat(os.Getenv("COOKIE"), bingchat_api.ConversationBalanceStyle, 2*time.Minute)
	if err != nil {
		panic(err)
	}
	t.Logf("%+v", chat.CheckAlive())
}
