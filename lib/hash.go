package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

/*
Dirhash walks into a provided directory and calculates its SHA256 checksum,
based on the checksum and name of th individual files within the directory.
A list of exlcudedPaths as glob patterns can be provided to make Dirhash ignore their matches
*/
func DirHash(path string, ignoredPaths []string) string {
	// Keep existing behavior for backward compatibility
	if !filepath.IsAbs(path) {
		baseDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("Using relative path ", baseDir)
		path = filepath.Join(baseDir, path)
	}
	var allFiles, err = walkDir(path)
	if err != nil {
		log.Fatal(err)
	}
	exlcudedFilesMatch, err := filesToIgnore(path, ignoredPaths)
	if err != nil {
		log.Fatal(err)
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		for i := 0; i < len(exlcudedFilesMatch); i++ {
			log.Debug("excluding: ", exlcudedFilesMatch[i])
		}
	}
	var filesToHash = []string{}
	for i := 0; i < len(allFiles); i++ {
		if !slices.Contains(exlcudedFilesMatch, allFiles[i]) {
			filesToHash = append(filesToHash, allFiles[i])
		}
	}
	var fileHashes = []string{}
	for i := 0; i < len(filesToHash); i++ {
		fileHash, err := fileSha256(filesToHash[i])
		if err != nil {
			log.Fatal(fmt.Sprintf("Error hashing file %s : %s", filesToHash[i], err))
		}
		relpath, _ := filepath.Rel(path, filesToHash[i])
		hashCombo := fmt.Sprintf("%s %s", relpath, fileHash)
		fileHashes = append(fileHashes, hashCombo)
		log.Debug("hashing: ", hashCombo)
	}
	return mergeAllHashes(fileHashes)
}

// FileHash represents a per-file digest entry used for JSON output.
type FileHash struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// DirHashDetails computes the directory hash and also returns per-file hashes.
// It mirrors DirHash behavior but provides structured details for consumers.
func DirHashDetails(path string, ignoredPaths []string) (string, []FileHash) {
	if !filepath.IsAbs(path) {
		baseDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("Using relative path ", baseDir)
		path = filepath.Join(baseDir, path)
	}

	allFiles, err := walkDir(path)
	if err != nil {
		log.Fatal(err)
	}
	exlcudedFilesMatch, err := filesToIgnore(path, ignoredPaths)
	if err != nil {
		log.Fatal(err)
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		for i := 0; i < len(exlcudedFilesMatch); i++ {
			log.Debug("excluding: ", exlcudedFilesMatch[i])
		}
	}

	var filesToHash = []string{}
	for i := 0; i < len(allFiles); i++ {
		if !slices.Contains(exlcudedFilesMatch, allFiles[i]) {
			filesToHash = append(filesToHash, allFiles[i])
		}
	}

	var pairs []FileHash
	var combined = []string{}
	for i := 0; i < len(filesToHash); i++ {
		h, err := fileSha256(filesToHash[i])
		if err != nil {
			log.Fatal(fmt.Sprintf("Error hashing file %s : %s", filesToHash[i], err))
		}
		relpath, _ := filepath.Rel(path, filesToHash[i])
		pairs = append(pairs, FileHash{Path: relpath, Hash: h})
		combined = append(combined, fmt.Sprintf("%s %s", relpath, h))
		log.Debug("hashing: ", relpath, " ", h)
	}

	overall := mergeAllHashes(combined)
	return overall, pairs
}

// mergeAllHashes returns hash of joint slice elements as lines
func mergeAllHashes(hashes []string) string {
	sort.Strings(hashes)
	hash := strings.Join(hashes, "\n")
	return stringSha256(strings.NewReader(hash))
}

// fileSha256 returns the hash of a file
func fileSha256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	return stringSha256(f), nil
}

// stringSha256 returns the SHA256 checksum of a given io.Reader
func stringSha256(f io.Reader) string {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error(err)
	}
	return hex.EncodeToString(h.Sum(nil)[:])
}
