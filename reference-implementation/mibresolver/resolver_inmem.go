package mibresolver

import (
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/sleepinggenius2/gosmi"
)

// NewFromMIBData creates a resolver using in-memory MIB data instead of the
// OS filesystem. This is the primary constructor for WASM environments where
// filesystem access is unavailable.
//
// The embedded core MIBs (SNMPv2-SMI, SNMPv2-TC, SNMPv2-CONF, RFC1155-SMI)
// are always loaded first. If mibFiles is non-nil, those MIB files are loaded
// as well. mibFiles maps filename (e.g., "IF-MIB.mib") to file content.
func NewFromMIBData(mibFiles map[string][]byte) (*Resolver, error) {
	// Build in-memory FS with embedded core MIBs.
	coreFS := newMemFS()
	if err := loadEmbeddedCoreMIBs(coreFS); err != nil {
		return nil, fmt.Errorf("loading core MIBs: %w", err)
	}

	gosmi.Init()

	// Set core MIBs as the primary search path.
	gosmi.SetFS(gosmi.NamedFS("core", coreFS))

	// If user MIBs provided, add them as a separate FS.
	if len(mibFiles) > 0 {
		userFS := newMemFS()
		for name, data := range mibFiles {
			userFS.Add(name, data)
		}
		gosmi.AppendFS(gosmi.NamedFS("user", userFS))
	}

	// Discover and load all modules from the core FS.
	coreModules := discoverModulesFromFS(coreFS)
	var loaded int
	for _, moduleName := range coreModules {
		if _, err := gosmi.LoadModule(moduleName); err != nil {
			log.Printf("mibresolver: failed to load core %s: %v", moduleName, err)
			continue
		}
		loaded++
	}

	// Load user modules.
	if len(mibFiles) > 0 {
		for name := range mibFiles {
			if !strings.HasSuffix(name, ".mib") {
				continue
			}
			if strings.Contains(name, "@") {
				continue
			}
			moduleName := strings.TrimSuffix(name, ".mib")
			if _, err := gosmi.LoadModule(moduleName); err != nil {
				log.Printf("mibresolver: failed to load user %s: %v", moduleName, err)
				continue
			}
			loaded++
		}
	}

	if loaded == 0 {
		gosmi.Exit()
		return nil, fmt.Errorf("no MIB modules loaded successfully")
	}

	return &Resolver{}, nil
}

// LoadAdditionalMIBs adds more MIB files to an existing in-memory resolver.
// This supports incremental loading — users can upload MIBs one batch at a
// time. mibFiles maps filename to content.
func (r *Resolver) LoadAdditionalMIBs(mibFiles map[string][]byte) (int, error) {
	if len(mibFiles) == 0 {
		return 0, nil
	}

	userFS := newMemFS()
	for name, data := range mibFiles {
		userFS.Add(name, data)
	}
	gosmi.AppendFS(gosmi.NamedFS("user", userFS))

	var loaded int
	for name := range mibFiles {
		if !strings.HasSuffix(name, ".mib") {
			continue
		}
		if strings.Contains(name, "@") {
			continue
		}
		moduleName := strings.TrimSuffix(name, ".mib")
		if _, err := gosmi.LoadModule(moduleName); err != nil {
			log.Printf("mibresolver: failed to load %s: %v", moduleName, err)
			continue
		}
		loaded++
	}

	if loaded == 0 {
		return 0, fmt.Errorf("no MIB modules loaded from provided files")
	}

	return loaded, nil
}

// loadEmbeddedCoreMIBs reads the embedded coremib/*.mib files into a memFS.
func loadEmbeddedCoreMIBs(mfs *memFS) error {
	return fs.WalkDir(coreMIBs, "coremib", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := coreMIBs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading embedded %s: %w", path, err)
		}
		mfs.Add(d.Name(), data)
		return nil
	})
}

// discoverModulesFromFS lists module names from a memFS (files ending in .mib,
// no "@" in name).
func discoverModulesFromFS(mfs *memFS) []string {
	var names []string
	for name := range mfs.files {
		if !strings.HasSuffix(name, ".mib") {
			continue
		}
		if strings.Contains(name, "@") {
			continue
		}
		names = append(names, strings.TrimSuffix(name, ".mib"))
	}
	return names
}
