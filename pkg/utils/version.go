package utils

var WellKnownVersions = map[string]string{
	"1.4":  `^1\.4.*`,
	"1.5":  `^1\.5.*`,
	"1.6":  `^1\.6.*`,
	"1.7":  `^1\.7.*`,
	"1.8":  `^1\.8.*`,
	"1.9":  `^1\.9.*`,
	"1.10": `^1\.10.*`,
	"1.11": `^1\.11.*`,
	"1.12": `^1\.12.*`,
	"1.13": `^1\.13.*`,
	"1.14": `^1\.14.*`,
	"1.15": `^1\.15.*`,
	"1.16": `^1\.16.*`,
	"1.17": `^1\.17.*`,
	"1.18": `^1\.18.*`,
	"1.19": `^1\.19.*`,
	"1.20": `^1\.20.*`,
	"1.21": `^1\.21.*`,
	"1.22": `^1\.22.*`,
	"1.23": `^1\.23.*`,
	"1.24": `^1\.24.*`,
	"1.25": `^1\.25.*`,
	"1.26": `^1\.26.*`,
}

func BuildEnvoyFilterNamesAllVersion(base string) []string {
	var names []string
	for version := range WellKnownVersions {
		names = append(names, base+"-"+version)
	}

	return names
}

func BuildEnvoyFilterNames(base string, versions []string) []string {
	var names []string
	for _, version := range versions {
		names = append(names, base+"-"+version)
	}

	return names
}
