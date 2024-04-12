// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the LGPLv3, see LICENCE file for details.

package config_test

import "github.com/juju/packaging/v2"

var (
	// testedPackages is a slice of random package tests to run tests on.
	testedPackages = []string{
		"awesome-wm",
		"archey3",
		"arch-chroot",
		"ranger",
	}

	testedSource = packaging.PackageSource{
		Name: "Some Totally Official Source.",
		URL:  "some-source.com/packages",
		Key:  "some-key",
	}

	testedPrefs = packaging.PackagePreferences{
		Path:        "/etc/my-package-manager.d/prefs_file.conf",
		Explanation: "don't judge me",
		Package:     "some-package",
		Pin:         "releases/extra-special",
		Priority:    42,
	}
)
