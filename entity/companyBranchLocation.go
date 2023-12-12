package entity

type CompanyBranchLocation struct {
	ID        uint    `gorm:"primary_key" json:"id"`
	Branch    string  `gorm:"size:255;not null;" json:"branch"`
	Address   string  `gorm:"size:255;not null;" json:"address"`
	Longitude float64 `gorm:"size:255;not null;" json:"longitude"`
	Latitude  float64 `gorm:"size:255;not null;" json:"latitude"`
}

type CompanyBranchLocationInput struct {
	Branch    string  `json:"branch" binding:"required"`
	Address   string  `json:"address" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}
