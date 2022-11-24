package types

import "time"

type ValidatorsQuery struct {
	Validators []ValidatorDetail
}

type ValidatorDetail struct {
	OperatorAddress string `json:"operator_address"`
	ConsensusPubKey struct {
		AType string `json:"@type"`
		Key   string `json:"key"`
	}
	Jailed          bool   `json:"jailed"`
	Status          string `json:"status"`
	Tokens          string `json:"tokens"`
	DelegatorShares string `json:"delegator_shares"`
	Description     struct {
		Moniker         string `json:"moniker"`
		Identity        string `json:"identity"`
		Website         string `json:"website"`
		SecurityContact string `json:"security_contact"`
		Details         string `json:"details"`
	}
	UnbondingHeight string `json:"unbonding_height"`
	UnbondingTime   string `json:"unbonding_time"`
	Commission      struct {
		CommissionRates CommissionRatesStruct
		UpdateTime      time.Time `json:"update_time"`
	}
	MinSelfDelegation string `json:"min_self_delegation"`
}

type CommissionRatesStruct struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}
