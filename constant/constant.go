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

	Packager = `[{"tag":"00","value":""},{"tag":"01","value":""},{"tag":"40","value":[{"tag":"00","value":""},{"tag":"01","value":""},{"tag":"02","value":""}]},{"tag":"52","value":""},{"tag":"53","value":""},{"tag":"54","value":""},{"tag":"58","value":""},{"tag":"59","value":""},{"tag":"60","value":""},{"tag":"61","value":""},{"tag":"62","value":[{"tag":"08","value":""},{"tag":"99","value":[{"tag":"00","value":""},{"tag":"01","value":""}]}]},{"tag":"63","value":""}]`
)
