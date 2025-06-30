package supplier

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Supplier struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" validate:"required"`
	Email     *string   `json:"email,omitempty" gorm:"type:varchar(100)" validate:"omitempty,email"`
	Phone     *string   `json:"phone,omitempty" gorm:"type:varchar(20)" validate:"omitempty,max=20"`
	Address   *string   `json:"address,omitempty" gorm:"type:varchar(255)" validate:"omitempty,max=255"`
	RFC       *string   `json:"rfc,omitempty" gorm:"type:varchar(13)" validate:"omitempty,len=13"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Supplier) BeforeCreate(*gorm.DB) {
	fmt.Println("BeforeCreate: Initializing timestamps")
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}
