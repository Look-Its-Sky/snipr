package schemas

import "gorm.io/gorm"

type Contract struct {
	gorm.Model
	Address            string `gorm:"uniqueIndex;not null"`
	BackingCoinAddress string
	Exchange           string `gorm:"index"`
	BlockNumber        uint64 `gorm:"index"`
}
