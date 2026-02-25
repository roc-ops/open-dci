//go:build !js

package mibresolver

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sleepinggenius2/gosmi"
)

// New creates a resolver that loads MIBs from the given root directory.
// mibsDir is the root mibs/ folder. The resolver loads from:
// ietf/, iana/, cablelabs/*/, vendors/*/
//
// gosmi natively finds files named MODULE.mib in its search paths, so the
// symlinks (IF-MIB.mib -> IF-MIB@2000-06-14.mib) are resolved automatically.
// Versioned files (IF-MIB@2000-06-14.mib) are ignored because the "@" prefix
// prevents module name matching. For version overrides, a temporary directory
// is created with MODULE.mib symlinks pointing to the requested version.
//
// Correct versions of foundational MIBs (SNMPv2-SMI, SNMPv2-TC, SNMPv2-CONF,
// RFC1155-SMI) are embedded in this package and written to a temporary
// directory that is searched before the repository's mibs/ tree.
func New(mibsDir string, opts ...Option) (*Resolver, error) {
	cfg := &config{
		versionOverrides: make(map[string]string),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Resolve mibsDir to absolute path.
	absDir, err := filepath.Abs(mibsDir)
	if err != nil {
		return nil, fmt.Errorf("resolving mibs directory: %w", err)
	}

	// Collect all MIB source directories.
	searchDirs := collectSearchDirs(absDir)
	if len(searchDirs) == 0 {
		return nil, fmt.Errorf("no MIB directories found in %s", absDir)
	}

	// Create a temporary directory for core MIB overrides and version overrides.
	tempDir, err := os.MkdirTemp("", "mibresolver-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp directory: %w", err)
	}

	// Write embedded core MIBs to the temp directory so they take precedence
	// over potentially broken copies in the repository.
	if err := writeCoreMIBs(tempDir); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("writing core MIBs: %w", err)
	}

	// If there are version overrides, create symlinks in the temp directory
	// mapping MODULE.mib -> the specific versioned file.
	for moduleName, versionSpec := range cfg.versionOverrides {
		versionedFile := versionSpec + ".mib"
		found := false
		for _, dir := range searchDirs {
			srcPath := filepath.Join(dir, versionedFile)
			if _, statErr := os.Stat(srcPath); statErr == nil {
				realPath, evalErr := filepath.EvalSymlinks(srcPath)
				if evalErr != nil {
					log.Printf("mibresolver: override %s: cannot resolve path: %v", versionSpec, evalErr)
					continue
				}
				linkPath := filepath.Join(tempDir, moduleName+".mib")
				// Remove any existing file (could be a core MIB).
				os.Remove(linkPath)
				if symlinkErr := os.Symlink(realPath, linkPath); symlinkErr != nil {
					log.Printf("mibresolver: override %s: cannot create symlink: %v", versionSpec, symlinkErr)
					continue
				}
				found = true
				break
			}
		}
		if !found {
			log.Printf("mibresolver: override %s: versioned file not found", versionSpec)
		}
	}

	// Initialize gosmi.
	gosmi.Init()

	// Add temp dir first so core MIBs and overrides take precedence.
	gosmi.AppendPath(tempDir)

	// Add all MIB source directories to gosmi's search path.
	for _, dir := range searchDirs {
		gosmi.AppendPath(dir)
	}

	// Discover all module names from symlinks (MODULE.mib, no "@" in name)
	// in the search directories and temp directory.
	moduleNames := discoverModuleNames(append([]string{tempDir}, searchDirs...))

	var loaded int
	for _, moduleName := range moduleNames {
		// gosmi prints parse errors via fmt.Println; suppress by
		// temporarily redirecting stdout to /dev/null.
		origStdout := os.Stdout
		devNull, _ := os.Open(os.DevNull)
		if devNull != nil {
			os.Stdout = devNull
		}
		_, loadErr := gosmi.LoadModule(moduleName)
		os.Stdout = origStdout
		if devNull != nil {
			devNull.Close()
		}
		if loadErr != nil {
			log.Printf("mibresolver: failed to load %s: %v", moduleName, loadErr)
			continue
		}
		loaded++
	}

	if loaded == 0 {
		gosmi.Exit()
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("no MIB modules loaded successfully")
	}

	return &Resolver{
		tempDir: tempDir,
	}, nil
}

// writeCoreMIBs extracts the embedded core MIB files to the given directory.
func writeCoreMIBs(dir string) error {
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
		outPath := filepath.Join(dir, d.Name())
		if err := os.WriteFile(outPath, data, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", outPath, err)
		}
		return nil
	})
}

// collectSearchDirs returns the list of directories under mibsDir that may
// contain MIB files: ietf/, iana/, cablelabs/*/, vendors/*/
func collectSearchDirs(mibsDir string) []string {
	var dirs []string

	// Direct subdirs: ietf, iana
	for _, sub := range []string{"ietf", "iana"} {
		d := filepath.Join(mibsDir, sub)
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			dirs = append(dirs, d)
		}
	}

	// Nested subdirs: cablelabs/*, vendors/*
	for _, sub := range []string{"cablelabs", "vendors"} {
		parent := filepath.Join(mibsDir, sub)
		entries, err := os.ReadDir(parent)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				dirs = append(dirs, filepath.Join(parent, entry.Name()))
			}
		}
	}

	return dirs
}

// discoverModuleNames scans directories for MIB files (MODULE.mib, no "@")
// and returns a deduplicated list of module names to load.
func discoverModuleNames(searchDirs []string) []string {
	seen := make(map[string]bool)
	var names []string

	for _, dir := range searchDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			name := entry.Name()
			if !strings.HasSuffix(name, ".mib") {
				continue
			}
			// Skip versioned files (contain "@").
			if strings.Contains(name, "@") {
				continue
			}
			moduleName := strings.TrimSuffix(name, ".mib")
			if !seen[moduleName] {
				seen[moduleName] = true
				names = append(names, moduleName)
			}
		}
	}

	return names
}
