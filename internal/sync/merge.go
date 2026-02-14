package sync

import (
	"fmt"
	"strings"
)

// MergeStrategy controls how non-equal blocks are handled during merge
type MergeStrategy int

const (
	// MergeUnion includes blocks from both sides (additions from either
	// side are kept, modifications result in both versions appearing)
	MergeUnion MergeStrategy = iota
	// MergePreferLocal keeps local blocks and discards remote-only blocks
	MergePreferLocal
	// MergePreferRemote keeps remote blocks and discards local-only blocks
	MergePreferRemote
)

// SyncResult holds the outcome of a sync operation
type SyncResult struct {
	Merged  string   // the merged markdown body
	Changes []DiffOp // the block-level diff operations
	Summary string   // human-readable change summary
}

// Sync compares two markdown document bodies at the block level, merges
// them according to the given strategy, and returns the result.
func Sync(local, remote string, strategy MergeStrategy) SyncResult {
	localBlocks := Parse(local)
	remoteBlocks := Parse(remote)
	ops := Diff(localBlocks, remoteBlocks)
	merged := Merge(ops, strategy)

	return SyncResult{
		Merged:  Render(merged),
		Changes: ops,
		Summary: summarize(ops),
	}
}

// Merge applies a merge strategy to a set of diff operations and returns
// the resulting block sequence.
func Merge(ops []DiffOp, strategy MergeStrategy) []Block {
	var result []Block

	for _, op := range ops {
		switch op.Type {
		case OpEqual:
			result = append(result, *op.Local)

		case OpLocalOnly:
			switch strategy {
			case MergeUnion, MergePreferLocal:
				result = append(result, *op.Local)
			case MergePreferRemote:
				// discard local-only blocks
			}

		case OpRemoteOnly:
			switch strategy {
			case MergeUnion, MergePreferRemote:
				result = append(result, *op.Remote)
			case MergePreferLocal:
				// discard remote-only blocks
			}
		}
	}

	return result
}

// HasChanges returns true if the diff contains any non-equal operations
func HasChanges(ops []DiffOp) bool {
	for _, op := range ops {
		if op.Type != OpEqual {
			return true
		}
	}
	return false
}

// summarize produces a human-readable summary of the diff operations
func summarize(ops []DiffOp) string {
	var localOnly, remoteOnly, equal int
	for _, op := range ops {
		switch op.Type {
		case OpEqual:
			equal++
		case OpLocalOnly:
			localOnly++
		case OpRemoteOnly:
			remoteOnly++
		}
	}

	if localOnly == 0 && remoteOnly == 0 {
		return "no changes detected"
	}

	var parts []string
	if equal > 0 {
		parts = append(parts, fmt.Sprintf("%d %s unchanged", equal, pluralBlock(equal)))
	}
	if localOnly > 0 {
		parts = append(parts, fmt.Sprintf("%d %s only in local", localOnly, pluralBlock(localOnly)))
	}
	if remoteOnly > 0 {
		parts = append(parts, fmt.Sprintf("%d %s only in remote", remoteOnly, pluralBlock(remoteOnly)))
	}

	return strings.Join(parts, ", ")
}

func pluralBlock(n int) string {
	if n == 1 {
		return "block"
	}
	return "blocks"
}
