package common

import (
	"bytes"
	"fmt"
	"hash/crc32"
)

// Generates a hash for the set hash function used by the IDs

func DataResourceIdsHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", HashString(buf.String()))
}

// HashString hashes a string to a unique hashcode.
//
// This will be removed in v2 without replacement. So we place here instead of import.
func HashString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
