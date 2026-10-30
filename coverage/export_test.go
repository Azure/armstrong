package coverage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/magodo/azure-rest-api-index/azidx"
)

const testIndexRelease = "index-2026-07-27"

var pinnedTestIndex *azidx.Index

// SetupPinnedIndex configures the pinned online snapshot for a non-parallel test.
// It downloads the index once and restores the previous globals after each test.
func SetupPinnedIndex(t *testing.T) string {
	t.Helper()

	previousIndexURL, previousRepoURL, previousCache := indexFileURL, azureRepoURL, indexCache
	t.Cleanup(func() {
		indexFileURL, azureRepoURL, indexCache = previousIndexURL, previousRepoURL, previousCache
	})

	indexFileURL = "https://github.com/teowa/azure-rest-api-index-file/releases/download/" + testIndexRelease + "/index.json.zip"
	indexCache = pinnedTestIndex
	index, err := GetIndex("")
	if err != nil {
		t.Fatalf("load pinned test index: %v", err)
	}
	if !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(index.Commit) {
		t.Fatalf("test index must record a full commit SHA, got %q", index.Commit)
	}
	pinnedTestIndex = index
	azureRepoURL = "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/" + index.Commit + "/specification/"
	return azureRepoURL
}

// SetupPinnedIndexCacheFile writes the pinned index to a temporary cache file for a non-parallel test.
// It clears the in-memory index and disables downloads, so index lookups must load the returned file.
func SetupPinnedIndexCacheFile(t *testing.T) (cacheFile, repoURL string) {
	t.Helper()

	repoURL = SetupPinnedIndex(t)
	data, err := json.Marshal(indexCache)
	if err != nil {
		t.Fatalf("marshal pinned test index: %v", err)
	}
	cacheFile = filepath.Join(t.TempDir(), "index.json")
	if err := os.WriteFile(cacheFile, data, 0o600); err != nil {
		t.Fatalf("write pinned test index cache: %v", err)
	}

	indexCache = nil
	indexFileURL = "http://127.0.0.1:0/index-download-disabled"
	return cacheFile, repoURL
}

func TestSetupPinnedIndex(t *testing.T) {
	previousIndexURL, previousRepoURL, previousCache := indexFileURL, azureRepoURL, indexCache
	t.Cleanup(func() {
		indexFileURL, azureRepoURL, indexCache = previousIndexURL, previousRepoURL, previousCache
	})
	staleIndex := &azidx.Index{Commit: "stale"}
	indexCache = staleIndex

	t.Run("pinned snapshot", func(t *testing.T) {
		repoURL := SetupPinnedIndex(t)
		index, err := GetIndex("")
		if err != nil {
			t.Fatal(err)
		}
		if index == staleIndex {
			t.Fatal("pinned setup reused the previous index cache")
		}
		expectedURL := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/" + index.Commit + "/specification/"
		if repoURL != expectedURL || azureRepoURL != expectedURL {
			t.Fatalf("spec URL %q does not match index commit %q", azureRepoURL, index.Commit)
		}
	})

	if indexFileURL != previousIndexURL || azureRepoURL != previousRepoURL || indexCache != staleIndex {
		t.Fatal("pinned setup did not restore the previous index configuration")
	}
}
