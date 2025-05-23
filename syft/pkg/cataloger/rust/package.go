package rust

import (
	"github.com/rust-secure-code/go-rustaudit"

	"github.com/anchore/packageurl-go"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/pkg"
)

// Pkg returns the standard `pkg.Package` representation of the package referenced within the Cargo.lock metadata.
func newPackageFromCargoMetadata(m pkg.RustCargoLockEntry, locations ...file.Location) pkg.Package {
	p := pkg.Package{
		Name:      m.Name,
		Version:   m.Version,
		Locations: file.NewLocationSet(locations...),
		PURL:      packageURL(m.Name, m.Version),
		Language:  pkg.Rust,
		Type:      pkg.RustPkg,
		Metadata:  m,
	}

	p.SetID()

	return p
}

func newPackageFromAudit(dep *rustaudit.Package, locations ...file.Location) pkg.Package {
	p := pkg.Package{
		Name:      dep.Name,
		Version:   dep.Version,
		PURL:      packageURL(dep.Name, dep.Version),
		Language:  pkg.Rust,
		Type:      pkg.RustPkg,
		Locations: file.NewLocationSet(locations...),
		FoundBy:   cargoAuditBinaryCatalogerName,
		Metadata: pkg.RustBinaryAuditEntry{
			Name:    dep.Name,
			Version: dep.Version,
			Source:  dep.Source,
		},
	}

	p.SetID()

	return p
}

// packageURL returns the PURL for the specific rust package (see https://github.com/package-url/purl-spec)
func packageURL(name, version string) string {
	return packageurl.NewPackageURL(
		packageurl.TypeCargo,
		"",
		name,
		version,
		nil,
		"",
	).ToString()
}
