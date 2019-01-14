package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtractLatestTagFromDockerHubResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HappyFlow",
			input:    `{"count": 2, "next": null, "previous": null, "results": [{"name": "v0.4.2", "full_size": 126289039, "images": [{"size": 126289039, "architecture": "amd64", "variant": null, "features": null, "os": "linux", "os_version": null, "os_features": null}], "id": 40717129, "repository": 6341803, "creator": 4691149, "last_updater": 4691149, "last_updated": "2018-11-26T14:02:32.560968Z", "image_id": null, "v2": true}, {"name": "v0.7.0", "full_size": 126289039, "images": [{"size": 126289039, "architecture": "amd64", "variant": null, "features": null, "os": "linux", "os_version": null, "os_features": null}], "id": 40686462, "repository": 6341803, "creator": 4691149, "last_updater": 4691149, "last_updated": "2018-11-26T14:02:27.625449Z", "image_id": null, "v2": true}]}`,
			expected: "v0.7.0",
		},
		{
			name:     "HappyFlowReversed",
			input:    `{"count": 2, "next": null, "previous": null, "results": [{"name": "v1.2.3", "full_size": 126289039, "images": [{"size": 126289039, "architecture": "amd64", "variant": null, "features": null, "os": "linux", "os_version": null, "os_features": null}], "id": 40717129, "repository": 6341803, "creator": 4691149, "last_updater": 4691149, "last_updated": "2018-11-26T14:02:32.560968Z", "image_id": null, "v2": true}, {"name": "v0.7.0", "full_size": 126289039, "images": [{"size": 126289039, "architecture": "amd64", "variant": null, "features": null, "os": "linux", "os_version": null, "os_features": null}], "id": 40686462, "repository": 6341803, "creator": 4691149, "last_updater": 4691149, "last_updated": "2018-11-26T14:02:27.625449Z", "image_id": null, "v2": true}]}`,
			expected: "v1.2.3",
		},
		{
			name:     "Empty",
			input:    ``,
			expected: "",
		},
		{
			name:     "NoResults",
			input:    `{"count": 0, "next": null, "previous": null, "results": []}`,
			expected: "",
		},
		{
			name:     "Corrupt",
			input:    `{"count": 2, "next": null, "previous": null, "results": [{"name": "v0.4.2", "full_size": 126289039, "ima`,
			expected: "",
		},
		{
			name:     "NonSemver",
			input:    `{"count": 2, "next": null, "previous": null, "results": [{"name": "latest", "full_size": 126289039, "images": [{"size": 126289039, "architecture": "amd64", "variant": null, "features": null, "os": "linux", "os_version": null, "os_features": null}], "id": 40717129, "repository": 6341803, "creator": 4691149, "last_updater": 4691149, "last_updated": "2018-11-26T14:02:32.560968Z", "image_id": null, "v2": true}, {"name": "v0.7.0", "full_size": 126289039, "images": [{"size": 126289039, "architecture": "amd64", "variant": null, "features": null, "os": "linux", "os_version": null, "os_features": null}], "id": 40686462, "repository": 6341803, "creator": 4691149, "last_updater": 4691149, "last_updated": "2018-11-26T14:02:27.625449Z", "image_id": null, "v2": true}]}`,
			expected: "v0.7.0",
		},
		{
			name:     "Experimental",
			input:    `{"count": 2, "next": null, "previous": null, "results": [{"name": "experimental", "full_size": 126289039, "images": [{"size": 126289039, "architecture": "amd64", "variant": null, "features": null, "os": "linux", "os_version": null, "os_features": null}], "id": 40717129, "repository": 6341803, "creator": 4691149, "last_updater": 4691149, "last_updated": "2018-11-26T14:02:32.560968Z", "image_id": null, "v2": true}, {"name": "v0.7.0", "full_size": 126289039, "images": [{"size": 126289039, "architecture": "amd64", "variant": null, "features": null, "os": "linux", "os_version": null, "os_features": null}], "id": 40686462, "repository": 6341803, "creator": 4691149, "last_updater": 4691149, "last_updated": "2018-11-26T14:02:27.625449Z", "image_id": null, "v2": true}]}`,
			expected: "v0.7.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag, err := extractLatestTagFromDockerHubResponse([]byte(tt.input))
			if tt.expected == "" {
				require.Error(t, err, "extract should return an error")
			} else {
				require.NoError(t, err, "extract should not return an error")
				require.Equal(t, tt.expected, tag, "extracted tag should match")
			}
		})
	}
}

func TestCmpTags(t *testing.T) {
	tests := []struct {
		name     string
		first    string
		second   string
		expected int
	}{
		{
			name:     "Equal",
			first:    "v1.2.3",
			second:   "v1.2.3",
			expected: 0, // latest, current -> no upgrade
		},
		{
			name:     "FirstNewer(patch)",
			first:    "v1.2.3",
			second:   "v1.2.0",
			expected: 1, // latest, current -> upgrade!
		},
		{
			name:     "FirstNewer(minor)",
			first:    "v1.2.3",
			second:   "v1.1.3",
			expected: 1, // latest, current -> upgrade!
		},
		{
			name:     "FirstNewer(major)",
			first:    "v2.1.3",
			second:   "v1.2.3",
			expected: 1, // latest, current -> upgrade!
		},
		{
			name:     "FirstOlder(patch)",
			first:    "v1.2.0",
			second:   "v1.2.3",
			expected: -1, // latest, current -> no upgrade
		},
		{
			name:     "FirstOlder(minor)",
			first:    "v1.1.3",
			second:   "v1.2.3",
			expected: -1, // latest, current -> no upgrade
		},
		{
			name:     "FirstOlder(major)",
			first:    "v1.2.3",
			second:   "v2.1.3",
			expected: -1, // latest, current -> no upgrade
		},
		{
			name:     "BothNoSemver",
			first:    "junk",
			second:   "junk",
			expected: -1, // latest, current -> no upgrade
		},
		{
			name:     "BothNoSemverDifferent",
			first:    "junk2",
			second:   "junk1",
			expected: -1, // latest, current -> no upgrade
		},
		{
			name:     "FirstNoSemver",
			first:    "junk",
			second:   "v1.2.3",
			expected: -1, // latest, current -> no upgrade
		},
		{
			name:     "SecondNoSemver",
			first:    "v1.2.3",
			second:   "junk",
			expected: 1, // latest, current -> upgrade!
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmpTags(tt.first, tt.second)
			switch tt.expected {
			case 1:
				require.True(t, result > 0, "result should be positive")
			case -1:
				require.True(t, result < 0, "result should be negative")
			case 0:
				require.Equal(t, 0, result, "result should be zero")
			}
		})
	}
}
