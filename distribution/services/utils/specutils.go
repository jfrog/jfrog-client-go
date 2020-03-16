package utils

type DistributionCommonParams struct {
	SiteName     string
	CityName     string
	CountryCodes []string
}

type DistributionGetter interface {
	GetSiteName() string
	SetSiteName(siteName string)
	GetCityName() string
	SetCityName(cityName string)
	GetCountryCodes() []string
	SetCountryCodes(countryCodes []string)
}

func (params *DistributionCommonParams) GetSiteName() string {
	return params.SiteName
}

func (params *DistributionCommonParams) SetSiteName(siteName string) {
	params.SiteName = siteName
}

func (params *DistributionCommonParams) GetCityName() string {
	return params.CityName
}

func (params *DistributionCommonParams) SetCityName(cityName string) {
	params.CityName = cityName
}

func (params *DistributionCommonParams) GetCountryCodes() []string {
	return params.CountryCodes
}

func (params *DistributionCommonParams) SetCountryCodes(countryCodes []string) {
	params.CountryCodes = countryCodes
}
