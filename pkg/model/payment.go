package model



type Payment struct {
	PaymentID     string  `gorm:"primaryKey;type:varchar(255);not null" json:"PaymentID"` // Make PaymentID the primary key if it holds the Stripe ID
	OrderID       uint    `gorm:"not null" json:"OrderID"`
	Amount        float64 `gorm:"not null" json:"PaymentAmount"`
	Status        string  `gorm:"not null" json:"Status"`
	ClientSecret  string  `gorm:"not null" json:"ClientSecret"`
	PaymentMethod string  `gorm:"not null" json:"PaymentMethod"`
	UserID        uint    `gorm:"not null" json:"UserID"`
}

// type Payment struct {
// 	gorm.Model

// 	OrderID       uint    `gorm:"not null" json:"OrderID"`
// 	Amount        float64 `gorm:"not null" json:"PaymentAmount"`
// 	Status        string  `gorm:"not null" json:"Status"`
// 	PaymentID     string  `gorm:"not null" json:"PaymentID"`
// 	ClientSecret  string  `gorm:"not null" json:"ClientSecret"`
// 	PaymentMethod string  `gorm:"not null" json:"PaymentMethod"`
// 	UserID        uint    `gorm:"not null" json:"UserID"`
// }

// type CuttingResult struct {
// 	gorm.Model
// 	OrderID    uint   `gorm:"not null"`
// 	ItemID     uint   `gorm:"not null"`
// 	CutPattern string `gorm:"type:text"`
// 	Wastage    float64
// }
