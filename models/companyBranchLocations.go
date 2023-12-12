package models

type CompanyBranchLocation struct {
	ID        uint    `gorm:"primary_key" json:"id"`
	Branch    string  `gorm:"size:255;not null;" json:"branch"`
	Address   string  `gorm:"size:255;not null;" json:"address"`
	Longitude float64 `gorm:"size:255;not null;" json:"longitude"`
	Latitude  float64 `gorm:"size:255;not null;" json:"latitude"`
}

func (c *CompanyBranchLocation) SaveCompanyBranchLocation() (*CompanyBranchLocation, error) {

	err := DB.Create(&c).Error

	if err != nil {
		return &CompanyBranchLocation{}, err
	}

	return c, nil
}

func GetCompanyBranchLocations() ([]CompanyBranchLocation, error) {

	var c []CompanyBranchLocation

	if err := DB.Find(&c).Error; err != nil {
		return c, err
	}

	return c, nil
}
