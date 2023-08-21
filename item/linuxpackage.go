package item

import (
	internalitem "github.com/cobbler/cobblerclient/internal/item"
	rawItem "github.com/cobbler/cobblerclient/internal/item/raw"
	resolvedItem "github.com/cobbler/cobblerclient/internal/item/resolved"
)

type Package struct {
	Item

	// Internal representation
	raw      rawItem.Package
	resolved resolvedItem.Package

	// Package specific attributes
	Installer internalitem.Property[string]
	Version   internalitem.Property[string]
}

type Installer struct {
	item Package
}

func (i Installer) Get() string {
	return i.item.raw.Installer
}

func (i Installer) Set(installer string) {
	i.item.raw.Installer = installer
}

type Version struct {
	item Package
}

func (v Version) Get() string {
	return v.item.raw.Version
}

func (v Version) Set(version string) {
	v.item.raw.Version = version
}
