package models

type UserData struct {
	Age    int32 `db:"age" json:"age"`
	Height int32 `db:"height" json:"height"`
	Female bool  `db:"female" json:"female"`
}

type CardInfo struct {
	CardType  int     `json:"card_type"`
	Available float64 `json:"available"`
	BankCard  string  `json:"bank_card"`
}
