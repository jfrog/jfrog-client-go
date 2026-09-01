package distribution

type DistributionCommonParams struct {
	SiteName     string
	CityName     string
	CountryCodes []string
	Priority     string
}

type DistributionGetter interface {
	GetSiteName() string
	SetSiteName(siteName string)
	GetCityName() string
	SetCityName(cityName string)
	GetCountryCodes() []string
	SetCountryCodes(countryCodes []string)
	GetPriority() string
	SetPriority(priority string)
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

func (params *DistributionCommonParams) GetPriority() string {
	return params.Priority
}

func (params *DistributionCommonParams) SetPriority(priority string) {
	params.Priority = priority
}
