package bingchat_api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type IBingChat interface {
	Reset(style ...ConversationStyle)
	SendMessage(msg string) (*MsgResp, error)
	Style() ConversationStyle
	Close()
	CheckAlive() bool
	SetRemote(url, ws string)
}

type BingChatHub struct {
	sync.Mutex
	wsConn            *websocket.Conn
	cookies           []*http.Cookie
	client            *http.Client
	chatSession       *ConversationSession
	invocationId      int
	sendMessage       *SendMessage
	conversationStyle ConversationStyle
	timeout           time.Duration
	conversationUrl   string
	conversationWs    string
}

func (b *BingChatHub) buildHeaders(data map[string]string) http.Header {

	headers := http.Header{}
	for key, value := range data {
		headers.Add(key, value)
	}
	return headers
}

func NewBingChat(cookiesJson string, style ConversationStyle, timeout time.Duration) (IBingChat, error) {
	var cookies []*http.Cookie
	_ = json.Unmarshal([]byte(cookiesJson), &cookies)
	return &BingChatHub{
		timeout: timeout,
		cookies: cookies,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Timeout: timeout,
		},
		conversationStyle: style,
		conversationUrl:   conversationUrl,
		conversationWs:    conversationWs,
	}, nil
}

// Reset conversation style,the supported style list is:
// ConversationCreateStyle
// ConversationBalanceStyle
// ConversationPreciseStyle
func (b *BingChatHub) Reset(style ...ConversationStyle) {
	if b.wsConn == nil {
		return
	}
	if len(style) > 0 {
		b.conversationStyle = style[0]
	}
	_ = b.wsConn.Close()
	b.chatSession = nil
	b.invocationId = 0
	b.sendMessage = nil
}

// Close the websocket collection with new bing chat
func (b *BingChatHub) Close() {
	if b.wsConn != nil {
		_ = b.wsConn.Close()
	}
}

// SetRemote can set your own url & ws which available with new bing chat
func (b *BingChatHub) SetRemote(url, ws string) {
	if len(url) > 0 {
		b.conversationUrl = url
	}

	if len(ws) > 0 {
		b.conversationWs = ws
	}

}

func (b *BingChatHub) CheckAlive() bool {
	if b.wsConn == nil {
		return false
	}
	err := b.wsConn.WriteMessage(websocket.PingMessage, []byte{})
	if err != nil {
		return false
	}
	_, _, err = b.wsConn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			return false
		}
	}
	return true
}

// Style return current conversation style
func (b *BingChatHub) Style() ConversationStyle {
	return b.conversationStyle
}

func (b *BingChatHub) createConversation() error {
	req, err := http.NewRequest(http.MethodGet, b.conversationUrl, nil)
	if err != nil {
		return err
	}
	req.Header = b.buildHeaders(reqHeader)
	req.Header.Set("x-ms-client-request-id", uuid.New().String())
	for _, cookie := range b.cookies {
		req.AddCookie(cookie)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request status code: %d", resp.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	b.chatSession = &ConversationSession{}
	err = json.NewDecoder(resp.Body).Decode(b.chatSession)
	if err != nil {
		return err
	}
	return nil
}

func (b *BingChatHub) initWsConnect() error {
	if b.CheckAlive() {
		return nil
	}
	dial := websocket.DefaultDialer
	dial.Proxy = http.ProxyFromEnvironment
	dial.HandshakeTimeout = b.timeout
	dial.EnableCompression = true
	dial.TLSClientConfig = &tls.Config{}
	conn, resp, err := dial.Dial(b.conversationWs, b.buildHeaders(wsHeader))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	b.wsConn = conn

	err = conn.WriteMessage(websocket.BinaryMessage, []byte(`{"protocol":"json","version":1}`+DELIMITER))
	if err != nil {
		return fmt.Errorf("write json response: %v", err)
	}
	_, _, err = conn.NextReader()
	go func() {
		for {
			b.Lock()
			err := conn.WriteMessage(websocket.BinaryMessage, []byte(`{"type":6}`+DELIMITER))
			b.Unlock()
			if err != nil {
				break
			}
			time.Sleep(time.Second * 5)
		}
	}()
	return err
}

type MsgResp struct {
	Suggest    []string
	Notify     chan string
	Title      string
	References map[string]string
}

// SendMessage send message to bing chat and return a response with message(string) channel
// which you should receive the element from channel to get truly response message
func (b *BingChatHub) SendMessage(msg string) (*MsgResp, error) {
	if b.chatSession == nil {
		err := b.createConversation()
		if err != nil {
			return nil, err
		}
	}
	err := b.initWsConnect()

	if err != nil {
		return nil, err
	}

	b.beforeSendMsg(msg)
	msgData, _ := json.Marshal(b.sendMessage)
	b.Lock()
	err = b.wsConn.WriteMessage(websocket.BinaryMessage, append(msgData, []byte(DELIMITER)...))
	b.Unlock()
	if err != nil {
		return nil, err
	}
	msgRespChannel := &MsgResp{
		Notify: make(chan string, 1),
	}

	go func(msgRespChannel *MsgResp) {
		var startRev bool
		var lastMsg string
		defer close(msgRespChannel.Notify)
		for {
			_, data, err := b.wsConn.ReadMessage()
			if err != nil {
				log.Println(err)
				b.Reset()
				break
			}
			if len(data) == 0 {
				continue
			}
			spData := bytes.Split(data, []byte(DELIMITER))
			if len(spData) == 0 {
				continue
			}
			data = spData[0]
			resp := MessageResp{}
			_ = json.Unmarshal(data, &resp)

			for _, message := range resp.Item.Messages {
				if message.MessageType == "Disengaged" {
					b.Reset()
					return
				}
			}

			if resp.Type == 1 && len(resp.Arguments) > 0 && len(resp.Arguments[0].Cursor.J) > 0 {
				startRev = true
				continue
			}
			if !startRev {
				continue
			}

			//response continue
			if resp.Type == 1 && len(resp.Arguments) > 0 && len(resp.Arguments[0].Messages) > 0 {

				if len(resp.Arguments[0].Messages[0].SuggestedResponses) > 0 {
					suggests := make([]string, 0, len(resp.Arguments[0].Messages[0].SuggestedResponses))
					for _, suggest := range resp.Arguments[0].Messages[0].SuggestedResponses {
						suggests = append(suggests, suggest.Text)
					}
					msgRespChannel.Suggest = suggests
				}

				if len(resp.Arguments[0].Messages[0].SourceAttributions) > 0 {
					references := make(map[string]string)
					for _, source := range resp.Arguments[0].Messages[0].SourceAttributions {
						references[source.ProviderDisplayName] = source.SeeMoreUrl
					}
					msgRespChannel.References = references
				}

				if resp.Arguments[0].Messages[0].MessageType == "Disengaged" {
					b.Reset()
					break
				}
				msg := strings.TrimSpace(resp.Arguments[0].Messages[0].Text)
				msgRespChannel.Title = msg
				if len(lastMsg) > len(msg) {
					continue
				}
				if len(msg) == 0 || len(msg[len(lastMsg):]) == 0 {
					continue
				}
				msgRespChannel.Notify <- msg[len(lastMsg):]
				lastMsg = msg
			}

			//response done
			if resp.Type == 2 {
				break
			}
		}
	}(msgRespChannel)

	return msgRespChannel, nil
}

func (b *BingChatHub) beforeSendMsg(msg string) {
	if b.sendMessage == nil {
		b.sendMessage = b.conversationStyle.TmpMessage()
		b.sendMessage.Arguments[0].ConversationSignature = b.chatSession.ConversationSignature
		b.sendMessage.Arguments[0].Participant.Id = b.chatSession.ClientID
		b.sendMessage.Arguments[0].ConversationId = b.chatSession.ConversationID
	}
	b.sendMessage.Arguments[0].TraceId = b.getTraceId()
	b.sendMessage.Arguments[0].IsStartOfSession = b.invocationId == 0
	b.sendMessage.Arguments[0].Message.Text = msg
	b.sendMessage.Arguments[0].Message.Timestamp = time.Now()
	b.sendMessage.InvocationId = fmt.Sprint(b.invocationId)
	b.invocationId += 1
}

func (b *BingChatHub) getTraceId() string {
	//Ignore error cause that not ganna happen when use constant 1.
	node, _ := snowflake.NewNode(1)

	return node.Generate().String()

}
