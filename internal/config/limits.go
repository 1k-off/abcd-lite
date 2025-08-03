package config

type PackageLimits struct {
	MaxProjects int `mapstructure:"-"`
	MaxWebsites int `mapstructure:"-"`
}

func PersonalPlanProjects() int {
	return 1
}

func PersonalPlanWebsites() int {
	return 5
}

func LitePlanProjects() int {
	return 0
}

func LitePlanWebsites() int {
	return 0
}

func AirGappedPlanProjects() int {
	return 0
}

func AirGappedPlanWebsites() int {
	return 0
}

func editionLimits(edition string) PackageLimits {
	switch edition {
	case editionPersonal:
		return PackageLimits{
			MaxProjects: PersonalPlanProjects(),
			MaxWebsites: PersonalPlanWebsites(),
		}
	case editionLite:
		return PackageLimits{
			MaxProjects: LitePlanProjects(),
			MaxWebsites: LitePlanWebsites(),
		}
	case editionAirGapped:
		return PackageLimits{
			MaxProjects: AirGappedPlanProjects(),
			MaxWebsites: AirGappedPlanWebsites(),
		}
	default:
		return PackageLimits{
			MaxProjects: PersonalPlanProjects(),
			MaxWebsites: PersonalPlanWebsites(),
		}
	}
}
