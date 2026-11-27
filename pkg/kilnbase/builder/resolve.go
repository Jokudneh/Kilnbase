package builder

import (
	"fmt"
	"strings"
)

// ValidateImageRef rejects a reference that carries neither a tag nor a digest.
func ValidateImageRef(ref string) error {
	if ref == "" || (!strings.Contains(ref, ":") && !strings.Contains(ref, "@sha256:")) {
		return fmt.Errorf("image reference %q has no tag and no digest", ref)
	}
	return nil
}
