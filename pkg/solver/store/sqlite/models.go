package store

import (
	"gorm.io/gorm"
	"github.com/lilypad-tech/lilypad/pkg/data"
)

// this is what the solver keeps track of so we can know
// what the current state of the deal is
type JobOfferData struct {
	gorm.Model
	CID         string   `json:"cid"`
	DealID     string   `json:"deal_id"`
	JobCreator string   `json:"job_creator"`
	State      uint8    `json:"state"`
	JobOffer   data.JobOffer `json:"job_offer"`
}

// this is what the solver keeps track of so we can know
// what the current state of the deal is
type ResourceOfferData struct {
	gorm.Model
	CID               string				`json:"cid"`
	DealID           string        			`json:"deal_id"`
	ResourceProvider string        			`json:"resource_provider"`
	State            uint8         			`json:"state"`
	ResourceOffer    data.ResourceOffer		`json:"job_offer"`
}

type DealData struct {
	gorm.Model
	CID               string				`json:"cid"`
	JobCreator       string           		`json:"job_creator"`
	ResourceProvider string           		`json:"resource_provider"`
	JobOffer         string           		`json:"job_offer"`
	ResourceOffer    string           		`json:"resource_offer"`
	State            uint8            		`json:"state"`
	Deal             data.Deal				`json:"deal"`
	Transactions     data.DealTransactions	`json:"transactions"`
	Mediator         string					`json:"mediator"`
}

type ResultData struct {
	gorm.Model
	// this is the cid of the result where ID is set to empty string
	CID     string			`json:"cid"`
	DealID string			`json:"deal_id"`
	// the CID of the actual results
	DataID           string `json:"results_id"`
	Error            string `json:"error"`
	InstructionCount uint64 `json:"instruction_count"`
}

// represents a solver decision
// the solver keeps track of "no" decisions to avoid trying to repeatedly match
// things it's already decided it can't match
type MatchDecisionData struct {
	gorm.Model
	JobOffer      string `json:"job_offer"`
	ResourceOffer string `json:"resource_offer"`
	Deal          string `json:"deal"`
	Result        bool   `json:"result"`
}
