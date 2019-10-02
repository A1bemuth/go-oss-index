package types

type OssIndexRequest struct {
	Coordinates []string `json:"coordinates"`
}

type Vulnerability struct {
	Id            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	CvssScore     float32  `json:"cvssScore"`
	CvssVector    string   `json:"cvssVector"`
	Cwe           string   `json:"cwe"`
	Cve           string   `json:"cve"`
	Reference     string   `json:"reference"`
	VersionRanges []string `json:"versionRanges"`
}

type ComponentReport struct {
	Coordinates     string          `json:"coordinates"`
	Description     string          `json:"description"`
	Reference       string          `json:"reference"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
}
