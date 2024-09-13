package chainalysis

import (
	"encoding/json"
	"fmt"
	"time"
)

type Direction string

const (
	DirectionSent    Direction = "sent"
	DirectionReceive Direction = "received"
)

// API Objects and Types
// Anything which is a pointer is nullable

type TransferRegisterReq struct {
	Network           string    `json:"network"`
	Asset             string    `json:"asset"`
	TransferReference string    `json:"transferReference"`
	Direction         Direction `json:"direction"`
	// Below are optional fields
	TransferTimestamp string   `json:"transferTimestamp,omitempty"`
	AssetAmount       float64  `json:"assetAmount,omitempty"`
	OutputAddress     string   `json:"outputAddress,omitempty"`
	InputAddresses    []string `json:"inputAddresses,omitempty"`
	AssetPrice        float64  `json:"assetPrice,omitempty"`
	AssetDenomination string   `json:"assetDenomination,omitempty"`
}

func (t TransferRegisterReq) Validate() error {
	if t.Network == "" {
		return fmt.Errorf("network is required")
	}
	if t.Asset == "" {
		return fmt.Errorf("asset is required")
	}
	if t.TransferReference == "" {
		return fmt.Errorf("transferReference is required")
	}
	if t.Direction == "" {
		return fmt.Errorf("direction is required")
	}
	return nil
}

type TransferRegisterResp struct {
	UpdatedAt         *string     `json:"updatedAt"`
	Asset             string      `json:"asset"`
	Network           string      `json:"network"`
	TransferReference string      `json:"transferReference"` //{transaction_hash}:{log_index}
	Tx                *string     `json:"tx"`
	Idx               *int        `json:"idx"`
	UsdAmount         json.Number `json:"usdAmount"`
	AssetAmount       json.Number `json:"assetAmount"`
	Timestamp         *string     `json:"timestamp"`
	OutputAddress     *string     `json:"outputAddress"`
	ExternalID        string      `json:"externalId"`
}

func (t TransferRegisterResp) GetUpdatedAt() (time.Time, error) {
	if t.UpdatedAt == nil {
		return time.Time{}, fmt.Errorf("transfer has not been processed yet")
	}
	return time.Parse("2006-01-02T15:04:05.000000", *t.UpdatedAt)
}

func (t TransferRegisterResp) KYTFinishedProcessing() bool {
	return t.UpdatedAt != nil
}

type ExposureType string

const (
	IndirectExposure ExposureType = "INDIRECT"
	DirectExposure   ExposureType = "DIRECT"
)

type AlertType string

func (a AlertType) Int() int {
	switch a {
	case AlertSevere:
		return 4
	case AlertHigh:
		return 3
	case AlertMEDIUM:
		return 2
	case AlertLOW:
		return 1
	default:
		return 0
	}
}

const (
	AlertSevere AlertType = "SEVERE"
	AlertHigh   AlertType = "HIGH"
	AlertMEDIUM AlertType = "MEDIUM"
	AlertLOW    AlertType = "LOW"
)

type Exposure struct {
	Name       string `json:"name"`
	Category   string `json:"category"`
	CategoryID int    `json:"categoryId"`
}

// Equal compares two exposures
func (e Exposure) Equal(other Exposure) bool {
	return e.CategoryID == other.CategoryID && e.Name == other.Name
}

type Alert struct {
	AlertLevel   AlertType    `json:"alertLevel"`
	Category     *string      `json:"category"`
	Service      string       `json:"service"`
	ExternalID   string       `json:"externalId"`
	AlertAmount  json.Number  `json:"alertAmount"`  // amount of $USD involved in the alert
	ExposureType ExposureType `json:"exposureType"` // direct or indirect exposure
	CategoryID   int          `json:"categoryId"`
}
