package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/skandansn/webDevBankBackend/entity"
	"github.com/skandansn/webDevBankBackend/service"
	"github.com/skandansn/webDevBankBackend/utils"
)

type CompanyBranchLocationController interface {
	GetCompanyBranchLocations(ctx *gin.Context) ([]entity.CompanyBranchLocation, error)
	SaveCompanyBranchLocation(ctx *gin.Context) (entity.CompanyBranchLocation, error)
}

type companyBranchLocationController struct {
	service service.CompanyBranchLocationService
}

func (c *companyBranchLocationController) GetCompanyBranchLocations(ctx *gin.Context) ([]entity.CompanyBranchLocation, error) {

	lat := ctx.Query("lat")
	long := ctx.Query("long")
	count := ctx.Query("count")
	countInt := utils.ParseStringAsInt(count)
	float64Lat := utils.ParseStringAsFloat64(lat)
	float64Long := utils.ParseStringAsFloat64(long)

	if lat == "" || long == "" {
		res, err := c.service.GetCompanyBranchLocations()
		if err != nil {
			return []entity.CompanyBranchLocation{}, err
		}
		return res, nil
	}

	if count == "" {
		countInt = 3
	}

	res, err := c.service.GetNClosestBranches(countInt, float64Lat, float64Long)

	if err != nil {
		return []entity.CompanyBranchLocation{}, err
	}

	return res, nil

}

func (c *companyBranchLocationController) SaveCompanyBranchLocation(ctx *gin.Context) (entity.CompanyBranchLocation, error) {

	var branch entity.CompanyBranchLocationInput

	if err := ctx.ShouldBindJSON(&branch); err != nil {
		return entity.CompanyBranchLocation{}, err
	}

	res, err := c.service.SaveCompanyBranchLocation(branch)
	if err != nil {
		return entity.CompanyBranchLocation{}, err
	}
	return res, nil
}

func NewCompanyBranchLocationController(service service.CompanyBranchLocationService) CompanyBranchLocationController {
	return &companyBranchLocationController{
		service: service,
	}
}
