
# Bing Chat  

*WARN Due to the use of robot check in bingChat, the api is temporarily unuseful.*

### Requirements
* Have access to https://bing.com/chat
* Supported country or proxy with NewBing

### Export Bingchat Cookie
- Install [Cookie-Editor](https://chrome.google.com/webstore/detail/cookie-editor/hlkenndednhfkekhgcdicdfddnkalmdm?hl=en) in your browser.
- Export `bing.com` cookies with json

### Warn
This project currently **support** parallel sessions (multi sessions) and context

**Please make sure that websocket be closed at the end of conversation**
### Use
```go
go get github.com/NoahAmethyst/bingchat-api
```

### Example

You can see example codes in [chat_test.go](gotest%2Fchat_test.go) 
which include conversation with context and multi conversations 

### Test
```go

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

    defer func() {
        chat.Close()
        t.Logf("%+v", chat.CheckAlive())
    }()
	
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

```