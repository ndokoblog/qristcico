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
		Tag08 string `json:"08"`
		Tag99 struct {
			Tag00 string `json:"00,omitempty"`
			Tag01 string `json:"01,omitempty"`
		} `json:"99,omitempty"`
	} `json:"62"`
	Tag63 string `json:"63"`
}

func Generate(benef, account, name, city, zip, refnum string) (stringQR string, e error) {
	var qr StructQR
	err := json.Unmarshal([]byte(constant.BaseQRBRI), &qr)
	if err != nil {
		return stringQR, fmt.Errorf("fail get base qr")
	}

	qr.Tag40.Tag01 = constant.NnsBRI + constant.BenefType[benef] + refnum[len(refnum)-10:]
	qr.Tag40.Tag02 = account
	qr.Tag59 = strings.ToUpper(name)
	qr.Tag60 = city
	qr.Tag61 = zip

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
		if va, ok := v.([]map[string]interface{}); ok {
			temp := ""
			for _, vb := range va {
				temp += Stringify(vb)
			}
			return temp
		} else if vc, ok := v.(map[string]interface{}); ok {
			return Stringify(vc)
		} else if vd, ok := v.(string); ok {
			return k + pad0(len(vd)) + vd
		}
	}
	return ""
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

func padLeft(s, pad string, l int) string {
	for len(s) < l {
		s = pad + s
	}
	return s
}
