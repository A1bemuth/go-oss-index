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

func TestProcessesMissingVersion(t *testing.T) {
	client := Client{}
	require := require.New(t)

	_, err := client.Get([]string{"pkg:npm/url-parse"})
	require.Equal(err, ErrMissingCoordinatesVersion)
}

func TestProcessesTooManyRequests(t *testing.T) {
	client := Client{}

	for i := 0; i < 100; i++ {
		_, err := client.Get([]string{"pkg:npm/url-parse@1.4.2"})
		if err == ErrTooManyRequests {
			return
		}
	}
	t.Fatal("Did not get a TooManyRequests error")
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
