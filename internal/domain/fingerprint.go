package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
)

func NormalizeText(value string) string {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	return strings.Join(strings.Fields(trimmed), " ")
}

func NormalizePath(value string) string {
	normalized := strings.ReplaceAll(strings.TrimSpace(value), "\\", "/")
	normalized = filepath.ToSlash(normalized)
	normalized = strings.TrimPrefix(normalized, "./")
	return NormalizeText(normalized)
}

func StableResourceKey(resourceGUID, resourcePath string) string {
	guid := NormalizeText(resourceGUID)
	if guid != "" {
		return "guid:" + guid
	}

	path := NormalizePath(resourcePath)
	if path != "" {
		return "path:" + path
	}

	return ""
}

func GenerateIssueFingerprint(projectCode, ruleCode, resourceKey, locationKey, currentValue, expectedValue, message string) string {
	source := strings.Join([]string{
		NormalizeText(projectCode),
		NormalizeText(ruleCode),
		NormalizeText(resourceKey),
		NormalizeText(locationKey),
		NormalizeText(currentValue),
		NormalizeText(expectedValue),
		NormalizeText(message),
	}, "|")

	sum := sha256.Sum256([]byte(source))
	return hex.EncodeToString(sum[:])
}
