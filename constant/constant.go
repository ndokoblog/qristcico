package constant

var (
	BaseQRBRI = `{		
		"00":"01",
		"01":"11",
		"40": {
		  "00":"id.co.bri",
		  "01":"",
		  "02":""
		},
		"52":"4829",
		"53":"360",
		"58":"ID",
		"59":"",
		"60":"",
		"61":"",
		"62":{
		  "08":"DMCT",
		  "99":{
			  "00":"",
			  "01":""
		  }
		},
		"63":"CDEF"
	  }`

	NnsBRI    = "93600014"
	BenefType = map[string]string{
		"SA": "10",
		"CA": "20",
		"UE": "60",
	}
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
