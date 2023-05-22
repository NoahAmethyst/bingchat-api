
# Bing Chat  

### Requirements
* Have access to https://bing.com/chat
* Supported country or proxy with NewBing

### Export Bingchat Cookie
- Install [Cookie-Editor](https://chrome.google.com/webstore/detail/cookie-editor/hlkenndednhfkekhgcdicdfddnkalmdm?hl=en) in your browser.
- Export `bing.com` cookies with json

### Use
```go
go get github.com/NoahAmethyst/bingchat-api
```

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