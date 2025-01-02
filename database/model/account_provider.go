package model

type AccountProvider struct {
	BaseModel
	UserID      string  `db:"user_id" json:"user_id"`
	Provider    string  `db:"provider" json:"provider"`
	AccessToken string  `db:"access_token" json:"access_token"`
	IdToken     *string `db:"id_token" json:"id_token"`
}
