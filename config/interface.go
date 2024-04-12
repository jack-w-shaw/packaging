// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the LGPLv3, see LICENCE file for details.

// The config package defines an interface which returns packaging-related
// configuration options and operations depending on the desired
// package-management system.
package config

import (
	"github.com/juju/packaging/v2"
)

// PackagingConfigurer is an interface which handles various packaging-related configuration
// functions for the specific distribution it represents.
type PackagingConfigurer interface {
	// DefaultPackages returns a list of default packages which should be
	// installed the vast majority of cases on any specific machine
	DefaultPackages() []string

	// IsCloudArchivePackage signals whether the given package is a
	// cloud archive package and thus should be set as such.
	IsCloudArchivePackage(pack string) bool

	// ApplyCloudArchiveTarget returns the package with the required target
	// release bits preceding it.
	ApplyCloudArchiveTarget(pack string) []string

	// RenderSource returns the os-specific full file contents
	// of a given PackageSource.
	RenderSource(src packaging.PackageSource) (string, error)

	// RenderPreferences returns the os-specific full file contents of a given
	// set of PackagePreferences.
	RenderPreferences(prefs packaging.PackagePreferences) (string, error)
}

func NewPackagingConfigurer(os string) (PackagingConfigurer, error) {
	switch os {
	case "centos":
		return NewYumPackagingConfigurer(), nil
	case "opensuse":
		return NewZypperPackagingConfigurer(), nil
	default:
		return NewAptPackagingConfigurer(), nil
	}
}

// NewAptPackagingConfigurer returns a PackagingConfigurer for apt-based systems.
func NewAptPackagingConfigurer() PackagingConfigurer {
	return &aptConfigurer{&baseConfigurer{
		defaultPackages:      UbuntuDefaultPackages,
		cloudArchivePackages: cloudArchivePackagesUbuntu,
	}}
}

// NewYumPackagingConfigurer returns a PackagingConfigurer for yum-based systems.
func NewYumPackagingConfigurer() PackagingConfigurer {
	return &yumConfigurer{&baseConfigurer{
		defaultPackages:      CentOSDefaultPackages,
		cloudArchivePackages: cloudArchivePackagesCentOS,
	}}
}

func NewZypperPackagingConfigurer() PackagingConfigurer {
	return &zypperConfigurer{&baseConfigurer{
		defaultPackages:      OpenSUSEDefaultPackages,
		cloudArchivePackages: cloudArchivePackagesOpenSUSE,
	}}
}
