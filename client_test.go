package ossindex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindsPackagesWithVulnerabilities(t *testing.T) {
	t.Run(
		"Single package",
		verifyReturnsRequestedPackages([]string{"pkg:npm/url-parse@1.4.2"}))
	t.Run(
		"Two packages",
		verifyReturnsRequestedPackages([]string{"pkg:npm/url-parse@1.4.2", "pkg:npm/macaddress@0.2.8"}))
}

func verifyReturnsRequestedPackages(purls []string) func(t *testing.T) {
	return func(t *testing.T) {
		client := Client{}
		require := require.New(t)
		reports, err := client.Get(purls)

		require.Nil(err)
		require.Len(reports, len(purls))
		for _, report := range reports {
			require.Contains(purls, report.Coordinates)
			require.NotNil(report.Description)
			require.NotNil(report.Reference)
			require.NotEmpty(report.Vulnerabilities)
		}
	}
}
