package entity

import "encoding/json"

type WalletResponse struct {
	ID       int64  `json:"id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Balance  int64  `json:"balance,omitempty"`
}

type WalletReq struct {
	ID       int64  `json:"id,omitempty" swaggerignore:"true"`
	UserID   int64  `json:"user_id,omitempty" validate:"required" swaggerignore:"true"`
	UserName string `json:"user_name" validate:"required"`
	Balance  int64  `json:"balance" validate:"required"`
}

func (c *WalletReq) LoadFromMap(m map[string]interface{}) error {
	data, err := json.Marshal(m)
	if err == nil {
		err = json.Unmarshal(data, c)
	}
	return err
}

func (r *WalletReq) SetID(ID int64) {
	r.ID = ID
}

func (r *WalletReq) SetUserID(UserID int64) {
	r.UserID = UserID
}
