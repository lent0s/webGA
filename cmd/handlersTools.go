package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type Page struct {
	HeadTitle        string `json:"HeadTitle"`
	IdInstance       string `json:"IdInstance"`
	ApiTokenInstance string `json:"ApiTokenInstance"`
	MsgNum           string `json:"MsgNum"`
	Msg              string `json:"Msg"`
	FileNum          string `json:"FileNum"`
	FPath            string `json:"FPath"`
	Answer           string `json:"Answer"`
}

func getDefaultPage() *Page {

	return &Page{
		HeadTitle:        "tools4green-API",
		IdInstance:       "7123456789",
		ApiTokenInstance: "",
		MsgNum:           "77771234567",
		Msg:              "Введите текст Вашего сообщения",
		FileNum:          "77771234567",
		FPath:            "Введите путь к файлу",
		Answer:           `Введите Ваши [IdInstance]&[ApiTokenInstance]`,
	}
}

type idInstance struct {
	Wid                               string `json:"wid"`
	CountryInstance                   string `json:"countryInstance"`
	TypeAccount                       string `json:"typeAccount"`
	WebhookUrl                        string `json:"webhookUrl"`
	WebhookUrlToken                   string `json:"webhookUrlToken"`
	DelaySendMessagesMilliseconds     int    `json:"delaySendMessagesMilliseconds"`
	MarkIncomingMessagesReaded        string `json:"markIncomingMessagesReaded"`
	MarkIncomingMessagesReadedOnReply string `json:"markIncomingMessagesReadedOnReply"`
	SharedSession                     string `json:"sharedSession"`
	OutgoingWebhook                   string `json:"outgoingWebhook"`
	OutgoingMessageWebhook            string `json:"outgoingMessageWebhook"`
	OutgoingAPIMessageWebhook         string `json:"outgoingAPIMessageWebhook"`
	IncomingWebhook                   string `json:"incomingWebhook"`
	DeviceWebhook                     string `json:"deviceWebhook"`
	StatusInstanceWebhook             string `json:"statusInstanceWebhook"`
	StateWebhook                      string `json:"stateWebhook"`
	EnableMessagesHistory             string `json:"enableMessagesHistory"`
	KeepOnlineStatus                  string `json:"keepOnlineStatus"`
	PollMessageWebhook                string `json:"pollMessageWebhook"`
	IncomingBlockWebhook              string `json:"incomingBlockWebhook"`
}

type apiTokenInstance struct {
	StateInstance string `json:"stateInstance"`
}

type sendMsg struct {
	IdInstance       string `json:"IdInstance"`
	ApiTokenInstance string `json:"ApiTokenInstance"`
	MsgNum           string `json:"MsgNum"`
	Msg              string `json:"Msg"`
}

type msg struct {
	ChatId          string `json:"chatId"`
	Message         string `json:"message"`
	QuotedMessageId string `json:"quotedMessageId"`
	LinkPreview     bool   `json:"linkPreview"`
}

type msgAnswer struct {
	IdMessage string `json:"idMessage"`
}

type sendFile struct {
	IdInstance       string `json:"IdInstance"`
	ApiTokenInstance string `json:"ApiTokenInstance"`
	FileNum          string `json:"FileNum"`
	FPath            string `json:"FPath"`
}

type propFile struct {
	ChatId          string `json:"chatId"`
	UrlFile         string `json:"urlFile"`
	FileName        string `json:"fileName"`
	Caption         string `json:"caption"`
	QuotedMessageId string `json:"quotedMessageId"`
}

func runTemplate(p *Page, w http.ResponseWriter) {

	content := []string{
		`web/sh.page.tmpl`,
	}

	t := template.Must(template.ParseFiles(content...))
	err := t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func getAccountInfo(body *io.ReadCloser) ([]byte, error) {

	buf, err := getAccountData(body, "getSettings")
	if err != nil {
		return nil, err
	}

	answer := idInstance{}
	if err = json.Unmarshal(buf, &answer); err != nil {
		return nil, err
	}

	buf, err = json.MarshalIndent(answer, "", "    ")
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func getAccountState(body *io.ReadCloser) ([]byte, error) {

	buf, err := getAccountData(body, "getStateInstance")
	if err != nil {
		return nil, err
	}

	answer := apiTokenInstance{}
	if err = json.Unmarshal(buf, &answer); err != nil {
		return nil, err
	}

	buf, err = json.MarshalIndent(answer, "", "    ")
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func getAccountData(body *io.ReadCloser, urlFrag string) ([]byte, error) {

	buf, err := io.ReadAll(*body)
	if err != nil {
		return nil, err
	}

	p := Page{}
	if err = json.Unmarshal(buf, &p); err != nil {
		return nil, err
	}

	if err = validInfo(p); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(`https://api.green-api.com/waInstance%s/%s/%s`,
		p.IdInstance, urlFrag, p.ApiTokenInstance)

	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	buf, err = io.ReadAll(res.Body)
	_ = res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if buf, err = readErrorAnswer(buf); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", buf)
	}

	return buf, nil
}

func sendMessageT(body *io.ReadCloser) ([]byte, error) {

	buf, err := io.ReadAll(*body)
	if err != nil {
		return nil, err
	}

	send := sendMsg{}
	if err = json.Unmarshal(buf, &send); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(`https://api.green-api.com/waInstance%s/sendMessage/%s`,
		send.IdInstance, send.ApiTokenInstance)

	msg := msg{}
	if err = json.Unmarshal([]byte(send.Msg), &msg); err != nil {
		msg.ChatId = send.MsgNum + "@c.us"
		msg.Message = send.Msg
	}

	if err = validMsg(send, msg); err != nil {
		return nil, err
	}

	buf, err = json.MarshalIndent(msg, "", "    ")
	if err != nil {
		return nil, err
	}

	req, err := http.Post(path, "application/json",
		strings.NewReader(string(buf)))
	if err != nil {
		return nil, err
	}

	buf, err = io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	_ = req.Body.Close()

	if req.StatusCode != http.StatusOK {
		if buf, err = readErrorAnswer(buf); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", buf)
	}

	answer := msgAnswer{}
	if err = json.Unmarshal(buf, &answer); err != nil {
		return nil, err
	}

	buf, err = json.MarshalIndent(answer, "", "    ")
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func sendFileByUrlT(body *io.ReadCloser) ([]byte, error) {

	buf, err := io.ReadAll(*body)
	if err != nil {
		return nil, err
	}

	send := sendFile{}
	if err = json.Unmarshal(buf, &send); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(`https://api.green-api.com/waInstance%s/sendFileByUrl/%s`,
		send.IdInstance, send.ApiTokenInstance)

	str := strings.Trim(send.FPath, `\`)
	prop := propFile{}
	if err = json.Unmarshal([]byte(str), &prop); err != nil {
		return nil, err
	}

	if err = validFile(send, prop); err != nil {
		return nil, err
	}

	buf, err = json.MarshalIndent(prop, "", "    ")
	if err != nil {
		return nil, err
	}

	req, err := http.Post(path, "application/json",
		strings.NewReader(string(buf)))
	if err != nil {
		return nil, err
	}

	buf, err = io.ReadAll(req.Body)
	_ = req.Body.Close()
	if err != nil {
		return nil, err
	}

	if req.StatusCode != http.StatusOK {
		if buf, err = readErrorAnswer(buf); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", buf)
	}

	answer := msgAnswer{}
	if err = json.Unmarshal(buf, &answer); err != nil {
		return nil, err
	}

	buf, err = json.MarshalIndent(answer, "", "    ")
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func readErrorAnswer(buf []byte) ([]byte, error) {

	m := make(map[string]interface{})
	if err := json.Unmarshal(buf, &m); err != nil {
		return nil, err
	}

	buf, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func notSupportedMethod(request, support string, w http.ResponseWriter) bool {

	if request != support {
		w.Header().Set("Support", support)
		err := fmt.Sprintf("method \"%s\" not allowed (%s)", request, support)
		http.Error(w, err, http.StatusMethodNotAllowed)
		return true
	}
	return false
}
