package qristcico

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ndokoblog/qristcico/constant"
	"github.com/snksoft/crc"
)

type StructQR struct {
	Tag00 string `json:"00"`
	Tag01 string `json:"01,omitempty"`
	Tag40 struct {
		Tag00 string `json:"00"`
		Tag01 string `json:"01"`
		Tag02 string `json:"02"`
	} `json:"40"`
	Tag52 string `json:"52"`
	Tag53 string `json:"53"`
	Tag54 string `json:"54,omitempty"`
	Tag58 string `json:"58"`
	Tag59 string `json:"59"`
	Tag60 string `json:"60"`
	Tag61 string `json:"61"`
	Tag62 struct {
		Tag01 string `json:"01,omitempty"`
		Tag05 string `json:"05,omitempty"`
		Tag07 string `json:"07,omitempty"`
		Tag08 string `json:"08"`
		Tag99 string `json:"99,omitempty"`
	} `json:"62"`
	Tag63 string `json:"63"`
}

func Generate(benef, account, name, city, zip, accountId string, amount int) (stringQR, mpan string, e error) {
	var qr StructQR
	err := json.Unmarshal([]byte(constant.BaseQRBRI), &qr)
	if err != nil {
		return stringQR, mpan, fmt.Errorf("fail get base qr")
	}

	digit := constant.BenefType[benef]
	if digit == "" {
		digit = constant.BenefType["default"]
	}

	qr.Tag40.Tag01 = constant.NnsBRI + digit + padLeft(accountId, "0", 9)
	qr.Tag40.Tag02 = account

	if amount > 0 {
		qr.Tag54 = strconv.Itoa(amount)
	}

	qr.Tag59 = strings.ToUpper(name)
	qr.Tag60 = city
	qr.Tag61 = zip

	stringQR, err = qr.tlv()
	if err != nil || stringQR == "" {
		return stringQR, mpan, fmt.Errorf("fail generate")
	}

	mpan = qr.Tag40.Tag01

	return stringQR, mpan, nil
}

func Decode(str string) (qr StructQR, e error) {
	var packager []map[string]interface{}
	err := json.Unmarshal([]byte(constant.Packager), &packager)
	if err != nil {
		return qr, fmt.Errorf("fail get packager")
	}

	t, err := funcTlv(str)
	if err != nil {
		return qr, err
	}
	tag, err := parsingTag(t, packager)
	if err != nil {
		return qr, err
	}

	c := crc.CalculateCRC(crc.CCITT, []byte(str[:len(str)-4]))
	crc := fmt.Sprintf("%X", c)
	crc = padLeft(crc, "0", 4)
	if crc != tag["63"].(string) {
		return qr, fmt.Errorf("invalid crc value")
	}

	js, err := json.Marshal(tag)
	if err != nil {
		return qr, fmt.Errorf("fail marshal")
	}

	err = json.Unmarshal(js, &qr)
	if err != nil {
		return qr, fmt.Errorf("fail unmarshal")
	}

	return qr, nil
}

func (x StructQR) tlv() (s string, e error) {
	var m map[string]interface{}
	js, err := json.Marshal(x)
	if err != nil {
		return s, fmt.Errorf("fail marshal")
	}
	err = json.Unmarshal(js, &m)
	if err != nil {
		return s, fmt.Errorf("fail unmarshal")
	}

	s += Stringify(map[string]interface{}{"00": m["00"]})
	s += Stringify(map[string]interface{}{"01": m["01"]})
	s += Stringify(map[string]interface{}{"40": m["40"]})
	s += Stringify(map[string]interface{}{"52": m["52"]})
	s += Stringify(map[string]interface{}{"53": m["53"]})
	s += Stringify(map[string]interface{}{"54": m["54"]})
	s += Stringify(map[string]interface{}{"58": m["58"]})
	s += Stringify(map[string]interface{}{"59": m["59"]})
	s += Stringify(map[string]interface{}{"60": m["60"]})
	s += Stringify(map[string]interface{}{"61": m["61"]})
	s += Stringify(map[string]interface{}{"62": m["62"]})
	s += Stringify(map[string]interface{}{"63": m["63"]})

	c := crc.CalculateCRC(crc.CCITT, []byte(s[:len(s)-4]))
	crc := fmt.Sprintf("%X", c)
	crc = padLeft(crc, "0", 4)
	s = s[:len(s)-4] + crc

	return s, nil
}

func Stringify(m map[string]interface{}) (s string) {
	for k, v := range m {
		s += k
		if b, ok := v.(map[string]interface{}); ok {
			temp := Stringify(b)
			length := pad0(len(temp))
			s += length
			s += temp
		} else {
			length := pad0(len(v.(string)))
			s += length
			s += v.(string)
		}
	}
	return s
}

func pad0(i int) string {
	length := strconv.Itoa(i)
	if len(length) == 1 {
		length = "0" + length
	}
	return length
}

func funcTlv(str string) (map[string]interface{}, error) {
	index := 0
	tag := map[string]interface{}{}

	for index < len(str) {
		if index+2 > len(str) {
			return tag, fmt.Errorf("invalid length format")
		}
		de := str[index : index+2]
		if index+4 > len(str) {
			return tag, fmt.Errorf("invalid length format")
		}
		length, _ := strconv.Atoi(str[index+2 : index+4])
		if index+4+length > len(str) {
			return tag, fmt.Errorf("invalid length format")
		}
		tag[de] = str[index+4 : index+4+length]
		index = index + 4 + length
	}

	return tag, nil
}

func parsingTag(tlv map[string]interface{}, packager []map[string]interface{}) (map[string]interface{}, error) {
	tag := map[string]interface{}{}

	for i := 0; i < len(packager); i++ {
		if t, ok := packager[i]["tag"].(string); ok {
			if _, ok := tlv[t]; ok {
				var m []map[string]interface{}
				js, _ := json.Marshal(packager[i]["value"])
				err := json.Unmarshal(js, &m)
				if err == nil {
					tl, err := funcTlv(tlv[t].(string))
					if err != nil {
						return tag, err
					}
					tag[t], err = parsingTag(tl, m)
					if err != nil {
						return tag, err
					}
				} else {
					tag[t] = tlv[packager[i]["tag"].(string)]
				}
			}
		}
	}

	return tag, nil
}

func padLeft(s, pad string, l int) string {
	for len(s) < l {
		s = pad + s
	}
	return s
}
