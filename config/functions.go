// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the LGPLv3, see LICENCE file for details.

package config

// SeriesRequiresCloudArchiveTools signals whether the given series
// requires the configuration of cloud archive cloud tools.
func SeriesRequiresCloudArchiveTools(series string) bool {
	return seriesRequiringCloudTools[series]
}
