package models

import "github.com/jinzhu/gorm"

type Access struct {
	gorm.Model
	EmployeeID  uint   `gorm:"not null" json:"employee_id"`
	AccessName  string `gorm:"size:255;not null;" json:"access_name"`
	AccessGiven bool   `gorm:"not null;" json:"access_given"`
}

func GetAllAccessForEmployee(eid uint) ([]Access, error) {

	var a []Access

	if err := DB.Where("employee_id = ?", eid).Find(&a).Error; err != nil {
		return a, err
	}

	return a, nil
}

func SaveAccessesForEmployee(a []Access) ([]Access, error) {
	for i := range a {
		var access Access
		if err := DB.Where(Access{EmployeeID: a[i].EmployeeID, AccessName: a[i].AccessName}).First(&access).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				err = DB.Create(&a[i]).Error
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			access.AccessGiven = a[i].AccessGiven
			err = DB.Save(&access).Error
			if err != nil {
				return nil, err
			}
		}
	}
	return a, nil
}

func IsAccessPresentForEmployee(eid uint, accessName string) (bool, error) {

	var a Access

	if err := DB.Where("employee_id = ? AND access_name = ? AND access_given = 1", eid, accessName).First(&a).Error; err != nil {
		return false, err
	}

	return true, nil
}
