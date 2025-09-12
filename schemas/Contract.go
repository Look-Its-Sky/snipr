package schemas

import (
	"gorm.io/gorm"
)

type Contract struct {
	gorm.Model
	Address 							string `gorm:"uniqueIndex;not null"`
	BackingCoinAddress    string
}
