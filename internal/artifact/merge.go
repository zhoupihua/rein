package artifact

import (
	"fmt"
	"regexp"
	"strings"
)

// MergeResult reports the outcome of merging a DeltaPlan into a base spec.
type MergeResult struct {
	Added    int
	Modified int
	Removed  int
	Renamed  int
}

// MergeDeltas applies a DeltaPlan to a base spec and returns the updated content.
// Operations are applied in strict order: RENAMED → REMOVED → MODIFIED → ADDED.
func MergeDeltas(baseContent string, plan *DeltaPlan) (string, *MergeResult, error) {
	result := &MergeResult{}

	// Extract the requirements section
	before, headerLine, preamble, blocks, after := ExtractRequirementsSection(baseContent)

	// Build name-to-block map
	blockMap := make(map[string]RequirementBlock)
	var blockOrder []string
	for _, b := range blocks {
		normalized := normalizeReqName(b.Name)
		blockMap[normalized] = b
		blockOrder = append(blockOrder, normalized)
	}

	// Validate cross-section conflicts
	if err := validateDeltaConflicts(plan, blockMap); err != nil {
		return "", nil, err
	}

	// Apply RENAMED
	for _, pair := range plan.Renamed {
		fromNorm := normalizeReqName(pair.From)
		toNorm := normalizeReqName(pair.To)

		block, ok := blockMap[fromNorm]
		if !ok {
			return "", nil, fmt.Errorf("RENAMED: requirement %q not found in base spec", pair.From)
		}
		if _, exists := blockMap[toNorm]; exists {
			return "", nil, fmt.Errorf("RENAMED: target name %q already exists in base spec", pair.To)
		}

		// Replace header line with new name
		newHeader := reqHeaderRe.ReplaceAllString(block.HeaderLine, "### Requirement: "+pair.To)
		block.HeaderLine = newHeader
		block.Name = pair.To
		block.Raw = strings.Replace(block.Raw, blockMap[fromNorm].HeaderLine, newHeader, 1)

		delete(blockMap, fromNorm)
		blockMap[toNorm] = block

		// Update blockOrder
		for i, name := range blockOrder {
			if name == fromNorm {
				blockOrder[i] = toNorm
				break
			}
		}

		result.Renamed++
	}

	// Apply REMOVED
	for _, name := range plan.Removed {
		norm := normalizeReqName(name)
		if _, ok := blockMap[norm]; !ok {
			// Not an error for new specs (no base), but warn
			continue
		}
		delete(blockMap, norm)
		result.Removed++
	}

	// Apply MODIFIED
	for _, block := range plan.Modified {
		norm := normalizeReqName(block.Name)
		if _, ok := blockMap[norm]; !ok {
			return "", nil, fmt.Errorf("MODIFIED: requirement %q not found in base spec", block.Name)
		}
		blockMap[norm] = block
		result.Modified++
	}

	// Apply ADDED
	for _, block := range plan.Added {
		norm := normalizeReqName(block.Name)
		if _, ok := blockMap[norm]; ok {
			return "", nil, fmt.Errorf("ADDED: requirement %q already exists in base spec", block.Name)
		}
		blockMap[norm] = block
		blockOrder = append(blockOrder, norm)
		result.Added++
	}

	// Reassemble the spec
	updated := reassembleSpec(before, headerLine, preamble, blockOrder, blockMap, after)
	return updated, result, nil
}

// validateDeltaConflicts checks for cross-section conflicts.
func validateDeltaConflicts(plan *DeltaPlan, blockMap map[string]RequirementBlock) error {
	addedNames := make(map[string]bool)
	for _, b := range plan.Added {
		n := normalizeReqName(b.Name)
		addedNames[n] = true
	}

	modifiedNames := make(map[string]bool)
	for _, b := range plan.Modified {
		n := normalizeReqName(b.Name)
		modifiedNames[n] = true
	}

	removedNames := make(map[string]bool)
	for _, name := range plan.Removed {
		removedNames[normalizeReqName(name)] = true
	}

	renamedFrom := make(map[string]bool)
	renamedTo := make(map[string]bool)
	for _, pair := range plan.Renamed {
		renamedFrom[normalizeReqName(pair.From)] = true
		renamedTo[normalizeReqName(pair.To)] = true
	}

	// MODIFIED and REMOVED conflict
	for name := range modifiedNames {
		if removedNames[name] {
			return fmt.Errorf("requirement %q appears in both MODIFIED and REMOVED", name)
		}
	}

	// MODIFIED and ADDED conflict
	for name := range modifiedNames {
		if addedNames[name] {
			return fmt.Errorf("requirement %q appears in both MODIFIED and ADDED", name)
		}
	}

	// ADDED and REMOVED conflict
	for name := range addedNames {
		if removedNames[name] {
			return fmt.Errorf("requirement %q appears in both ADDED and REMOVED", name)
		}
	}

	// MODIFIED must reference the NEW name after a rename
	for name := range modifiedNames {
		if renamedFrom[name] {
			return fmt.Errorf("MODIFIED references %q which is a RENAMED FROM — use the new name instead", name)
		}
	}

	// RENAMED TO cannot collide with ADDED
	for name := range renamedTo {
		if addedNames[name] {
			return fmt.Errorf("RENAMED TO %q collides with an ADDED requirement", name)
		}
	}

	return nil
}

// reassembleSpec rebuilds the spec content from parsed components.
func reassembleSpec(before, headerLine, preamble string, blockOrder []string, blockMap map[string]RequirementBlock, after string) string {
	var parts []string

	if before != "" {
		parts = append(parts, before)
	}

	if headerLine != "" {
		parts = append(parts, headerLine)
	}

	if preamble != "" {
		parts = append(parts, preamble)
	}

	seen := make(map[string]bool)
	for _, name := range blockOrder {
		block, ok := blockMap[name]
		if !ok || seen[name] {
			continue
		}
		seen[name] = true
		parts = append(parts, block.Raw)
	}

	// Append any blocks not in original order (newly ADDED)
	for name, block := range blockMap {
		if !seen[name] {
			parts = append(parts, block.Raw)
		}
	}

	if after != "" {
		parts = append(parts, after)
	}

	result := strings.Join(parts, "\n\n")

	// Collapse 3+ consecutive newlines to 2
	collapseRe := regexp.MustCompile(`\n{3,}`)
	result = collapseRe.ReplaceAllString(result, "\n\n")

	return result
}

// normalizeReqName normalizes a requirement name for comparison.
func normalizeReqName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}
