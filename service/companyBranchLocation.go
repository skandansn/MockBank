package service

import (
	"github.com/skandansn/webDevBankBackend/entity"
	"github.com/skandansn/webDevBankBackend/models"
	"log"
	"math"
)

type CompanyBranchLocationService interface {
	GetCompanyBranchLocations() ([]entity.CompanyBranchLocation, error)
	SaveCompanyBranchLocation(branch entity.CompanyBranchLocationInput) (entity.CompanyBranchLocation, error)
	GetNClosestBranches(n uint, latitude float64, longitude float64) ([]entity.CompanyBranchLocation, error)
}

type companyBranchLocationService struct {
}

func (c *companyBranchLocationService) GetNClosestBranches(n uint, latitude float64, longitude float64) ([]entity.CompanyBranchLocation, error) {
	dbBranches, err := models.GetCompanyBranchLocations()
	if err != nil {
		return []entity.CompanyBranchLocation{}, err
	}

	var branches []entity.CompanyBranchLocation

	for _, branch := range dbBranches {
		branches = append(branches, entity.CompanyBranchLocation{
			ID:        branch.ID,
			Branch:    branch.Branch,
			Address:   branch.Address,
			Longitude: branch.Longitude,
			Latitude:  branch.Latitude,
		})
	}

	branchDistance := make([]interface{}, len(branches))

	for i, branch := range branches {
		branchDistance[i] = map[string]interface{}{
			"branch":   branch,
			"distance": calculateDistance(latitude, longitude, branch.Latitude, branch.Longitude),
		}
	}

	// Sort by distance

	for i := 0; i < len(branchDistance); i++ {
		for j := 0; j < len(branchDistance)-i-1; j++ {
			if branchDistance[j].(map[string]interface{})["distance"].(float64) > branchDistance[j+1].(map[string]interface{})["distance"].(float64) {
				branchDistance[j], branchDistance[j+1] = branchDistance[j+1], branchDistance[j]
			}
		}
	}

	// Get n closest branches

	var closestBranches []entity.CompanyBranchLocation

	numberOfBranches := min(len(branchDistance), int(n))

	for i := 0; i < numberOfBranches; i++ {
		log.Println(n)
		closestBranches = append(closestBranches, branchDistance[i].(map[string]interface{})["branch"].(entity.CompanyBranchLocation))
	}

	return closestBranches, nil
}

func (c *companyBranchLocationService) GetCompanyBranchLocations() ([]entity.CompanyBranchLocation, error) {
	dbBranches, err := models.GetCompanyBranchLocations()
	if err != nil {
		return []entity.CompanyBranchLocation{}, err
	}

	return convertToCompanyBranchLocationDTO(dbBranches), nil
}

func (c *companyBranchLocationService) SaveCompanyBranchLocation(branch entity.CompanyBranchLocationInput) (entity.CompanyBranchLocation, error) {
	dbBranch := models.CompanyBranchLocation{
		Branch:    branch.Branch,
		Address:   branch.Address,
		Longitude: branch.Longitude,
		Latitude:  branch.Latitude,
	}

	res, err := dbBranch.SaveCompanyBranchLocation()
	if err != nil {
		return entity.CompanyBranchLocation{}, err
	}

	return entity.CompanyBranchLocation{
		ID:        res.ID,
		Branch:    res.Branch,
		Address:   res.Address,
		Longitude: res.Longitude,
		Latitude:  res.Latitude,
	}, nil
}

func NewCompanyBranchLocationService() CompanyBranchLocationService {
	return &companyBranchLocationService{}
}

func convertToCompanyBranchLocationDTO(dbBranches []models.CompanyBranchLocation) []entity.CompanyBranchLocation {
	var branches []entity.CompanyBranchLocation
	for _, branch := range dbBranches {
		branches = append(branches, entity.CompanyBranchLocation{
			ID:        branch.ID,
			Branch:    branch.Branch,
			Address:   branch.Address,
			Longitude: branch.Longitude,
			Latitude:  branch.Latitude,
		})
	}
	return branches
}

func calculateDistance(customerLat, customerLong, branchLat, branchLong float64) float64 {
	// Convert latitude and longitude from degrees to radians
	customerLatRad := toRadians(customerLat)
	customerLongRad := toRadians(customerLong)
	branchLatRad := toRadians(branchLat)
	branchLongRad := toRadians(branchLong)

	// Haversine formula
	dLat := branchLatRad - customerLatRad
	dLong := branchLongRad - customerLongRad

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(customerLatRad)*math.Cos(branchLatRad)*math.Pow(math.Sin(dLong/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Earth radius in miles
	earthRadius := 3959.0

	// Calculate distance
	distance := earthRadius * c

	return distance
}

func toRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180.0)
}
