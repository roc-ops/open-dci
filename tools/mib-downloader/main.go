// Command mib-downloader fetches SNMP MIB files from IETF (RFCs), IANA, and
// CableLabs, renames them to a versioned naming convention, creates symlinks
// for the latest version, and generates a manifest.json.
//
// Usage (run from project root):
//
//	go run ./tools/mib-downloader/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Data types
// ---------------------------------------------------------------------------

// MIBEntry represents a single MIB module discovered during download.
type MIBEntry struct {
	Name     string `json:"name"`
	Revision string `json:"revision"`
	Source   string `json:"source"`
	Path     string `json:"path"`
	Latest   bool   `json:"latest"`
}

// Manifest is the top-level structure written to mibs/manifest.json.
type Manifest struct {
	GeneratedAt string                  `json:"generated_at"`
	Sources     map[string]SourceInfo   `json:"sources"`
	MIBs        []MIBEntry              `json:"mibs"`
}

// SourceInfo records per-source summary data.
type SourceInfo struct {
	Count  int    `json:"count"`
	Source string `json:"source"`
}

// ---------------------------------------------------------------------------
// Hardcoded RFC list (from snmp-mibs-downloader rfclist)
// ---------------------------------------------------------------------------

var ietfRFCs = []int{
	1155, 1212, 1213, 1215, 1227, 1381, 1382, 1414, 1461, 1471,
	1472, 1473, 1474, 1493, 1512, 1513, 1525, 1559, 1567, 1573,
	1592, 1593, 1611, 1612, 1628, 1657, 1658, 1659, 1660, 1666,
	1694, 1696, 1697, 1724, 1742, 1747, 1748, 1749, 1792, 2006,
	2011, 2012, 2013, 2020, 2024, 2051, 2096, 2108, 2115, 2127,
	2128, 2206, 2213, 2214, 2232, 2233, 2238, 2266, 2287, 2320,
	2417, 2452, 2454, 2455, 2456, 2457, 2494, 2512, 2513, 2514,
	2515, 2561, 2562, 2564, 2570, 2571, 2572, 2573, 2574, 2575,
	2576, 2578, 2579, 2580, 2584, 2594, 2605, 2613, 2618, 2619,
	2620, 2621, 2662, 2665, 2666, 2667, 2668, 2669, 2670, 2671,
	2672, 2674, 2677, 2707, 2720, 2737, 2742, 2786, 2787, 2788,
	2789, 2790, 2819, 2837, 2856, 2863, 2864, 2922, 2925, 2932,
	2933, 2934, 2940, 2954, 2955, 2959, 2981, 2982, 3014, 3019,
	3020, 3159, 3164, 3165, 3201, 3202, 3231, 3273, 3276, 3289,
	3291, 3371, 3395, 3411, 3412, 3413, 3414, 3415, 3417, 3418,
	3419, 3433, 3498, 3512, 3555, 3559, 3560, 3561, 3562, 3566,
	3567, 3584, 3591, 3592, 3593, 3606, 3621, 3635, 3637, 3728,
	3729, 3805, 3806, 3813, 3815, 3826, 3873, 3877, 3895, 3896,
	4022, 4044, 4087, 4088, 4089, 4113, 4131, 4133, 4188, 4220,
	4268, 4292, 4293, 4318, 4319, 4323, 4363, 4368, 4369, 4502,
	4546, 4560, 4631, 4668, 4669, 4670, 4671, 4672, 4710, 4747,
	4801, 4802, 4803, 4837, 4878, 4898, 4935, 5017, 5060, 5066,
	5097, 5131, 5132, 5190, 5240, 5249, 5324, 5427, 5428, 5519,
	5525, 5580, 5601, 5602, 5603, 5604, 5605, 5676, 5813, 5815,
	6011, 6065, 6241, 6340, 6353, 6368, 6445, 6461, 6474, 6493,
	6643, 6655, 6660, 6727, 6768, 6786, 6840, 6850, 6933, 6978,
	7020, 7069, 7076, 7079, 7098, 7109, 7124, 7131, 7132, 7133,
	7166, 7167, 7186, 7209, 7257, 7311, 7330, 7420, 7461, 7495,
	7559, 7577, 7580, 7630, 7666, 8247, 8502, 8503,
}

// IANA MIB names (the URL path component used at iana.org).
var ianaMIBs = []string{
	"ianaiftype-mib",
	"ianalanguage-mib",
	"ianaaddressfamilynumbers-mib",
	"ianaiprouteprotocol-mib",
	"ianatn3270etc-mib",
	"ianamalloc-mib",
	"ianacharset-mib",
	"ianaprinter-mib",
	"ianafinisher-mib",
	"ianaitualarmtc-mib",
	"ianagmplstc-mib",
	"ianaippmmetricsregistry-mib",
	"ianamau-mib",
	"ianaentity-mib",
	"ianapowerstateset-mib",
	"ianastoragemediatype-mib",
}

// CableLabs categories to crawl.
var cableLabsCategories = []string{
	"DOCSIS",
	"common",
	"OpenCable",
	"wireless",
}

// ---------------------------------------------------------------------------
// Regex patterns
// ---------------------------------------------------------------------------

var (
	// Matches a MIB module header: NAME DEFINITIONS ::= BEGIN
	reDefBegin = regexp.MustCompile(`(?m)^[ \t]*(([A-Z][A-Za-z0-9-]*)\s+DEFINITIONS\s*::=\s*BEGIN)\s*$`)
	// Matches a standalone END statement.
	reEnd = regexp.MustCompile(`(?m)^[ \t]*END\s*$`)
	// Matches LAST-UPDATED value.
	reLastUpdated = regexp.MustCompile(`LAST-UPDATED\s+"([0-9]+Z?)"`)
	// Matches href in Apache directory listing.
	reHref = regexp.MustCompile(`href="([^"]+)"`)
	// Matches CableLabs archive filename with date: NAME-YYYY-MM-DD.txt
	reCLArchiveDate = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})\.txt$`)
	// Page break pattern: form feed character (0x0C).
	reFormFeed = regexp.MustCompile(`\x0c`)
)

// ---------------------------------------------------------------------------
// HTTP client with timeout
// ---------------------------------------------------------------------------

var httpClient = &http.Client{
	Timeout: 60 * time.Second,
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	log.SetFlags(0)

	// Determine project root (tool expects to run from project root).
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot determine working directory: %v", err)
	}

	mibsDir := filepath.Join(root, "mibs")

	// Ensure base directories exist.
	for _, d := range []string{
		filepath.Join(mibsDir, "ietf"),
		filepath.Join(mibsDir, "iana"),
		filepath.Join(mibsDir, "cablelabs", "DOCSIS"),
		filepath.Join(mibsDir, "cablelabs", "common"),
		filepath.Join(mibsDir, "cablelabs", "OpenCable"),
		filepath.Join(mibsDir, "cablelabs", "wireless"),
		filepath.Join(mibsDir, "vendors"),
	} {
		os.MkdirAll(d, 0o755)
	}

	var allEntries []MIBEntry

	// ---- IETF ----
	ietfEntries := downloadIETF(filepath.Join(mibsDir, "ietf"))
	allEntries = append(allEntries, ietfEntries...)

	// ---- IANA ----
	ianaEntries := downloadIANA(filepath.Join(mibsDir, "iana"))
	allEntries = append(allEntries, ianaEntries...)

	// ---- CableLabs ----
	clEntries := downloadCableLabs(filepath.Join(mibsDir, "cablelabs"))
	allEntries = append(allEntries, clEntries...)

	// Deduplicate entries (same Path means same file).
	allEntries = deduplicateEntries(allEntries)

	// Create latest-version symlinks and mark entries.
	allEntries = createSymlinks(allEntries, mibsDir)

	// Generate manifest.
	ietfCount := countBySource(allEntries, "ietf")
	ianaCount := countBySource(allEntries, "iana")
	clCount := countBySource(allEntries, "cablelabs")

	manifest := Manifest{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Sources: map[string]SourceInfo{
			"ietf":      {Count: ietfCount, Source: "https://www.rfc-editor.org/rfc/"},
			"iana":      {Count: ianaCount, Source: "https://www.iana.org/assignments/"},
			"cablelabs": {Count: clCount, Source: "https://mibs.cablelabs.com/MIBs/"},
		},
		MIBs: allEntries,
	}

	manifestPath := filepath.Join(mibsDir, "manifest.json")
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		log.Fatalf("error marshaling manifest: %v", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(manifestPath, data, 0o644); err != nil {
		log.Fatalf("error writing manifest: %v", err)
	}

	log.Printf("Manifest written to %s", manifestPath)
	log.Printf("Summary: %d IETF, %d IANA, %d CableLabs MIBs (%d total)",
		ietfCount, ianaCount, clCount, len(allEntries))
}

// ---------------------------------------------------------------------------
// IETF: download RFCs and extract embedded MIB modules
// ---------------------------------------------------------------------------

func downloadIETF(destDir string) []MIBEntry {
	rfcs := uniqueInts(ietfRFCs)
	log.Printf("Downloading IETF MIBs from %d RFCs...", len(rfcs))

	var entries []MIBEntry
	for i, rfc := range rfcs {
		url := fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", rfc)
		body, err := httpGet(url)
		if err != nil {
			log.Printf("  WARNING: rfc%d: %v", rfc, err)
			rateSleep()
			continue
		}

		extracted := extractMIBsFromRFC(string(body), rfc)
		for _, em := range extracted {
			revision := parseLastUpdated(em.body)
			if revision == "" {
				revision = "unknown"
			}

			fname := fmt.Sprintf("%s@%s.mib", em.name, revision)
			fpath := filepath.Join(destDir, fname)

			// Idempotent: skip if exists.
			if fileExists(fpath) {
				log.Printf("  [%d/%d] rfc%d: %s (already exists, skipped)", i+1, len(rfcs), rfc, em.name)
			} else {
				if err := os.WriteFile(fpath, []byte(em.body), 0o644); err != nil {
					log.Printf("  WARNING: writing %s: %v", fname, err)
					continue
				}
				log.Printf("  [%d/%d] rfc%d: extracted %s -> %s", i+1, len(rfcs), rfc, em.name, fname)
			}

			entries = append(entries, MIBEntry{
				Name:     em.name,
				Revision: revision,
				Source:   "ietf",
				Path:     filepath.Join("ietf", fname),
			})
		}

		rateSleep()
	}

	return entries
}

// extractedMIB holds a single MIB module extracted from an RFC.
type extractedMIB struct {
	name string
	body string
}

// extractMIBsFromRFC scans RFC text for embedded MIB modules.
func extractMIBsFromRFC(text string, rfcNum int) []extractedMIB {
	// First, strip \r.
	text = strings.ReplaceAll(text, "\r", "")

	// Strip page breaks: blank line, footer, FF, header, blank line.
	// The form feed is 0x0C. The pattern around it varies, but generally:
	//   \n<optional whitespace>\n<footer line>\n\x0c<optional whitespace>\n<header line>\n<optional whitespace>\n
	// We use a broad but safe pattern.
	text = stripPageBreaks(text)

	var results []extractedMIB

	// Find all DEFINITIONS ::= BEGIN locations.
	locs := reDefBegin.FindAllStringSubmatchIndex(text, -1)
	for _, loc := range locs {
		// loc[4]:loc[5] is the MIB name (submatch 2).
		mibName := text[loc[4]:loc[5]]
		startIdx := loc[0]

		// Find the matching END after this BEGIN.
		remainder := text[loc[1]:]
		endLoc := reEnd.FindStringIndex(remainder)
		if endLoc == nil {
			log.Printf("  WARNING: rfc%d: no END found for %s", rfcNum, mibName)
			continue
		}

		endIdx := loc[1] + endLoc[1]
		mibBody := text[startIdx:endIdx]

		results = append(results, extractedMIB{
			name: mibName,
			body: mibBody,
		})
	}

	return results
}

// stripPageBreaks removes RFC page break artifacts from text.
// Pattern: lines around a form feed character.
func stripPageBreaks(text string) string {
	if !reFormFeed.MatchString(text) {
		return text
	}

	lines := strings.Split(text, "\n")
	var out []string
	i := 0
	for i < len(lines) {
		if strings.Contains(lines[i], "\x0c") {
			// Found a form feed line. Remove surrounding context:
			// Typically: blank line before, footer line before FF, the FF line,
			// header line after FF, blank line after.
			//
			// Work backwards to remove up to 3 lines before (blank + footer + possible blank).
			removeBack := 0
			for removeBack < 3 && len(out) > 0 && (strings.TrimSpace(out[len(out)-1]) == "" || removeBack < 2) {
				out = out[:len(out)-1]
				removeBack++
			}

			// Skip the FF line itself.
			i++

			// Skip up to 3 lines after (header + blank + possible blank).
			skipped := 0
			for skipped < 3 && i < len(lines) && (strings.TrimSpace(lines[i]) == "" || skipped < 2) {
				i++
				skipped++
			}
			continue
		}
		out = append(out, lines[i])
		i++
	}
	return strings.Join(out, "\n")
}

// ---------------------------------------------------------------------------
// IANA: download pre-formatted MIB files
// ---------------------------------------------------------------------------

func downloadIANA(destDir string) []MIBEntry {
	log.Printf("Downloading %d IANA MIBs...", len(ianaMIBs))

	var entries []MIBEntry
	for i, name := range ianaMIBs {
		url := fmt.Sprintf("https://www.iana.org/assignments/%s/%s", name, name)
		body, err := httpGet(url)
		if err != nil {
			log.Printf("  WARNING: %s: %v", name, err)
			rateSleep()
			continue
		}

		content := string(body)

		// Extract the actual MIB module name from the content.
		mibName := extractMIBName(content)
		if mibName == "" {
			// Fall back to uppercased filename.
			mibName = strings.ToUpper(name)
		}

		revision := parseLastUpdated(content)
		if revision == "" {
			revision = "unknown"
		}

		fname := fmt.Sprintf("%s@%s.mib", mibName, revision)
		fpath := filepath.Join(destDir, fname)

		if fileExists(fpath) {
			log.Printf("  [%d/%d] %s (already exists, skipped)", i+1, len(ianaMIBs), mibName)
		} else {
			if err := os.WriteFile(fpath, body, 0o644); err != nil {
				log.Printf("  WARNING: writing %s: %v", fname, err)
				rateSleep()
				continue
			}
			log.Printf("  [%d/%d] %s -> %s", i+1, len(ianaMIBs), name, fname)
		}

		entries = append(entries, MIBEntry{
			Name:     mibName,
			Revision: revision,
			Source:   "iana",
			Path:     filepath.Join("iana", fname),
		})

		rateSleep()
	}

	return entries
}

// ---------------------------------------------------------------------------
// CableLabs: crawl directory listings, download current + archive files
// ---------------------------------------------------------------------------

func downloadCableLabs(destDir string) []MIBEntry {
	log.Printf("Downloading CableLabs MIBs from %d categories...", len(cableLabsCategories))

	var entries []MIBEntry
	baseURL := "https://mibs.cablelabs.com/MIBs"

	for _, cat := range cableLabsCategories {
		catDir := filepath.Join(destDir, cat)
		catURL := fmt.Sprintf("%s/%s/", baseURL, cat)

		// Download current .mib files from the category page.
		currentFiles := crawlDirectoryListing(catURL, ".mib")
		log.Printf("  Category %s: %d current files", cat, len(currentFiles))

		for _, fname := range currentFiles {
			fileURL := catURL + fname
			body, err := httpGet(fileURL)
			if err != nil {
				log.Printf("    WARNING: %s: %v", fname, err)
				rateSleep()
				continue
			}

			content := string(body)
			mibName := extractMIBName(content)
			if mibName == "" {
				mibName = strings.TrimSuffix(fname, ".mib")
			}

			revision := parseLastUpdated(content)
			if revision == "" {
				revision = "unknown"
			}

			outFname := fmt.Sprintf("%s@%s.mib", mibName, revision)
			outPath := filepath.Join(catDir, outFname)

			if fileExists(outPath) {
				log.Printf("    %s (already exists, skipped)", outFname)
			} else {
				if err := os.WriteFile(outPath, body, 0o644); err != nil {
					log.Printf("    WARNING: writing %s: %v", outFname, err)
					rateSleep()
					continue
				}
				log.Printf("    %s -> %s", fname, outFname)
			}

			relPath := filepath.Join("cablelabs", cat, outFname)
			entries = append(entries, MIBEntry{
				Name:     mibName,
				Revision: revision,
				Source:   "cablelabs",
				Path:     relPath,
			})

			rateSleep()
		}

		// Download archive files.
		archiveURL := catURL + "Archive-PreviousVersions/"
		archiveFiles := crawlDirectoryListing(archiveURL, ".txt")
		log.Printf("  Category %s: %d archive files", cat, len(archiveFiles))

		for _, fname := range archiveFiles {
			fileURL := archiveURL + fname
			body, err := httpGet(fileURL)
			if err != nil {
				log.Printf("    WARNING: archive %s: %v", fname, err)
				rateSleep()
				continue
			}

			content := string(body)
			mibName := extractMIBName(content)
			if mibName == "" {
				// Use filename stem, stripping date suffix.
				stem := strings.TrimSuffix(fname, ".txt")
				if m := reCLArchiveDate.FindStringSubmatch(fname); m != nil {
					stem = strings.TrimSuffix(stem, "-"+m[1])
				}
				mibName = stem
			}

			revision := ""
			// Try to get date from filename first for archive files.
			if m := reCLArchiveDate.FindStringSubmatch(fname); m != nil {
				revision = m[1]
			}
			if revision == "" {
				revision = parseLastUpdated(content)
			}
			if revision == "" {
				revision = "unknown"
			}

			outFname := fmt.Sprintf("%s@%s.mib", mibName, revision)
			outPath := filepath.Join(catDir, outFname)

			if fileExists(outPath) {
				log.Printf("    %s (already exists, skipped)", outFname)
			} else {
				if err := os.WriteFile(outPath, body, 0o644); err != nil {
					log.Printf("    WARNING: writing %s: %v", outFname, err)
					rateSleep()
					continue
				}
				log.Printf("    archive %s -> %s", fname, outFname)
			}

			relPath := filepath.Join("cablelabs", cat, outFname)
			entries = append(entries, MIBEntry{
				Name:     mibName,
				Revision: revision,
				Source:   "cablelabs",
				Path:     relPath,
			})

			rateSleep()
		}
	}

	return entries
}

// crawlDirectoryListing fetches an Apache-style directory listing page and
// returns filenames matching the given suffix.
func crawlDirectoryListing(url, suffix string) []string {
	body, err := httpGet(url)
	if err != nil {
		log.Printf("  WARNING: listing %s: %v", url, err)
		return nil
	}

	html := string(body)
	matches := reHref.FindAllStringSubmatch(html, -1)
	var files []string
	seen := make(map[string]bool)
	for _, m := range matches {
		href := m[1]
		// Skip parent directory links, query strings, and absolute URLs.
		if strings.HasPrefix(href, "/") || strings.HasPrefix(href, "?") || strings.Contains(href, "://") {
			continue
		}
		// Skip directory links (end with /).
		if strings.HasSuffix(href, "/") {
			continue
		}
		if strings.HasSuffix(strings.ToLower(href), strings.ToLower(suffix)) {
			if !seen[href] {
				seen[href] = true
				files = append(files, href)
			}
		}
	}
	return files
}

// ---------------------------------------------------------------------------
// Symlink creation
// ---------------------------------------------------------------------------

// createSymlinks groups entries by (source, name), picks the latest revision,
// and creates symlinks. Returns updated entries with Latest flag set.
func createSymlinks(entries []MIBEntry, mibsDir string) []MIBEntry {
	log.Println("Creating latest-version symlinks...")

	// Group by source directory + MIB name.
	type key struct {
		dir  string // relative dir within mibs/
		name string
	}
	groups := make(map[key][]int) // key -> indices into entries

	for i, e := range entries {
		dir := filepath.Dir(e.Path)
		k := key{dir: dir, name: e.Name}
		groups[k] = append(groups[k], i)
	}

	for k, indices := range groups {
		// Sort by revision descending.
		sort.Slice(indices, func(a, b int) bool {
			return entries[indices[a]].Revision > entries[indices[b]].Revision
		})

		latestIdx := indices[0]
		entries[latestIdx].Latest = true

		// Create symlink: MIB-NAME.mib -> MIB-NAME@YYYY-MM-DD.mib
		linkName := k.name + ".mib"
		targetName := filepath.Base(entries[latestIdx].Path)
		linkPath := filepath.Join(mibsDir, k.dir, linkName)

		// Remove existing symlink if any.
		os.Remove(linkPath)

		if err := os.Symlink(targetName, linkPath); err != nil {
			log.Printf("  WARNING: symlink %s: %v", linkPath, err)
		}
	}

	return entries
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// httpGet fetches a URL and returns the body bytes.
func httpGet(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", url, err)
	}
	return body, nil
}

// extractMIBName looks for the first DEFINITIONS ::= BEGIN line and returns
// the MIB module name.
func extractMIBName(text string) string {
	m := reDefBegin.FindStringSubmatch(text)
	if m == nil {
		return ""
	}
	return m[2]
}

// parseLastUpdated extracts the LAST-UPDATED revision date from MIB text.
// Returns "YYYY-MM-DD" or empty string.
func parseLastUpdated(text string) string {
	m := reLastUpdated.FindStringSubmatch(text)
	if m == nil {
		return ""
	}
	val := m[1]

	// Try 12-digit format: YYYYMMDDHHmmZ
	t, err := time.Parse("200601021504Z", val)
	if err == nil {
		return t.Format("2006-01-02")
	}

	// Try 10-digit format: YYMMDDHHmmZ (old format)
	t, err = time.Parse("0601021504Z", val)
	if err == nil {
		return t.Format("2006-01-02")
	}

	// Try without Z suffix variations.
	cleaned := strings.TrimSuffix(val, "Z")
	if len(cleaned) >= 8 {
		year := cleaned[0:4]
		month := cleaned[4:6]
		day := cleaned[6:8]
		return year + "-" + month + "-" + day
	}
	if len(cleaned) >= 6 {
		// YYMMDDHHmm
		yy := cleaned[0:2]
		mm := cleaned[2:4]
		dd := cleaned[4:6]
		year := "19" + yy
		if yy < "70" {
			year = "20" + yy
		}
		return year + "-" + mm + "-" + dd
	}

	return ""
}

// rateSleep pauses briefly between HTTP requests.
func rateSleep() {
	time.Sleep(100 * time.Millisecond)
}

// fileExists checks whether a file exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// countBySource counts entries with the given source value.
func countBySource(entries []MIBEntry, source string) int {
	n := 0
	for _, e := range entries {
		if e.Source == source {
			n++
		}
	}
	return n
}

// deduplicateEntries removes duplicate entries with the same Path.
func deduplicateEntries(entries []MIBEntry) []MIBEntry {
	seen := make(map[string]bool)
	var result []MIBEntry
	for _, e := range entries {
		if !seen[e.Path] {
			seen[e.Path] = true
			result = append(result, e)
		}
	}
	return result
}

// uniqueInts removes duplicate integers from a slice, preserving order.
func uniqueInts(vals []int) []int {
	seen := make(map[int]bool)
	var result []int
	for _, v := range vals {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
