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
	Num00 string `json:"00"`
	Num01 string `json:"01,omitempty"`
	Num40 struct {
		Num00 string `json:"00"`
		Num01 string `json:"01"`
		Num02 string `json:"02"`
	} `json:"40"`
	Num52 string `json:"52"`
	Num53 string `json:"53"`
	Num54 string `json:"54,omitempty"`
	Num58 string `json:"58"`
	Num59 string `json:"59"`
	Num60 string `json:"60"`
	Num61 string `json:"61"`
	Num62 struct {
		Num08 string `json:"08"`
		Num99 struct {
			Num00 string `json:"00,omitempty"`
			Num01 string `json:"01,omitempty"`
		} `json:"99,omitempty"`
	} `json:"62"`
	Num63 string `json:"63"`
}

func Generate(benef, account, name, city, zip, refnum string) (stringQR string, e error) {
	var qr StructQR
	err := json.Unmarshal([]byte(constant.BaseQRBRI), &qr)
	if err != nil {
		return stringQR, fmt.Errorf("fail get base qr")
	}

	qr.Num40.Num01 = constant.NnsBRI + constant.BenefType[benef] + refnum[len(refnum)-10:]
	qr.Num40.Num02 = account
	qr.Num59 = strings.ToUpper(name)
	qr.Num60 = city
	qr.Num61 = zip

	stringQR, err = qr.tlv()
	if err != nil || stringQR == "" {
		return stringQR, fmt.Errorf("fail generate")
	}

	return stringQR, nil
}

func Decode(str string) (qr StructQR, e error) {
	var packager []map[string]interface{}
	err := json.Unmarshal([]byte(constant.Packager), &packager)
	if err != nil {
		return qr, fmt.Errorf("fail get packager")
	}

	tag := parsingTag(funcTlv(str), packager)

	c := crc.CalculateCRC(crc.CCITT, []byte(str[:len(str)-4]))
	if fmt.Sprintf("%X", c) != tag["63"].(string) {
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

	s += stringify(map[string]interface{}{"00": m["00"]})
	s += stringify(map[string]interface{}{"01": m["01"]})
	s += stringify(map[string]interface{}{"40": m["40"]})
	s += stringify(map[string]interface{}{"52": m["52"]})
	s += stringify(map[string]interface{}{"53": m["53"]})
	s += stringify(map[string]interface{}{"54": m["54"]})
	s += stringify(map[string]interface{}{"58": m["58"]})
	s += stringify(map[string]interface{}{"59": m["59"]})
	s += stringify(map[string]interface{}{"60": m["60"]})
	s += stringify(map[string]interface{}{"61": m["61"]})
	s += stringify(map[string]interface{}{"62": m["62"]})
	s += stringify(map[string]interface{}{"63": m["63"]})

	c := crc.CalculateCRC(crc.CCITT, []byte(s[:len(s)-4]))
	s = s[:len(s)-4] + fmt.Sprintf("%X", c)

	return s, nil
}

func stringify(m map[string]interface{}) (s string) {
	for k, v := range m {
		s += k
		if b, ok := v.(map[string]interface{}); ok {
			temp := stringify(b)
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

func funcTlv(str string) map[string]interface{} {
	index := 0
	tag := map[string]interface{}{}

	for index < len(str) {
		de := str[index : index+2]
		length, _ := strconv.Atoi(str[index+2 : index+4])
		tag[de] = str[index+4 : index+4+length]
		index = index + 4 + length
	}

	return tag
}

func parsingTag(tlv map[string]interface{}, packager []map[string]interface{}) map[string]interface{} {
	tag := map[string]interface{}{}

	for i := 0; i < len(packager); i++ {
		if t, ok := packager[i]["tag"].(string); ok {
			if _, ok := tlv[t]; ok {
				var m []map[string]interface{}
				js, _ := json.Marshal(packager[i]["value"])
				err := json.Unmarshal(js, &m)
				if err == nil {
					tag[t] = parsingTag(funcTlv(tlv[t].(string)), m)
				} else {
					tag[t] = tlv[packager[i]["tag"].(string)]
				}
			}
		}
	}

	return tag
}
