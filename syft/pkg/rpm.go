package pkg

import (
	"sort"
	"strings"

	"github.com/scylladb/go-set/strset"

	"github.com/anchore/syft/syft/file"
)

// RpmDBGlob is the glob pattern used to find RPM DB files. Where:
// - /var/lib/rpm/... is the typical path for most distributions
// - /usr/share/rpm/... is common for rpm-ostree distributions (coreos-like)
// - Packages is the legacy Berkeley db based format
// - Packages.db is the "ndb" format used in SUSE
// - rpmdb.sqlite is the sqlite format used in fedora + derivates
const RpmDBGlob = "**/{var/lib,usr/share,usr/lib/sysimage}/rpm/{Packages,Packages.db,rpmdb.sqlite}"

// RpmManifestGlob is used in CBL-Mariner distroless images
const RpmManifestGlob = "**/var/lib/rpmmanifest/container-manifest-2"

var _ FileOwner = (*RpmDBEntry)(nil)

// RpmArchive represents all captured data from a RPM package archive.
type RpmArchive RpmDBEntry

// RpmDBEntry represents all captured data from a RPM DB package entry.
type RpmDBEntry struct {
	Name            string          `json:"name"`
	Version         string          `json:"version"`
	Epoch           *int            `json:"epoch"  cyclonedx:"epoch" jsonschema:"nullable"`
	Arch            string          `json:"architecture"`
	Release         string          `json:"release" cyclonedx:"release"`
	SourceRpm       string          `json:"sourceRpm" cyclonedx:"sourceRpm"`
	Signatures      []RpmSignature  `json:"signatures,omitempty" cyclonedx:"signatures"`
	Size            int             `json:"size" cyclonedx:"size"`
	Vendor          string          `json:"vendor"`
	ModularityLabel *string         `json:"modularityLabel,omitempty"`
	Provides        []string        `json:"provides,omitempty"`
	Requires        []string        `json:"requires,omitempty"`
	Files           []RpmFileRecord `json:"files"`
}

type RpmSignature struct {
	PublicKeyAlgorithm string `json:"algo"`
	HashAlgorithm      string `json:"hash"`
	Created            string `json:"created"`
	IssuerKeyID        string `json:"issuer"`
}

func (s RpmSignature) String() string {
	if s.PublicKeyAlgorithm == "" && s.HashAlgorithm == "" && s.Created == "" && s.IssuerKeyID == "" {
		return ""
	}
	// mimics the output you would see from rpm -q --qf "%{RSAHEADER}"
	// e.g."RSA/SHA256, Mon May 16 12:32:55 2022, Key ID 702d426d350d275d"
	return strings.Join([]string{s.PublicKeyAlgorithm + "/" + s.HashAlgorithm, s.Created, "Key ID " + s.IssuerKeyID}, ", ")
}

// RpmFileRecord represents the file metadata for a single file attributed to a RPM package.
type RpmFileRecord struct {
	Path      string      `json:"path"`
	Mode      RpmFileMode `json:"mode"`
	Size      int         `json:"size"`
	Digest    file.Digest `json:"digest"`
	UserName  string      `json:"userName"`
	GroupName string      `json:"groupName"`
	Flags     string      `json:"flags"`
}

// RpmFileMode is the raw file mode for a single file. This can be interpreted as the linux stat.h mode (see https://pubs.opengroup.org/onlinepubs/007908799/xsh/sysstat.h.html)
type RpmFileMode uint16

func (m RpmDBEntry) OwnedFiles() (result []string) {
	s := strset.New()
	for _, f := range m.Files {
		if f.Path != "" {
			s.Add(f.Path)
		}
	}
	result = s.List()
	sort.Strings(result)
	return result
}
