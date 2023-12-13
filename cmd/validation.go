package cmd

import (
	"fmt"
	"strings"
	"unicode"
)

func validInfo(p Page) error {

	if err := validNumField(p.IdInstance); err != nil {
		return fmt.Errorf("поле для IdInstance %s", err)
	}
	if len(p.IdInstance) != 10 {
		return fmt.Errorf("длина IdInstance должна быть 10 цифр\n")
	}

	if len(p.ApiTokenInstance) != 50 {
		return fmt.Errorf("длина ApiTokenInstance должна быть 50 символов\n")
	}

	var legalChar = &unicode.RangeTable{
		R16: []unicode.Range16{
			{0x0030, 0x0039, 1}, // 0-9
			{0x0041, 0x0046, 1}, // A-F
			{0x0061, 0x0066, 1}, // a-f
		},
		R32: []unicode.Range32{},
	}

	for _, letter := range p.ApiTokenInstance {
		if !unicode.In(letter, legalChar) {
			return fmt.Errorf("недопустимый символ [%s] в поле для ApiTokenInstance\nдопустимы [0123456789abcdef]\n",
				string(letter))
		}
	}

	return nil
}

func validMsg(send sendMsg, msg msg) error {

	if err := validInfo(Page{IdInstance: send.IdInstance, ApiTokenInstance: send.ApiTokenInstance}); err != nil {
		return err
	}

	if err := validNumField(send.MsgNum); err != nil {
		return fmt.Errorf("поле номера получателя %s", err)
	}

	chat := strings.Split(msg.ChatId, "@")
	if len(chat) != 2 ||
		chat[1] == "c.us" && chat[0] != send.MsgNum ||
		chat[1] != "c.us" && chat[1] != "g.us" {
		return fmt.Errorf("неверно указан получатель: %s", msg.ChatId)
	}

	return nil
}

func validNumField(num string) error {

	if strings.ContainsAny(num, "-+e") {
		return fmt.Errorf("должно содержать только цифры\n")
	}
	return nil
}

func validFile(send sendFile, prop propFile) error {

	prop.UrlFile = send.FPath

	if err := validInfo(Page{IdInstance: send.IdInstance, ApiTokenInstance: send.ApiTokenInstance}); err != nil {
		return err
	}

	if err := validNumField(send.FileNum); err != nil {
		return fmt.Errorf("поле номера получателя %s", err)
	}

	return nil
}
