package nhealth

import (
	"encoding/json"
)

const (
	bodyMeasureURL = "https://api.health.nokia.com/measure"
)

// Measurement types to query
const (
	MeasureTypeWeight                 = 1
	MeasureTypeHeight                 = 4
	MeasureTypeFatFreeMass            = 5
	MeasureTypeFatRatio               = 6
	MeasureTypeFatMassWeight          = 8
	MeasureTypeDiastolicBloodPressure = 9
	MeasureTypeSystolicBloodPressure  = 10
	MeasureTypeHeartPulse             = 11
	MeasureTypeTemperature            = 12
	MeasureTypeSP02                   = 54
	MeasureTypeBodyTemperature        = 71
	MeasureTypeSkinTemperature        = 73
	MeasureTypeMuscleMass             = 76
	MeasureTypeHydration              = 77
	MeasureTypeBoneMass               = 88
	MeasureTypePulseWaveVelocity      = 91
)

// Categories for measurements
const (
	CategoryReal = 1
	CategoryGoal = 2
)

// MeasureQuery measurement query object
type MeasureQuery struct {
	Action      string `url:"action,omitempty"`
	UserID      uint32 `url:"userid,omitempty"`
	StartDate   uint32 `url:"startdate,omitempty"`
	EndDate     uint32 `url:"enddate,omitempty"`
	LastUpdate  uint32 `url:"lastupdate,omitempty"`
	MeasureType uint32 `url:"meastype,omitempty"`
	Category    uint32 `url:"category,omitempty"`
	Limit       uint32 `url:"limit,omitempty"`
	Offset      uint32 `url:"offset,omitempty"`
}

// MeasureResponse measurement response object
type MeasureResponse struct {
	UpdateTime  uint32 `json:"updatetime"`
	TimeZone    string `json:"timezone"`
	More        uint32 `json:"more"`
	MeasureGrps []struct {
		GroupID     string `json:"grpid"`
		Attribution string `json:"attrib"`
		Date        uint32 `json:"date"`
		Category    uint32 `json:"category"`
		Measures    []struct {
			Value uint32 `json:"value"`
			Unit  uint32 `json:"unit"`
			Type  uint32 `json:"type"`
		} `json:"measures"`
	} `json:"measuregrps"`
	Status  uint32 `json:"status"`
	Offset  uint32 `json:"offset"`
	Message string `json:"message"`
}

//http://api.health.nokia.com/measure?action=getmeas&oauth_consumer_key=49c99e905b3901800b96e070c88f278087fb5465d452b38046b7c7dd4b2c7f2&oauth_nonce=013d75812a4ae4fab0bd96a519b58425&oauth_signature=XPx1NMUYGccS6bkbXhdXuJNnVSI%3D&oauth_signature_method=HMAC-SHA1&oauth_timestamp=1504159331&oauth_token=d3758d01eb95e19bf0c751844b2fcbb79ffe0d834738d39d0669751b71&oauth_version=1.0&userid=13644360

// GetMeasurement fetches a measurement for a given user
func (hapi *HealthAPI) GetMeasurement(accessToken, accessSecret string, mq MeasureQuery) (*MeasureResponse, error) {
	mq.Action = "getmeas"

	resp, err := hapi.Get(accessSecret, accessSecret, bodyMeasureURL, &mq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	inst := MeasureResponse{}
	err = json.NewDecoder(resp.Body).Decode(&inst)
	if err != nil {
		return nil, err
	}

	return &inst, nil
}
