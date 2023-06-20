package entity

type Wallet struct {
	Balance  int64  `gorm:"column:balance"`
	UserName string `gorm:"column:user_name"`
	ID       int64  `gorm:"column:id"`
}

func (Wallet) TableName() string {
	return "wallets"
}
