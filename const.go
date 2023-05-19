package bingchat_api

import "encoding/json"

var wsHeader map[string]string
var reqHeader map[string]string

const (
	balanceStyle    = `{"arguments":[{"source":"cib","optionsSets":["nlu_direct_response_filter","deepleo","disable_emoji_spoken_text","responsible_ai_policy_235","enablemm","galileo","rcallowlist","responseos","jb090","jbfv202","dv3sugg"],"allowedMessageTypes":["Chat","InternalSearchQuery","InternalSearchResult","Disengaged","InternalLoaderMessage","RenderCardRequest","AdsQuery","SemanticSerp","GenerateContentQuery","SearchQuery"],"sliceIds":["contctxp2tf","delayglobjscf","0417bicunivs0","ssoverlap50","sspltop5","sswebtop1","audseq","sbsvgopt","nopreloadsstf","winlongmsg2tf","perfimpcomb","sugdivdis","sydnoinputt","wpcssopt","414suggs0","scctl","418glpv6ps0","417rcallow","321slocs0","407pgparsers0","0329resp","asfixescf","udscahrfoncf","414jbfv202"],"verbosity":"verbose","traceId":"6441f4712452428ab53b745af65c5089","isStartOfSession":true,"message":{"locale":"zh-CN","market":"zh-CN","region":"WW","location":"lat:47.639557;long:-122.128159;re=1000m;","locationHints":[{"country":"Singapore","timezoneoffset":8,"countryConfidence":8,"Center":{"Latitude":1.2929,"Longitude":103.8547},"RegionType":2,"SourceType":1}],"timestamp":"2023-04-21T10:27:06+08:00","author":"user","inputMethod":"Keyboard","text":"我需要帮助制定计划","messageType":"Chat"},"conversationSignature":"AdctQf6lU2LbVhUuyTMDCcchrXUbEaX2jyBbeQ2iXuY=","participant":{"id":"914798353003051"},"conversationId":"51D|BingProd|92EBE954F84335CBF36EE4E86BF8A028765E8922A38BC8A1E248457AF7342CA3"}],"invocationId":"","target":"chat","type":4}`
	createStyle     = `{"arguments":[{"source":"cib","optionsSets":["nlu_direct_response_filter","deepleo","disable_emoji_spoken_text","responsible_ai_policy_235","enablemm","h3imaginative","rcallowlist","responseos","jb090","jbfv202","dv3sugg","clgalileo","gencontentv3"],"allowedMessageTypes":["Chat","InternalSearchQuery","InternalSearchResult","Disengaged","InternalLoaderMessage","RenderCardRequest","AdsQuery","SemanticSerp","GenerateContentQuery","SearchQuery"],"sliceIds":["contctxp2tf","delayglobjscf","0417bicunivs0","ssoverlap50","sspltop5","sswebtop1","audseq","sbsvgopt","nopreloadsstf","winlongmsg2tf","perfimpcomb","sugdivdis","sydnoinputt","wpcssopt","414suggs0","scctl","418glpv6ps0","417rcallow","321slocs0","407pgparsers0","0329resp","asfixescf","udscahrfoncf","414jbfv202"],"verbosity":"verbose","traceId":"6441f4712452428ab53b745af65c5089","isStartOfSession":true,"message":{"locale":"zh-CN","market":"zh-CN","region":"WW","location":"lat:47.639557;long:-122.128159;re=1000m;","locationHints":[{"country":"Singapore","timezoneoffset":8,"countryConfidence":8,"Center":{"Latitude":1.2929,"Longitude":103.8547},"RegionType":2,"SourceType":1}],"timestamp":"2023-04-21T10:27:06+08:00","author":"user","inputMethod":"Keyboard","text":"告诉我的星座","messageType":"Chat"},"conversationSignature":"uk3kLopdE2Zb8nTXHFx/smV2IWyec3G11B0y8ehSC4k=","participant":{"id":"914798353003051"},"conversationId":"51D|BingProd|29F20DF6A2946BAD80F2B98E87138C4E14A6D8C5D29E95056F41D5BE3539D4B4"}],"invocationId":"","target":"chat","type":4}`
	preciseStyle    = `{"arguments":[{"source":"cib","optionsSets":["nlu_direct_response_filter","deepleo","disable_emoji_spoken_text","responsible_ai_policy_235","enablemm","h3precise","rcallowlist","responseos","jb090","jbfv202","dv3sugg","clgalileo"],"allowedMessageTypes":["Chat","InternalSearchQuery","InternalSearchResult","Disengaged","InternalLoaderMessage","RenderCardRequest","AdsQuery","SemanticSerp","GenerateContentQuery","SearchQuery"],"sliceIds":["contctxp2tf","delayglobjscf","0417bicunivs0","ssoverlap50","sspltop5","sswebtop1","audseq","sbsvgopt","nopreloadsstf","winlongmsg2tf","perfimpcomb","sugdivdis","sydnoinputt","wpcssopt","414suggs0","scctl","418glpv6ps0","417rcallow","321slocs0","407pgparsers0","0329resp","asfixescf","udscahrfoncf","414jbfv202"],"verbosity":"verbose","traceId":"6441f4712452428ab53b745af65c5089","isStartOfSession":true,"message":{"locale":"zh-CN","market":"zh-CN","region":"WW","location":"lat:47.639557;long:-122.128159;re=1000m;","locationHints":[{"country":"Singapore","timezoneoffset":8,"countryConfidence":8,"Center":{"Latitude":1.2929,"Longitude":103.8547},"RegionType":2,"SourceType":1}],"timestamp":"2023-04-21T10:27:06+08:00","author":"user","inputMethod":"Keyboard","text":"我需要帮助做研究","messageType":"Chat"},"conversationSignature":"1F8e/oVRPtqkMq+/hrKWphxvXbc5DTQTsItUsoaxedE=","participant":{"id":"914798353003051"},"conversationId":"51D|BingProd|23B9F05272D0D7471D94F332A995F6996B9E192846CB8B7007092B4B6DE6FDEC"}],"invocationId":"14","target":"chat","type":4}`
	DELIMITER       = "\x1e"
	conversationUrl = "https://www.bing.com/turing/conversation/create"
	conversationWs  = "wss://sydney.bing.com/sydney/ChatHub"
)

type ConversationStyle uint8

const (
	ConversationCreateStyle ConversationStyle = iota + 1
	ConversationBalanceStyle
	ConversationPreciseStyle
)

func (c ConversationStyle) String() string {
	switch c {
	case ConversationBalanceStyle:
		return "Balance"
	case ConversationCreateStyle:
		return "Create"
	case ConversationPreciseStyle:
		return "Precise"
	}
	return ""
}

func (c ConversationStyle) TmpMessage() *SendMessage {
	var data string
	switch c {
	case ConversationBalanceStyle:
		data = balanceStyle
	case ConversationCreateStyle:
		data = createStyle
	case ConversationPreciseStyle:
		data = preciseStyle
	}
	msg := SendMessage{}
	json.Unmarshal([]byte(data), &msg)
	return &msg
}

func init() {
	wsHeader = map[string]string{
		"authority":                   "edgeservices.bing.com",
		"accept":                      "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language":             "en-US,en;q=0.9",
		"cache-control":               "max-age=0",
		"sec-ch-ua":                   `"Chromium";v="110", "Not A(Brand";v="24", "Microsoft Edge";v="110"`,
		"sec-ch-ua-arch":              `"x86"`,
		"sec-ch-ua-bitness":           `"64"`,
		"sec-ch-ua-full-version":      `"110.0.1587.69"`,
		"sec-ch-ua-full-version-list": `"Chromium";v="110.0.5481.192", "Not A(Brand";v="24.0.0.0", "Microsoft Edge";v="110.0.1587.69"`,
		"sec-ch-ua-mobile":            "?0",
		"sec-ch-ua-model":             `""`,
		"sec-ch-ua-platform":          `"Windows"`,
		"sec-ch-ua-platform-version":  `"15.0.0"`,
		"sec-fetch-dest":              "document",
		"sec-fetch-mode":              "navigate",
		"sec-fetch-site":              "none",
		"sec-fetch-user":              "?1",
		"upgrade-insecure-requests":   "1",
		"user-agent":                  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.69",
		"x-edge-shopping-flag":        "1",
	}

	reqHeader = map[string]string{
		"accept":                      "application/json",
		"accept-language":             "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"sec-ch-ua":                   "\"Chromium\";v=\"112\", \"Microsoft Edge\";v=\"112\", \"Not:A-Brand\";v=\"99\"",
		"sec-ch-ua-arch":              "\"x86\"",
		"sec-ch-ua-bitness":           "\"64\"",
		"sec-ch-ua-full-version":      "\"112.0.1722.48\"",
		"sec-ch-ua-full-version-list": "\"Chromium\";v=\"112.0.5615.121\", \"Microsoft Edge\";v=\"112.0.1722.48\", \"Not:A-Brand\";v=\"99.0.0.0\"",
		"sec-ch-ua-mobile":            "?0",
		"sec-ch-ua-model":             "\"\"",
		"sec-ch-ua-platform":          "\"Windows\"",
		"sec-ch-ua-platform-version":  "\"15.0.0\"",
		"sec-fetch-dest":              "empty",
		"sec-fetch-mode":              "cors",
		"sec-fetch-site":              "same-origin",
		"sec-ms-gec":                  "673E82A42CAB0AF8C4F97398D164CA4F1F69BEC0D5E41226FD5375F12B17F341",
		"sec-ms-gec-version":          "1-112.0.1722.48",
		"x-ms-client-request-id":      "8f0c8a85-28bb-49eb-9c5d-c35dfc5112dd",
		"x-ms-useragent":              "azsdk-js-api-client-factory/1.0.0-beta.1 core-rest-pipeline/1.10.0 OS/Win32",
		"user-agent":                  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48",
	}

}
