package main

import "golang.org/x/mod/modfile"

// ParseGoMod reads and parses a go.mod file
func ParseGoMod(filename string, data []byte) (*modfile.File, error) {
	return modfile.Parse(filename, data, nil)
}

// BuildVersionMap creates a map of module paths to versions from a modfile
func BuildVersionMap(mod *modfile.File) VersionMap {
	versions := make(VersionMap)
	for _, req := range mod.Require {
		versions[req.Mod.Path] = req.Mod.Version
	}
	return versions
}
