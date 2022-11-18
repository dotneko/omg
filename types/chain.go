package types

type BalancesQuery struct {
	Balances   []DenomAmount
	Pagination PaginationStruct
}

type RewardsQuery struct {
	Rewards []ValidatorReward
	Total   []DenomAmount
}

type ValidatorReward struct {
	ValidatorAddress string `json:"validator_address"`
	Reward           []DenomAmount
}

type DenomAmount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type PaginationStruct struct {
	NextKey string `json:"next_key"`
	Total   string `json:"total"`
}

type KeysListQuery struct {
	Key []KeyStruct
}
type KeyStruct struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
	Pubkey  string `json:"pubkey"`
}