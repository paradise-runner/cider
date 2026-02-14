package sync

import (
	"strings"
	"testing"
)

func TestMergeUnionKeepsBothSides(t *testing.T) {
	ops := []DiffOp{
		{Type: OpEqual, Local: blockPtr(BlockHeading, "# Title"), Remote: blockPtr(BlockHeading, "# Title")},
		{Type: OpLocalOnly, Local: blockPtr(BlockParagraph, "Local addition.")},
		{Type: OpRemoteOnly, Remote: blockPtr(BlockParagraph, "Remote addition.")},
		{Type: OpEqual, Local: blockPtr(BlockParagraph, "Shared."), Remote: blockPtr(BlockParagraph, "Shared.")},
	}

	result := Merge(ops, MergeUnion)

	if len(result) != 4 {
		t.Fatalf("Merge(Union) returned %d blocks, want 4", len(result))
	}
	if result[0].Content != "# Title" {
		t.Errorf("block[0] = %q, want %q", result[0].Content, "# Title")
	}
	if result[1].Content != "Local addition." {
		t.Errorf("block[1] = %q, want %q", result[1].Content, "Local addition.")
	}
	if result[2].Content != "Remote addition." {
		t.Errorf("block[2] = %q, want %q", result[2].Content, "Remote addition.")
	}
	if result[3].Content != "Shared." {
		t.Errorf("block[3] = %q, want %q", result[3].Content, "Shared.")
	}
}

func TestMergePreferLocalDropsRemoteOnly(t *testing.T) {
	ops := []DiffOp{
		{Type: OpEqual, Local: blockPtr(BlockHeading, "# Title"), Remote: blockPtr(BlockHeading, "# Title")},
		{Type: OpLocalOnly, Local: blockPtr(BlockParagraph, "Local.")},
		{Type: OpRemoteOnly, Remote: blockPtr(BlockParagraph, "Remote.")},
		{Type: OpEqual, Local: blockPtr(BlockParagraph, "End."), Remote: blockPtr(BlockParagraph, "End.")},
	}

	result := Merge(ops, MergePreferLocal)

	if len(result) != 3 {
		t.Fatalf("Merge(PreferLocal) returned %d blocks, want 3", len(result))
	}
	for _, b := range result {
		if b.Content == "Remote." {
			t.Error("PreferLocal should not include remote-only blocks")
		}
	}
}

func TestMergePreferRemoteDropsLocalOnly(t *testing.T) {
	ops := []DiffOp{
		{Type: OpEqual, Local: blockPtr(BlockHeading, "# Title"), Remote: blockPtr(BlockHeading, "# Title")},
		{Type: OpLocalOnly, Local: blockPtr(BlockParagraph, "Local.")},
		{Type: OpRemoteOnly, Remote: blockPtr(BlockParagraph, "Remote.")},
		{Type: OpEqual, Local: blockPtr(BlockParagraph, "End."), Remote: blockPtr(BlockParagraph, "End.")},
	}

	result := Merge(ops, MergePreferRemote)

	if len(result) != 3 {
		t.Fatalf("Merge(PreferRemote) returned %d blocks, want 3", len(result))
	}
	for _, b := range result {
		if b.Content == "Local." {
			t.Error("PreferRemote should not include local-only blocks")
		}
	}
}

func TestMergeAllEqual(t *testing.T) {
	ops := []DiffOp{
		{Type: OpEqual, Local: blockPtr(BlockHeading, "# Title"), Remote: blockPtr(BlockHeading, "# Title")},
		{Type: OpEqual, Local: blockPtr(BlockParagraph, "Same."), Remote: blockPtr(BlockParagraph, "Same.")},
	}

	for _, strategy := range []MergeStrategy{MergeUnion, MergePreferLocal, MergePreferRemote} {
		result := Merge(ops, strategy)
		if len(result) != 2 {
			t.Errorf("strategy %d: got %d blocks, want 2", strategy, len(result))
		}
	}
}

func TestMergeEmpty(t *testing.T) {
	result := Merge(nil, MergeUnion)
	if len(result) != 0 {
		t.Errorf("Merge(nil) returned %d blocks, want 0", len(result))
	}
}

func TestHasChanges(t *testing.T) {
	t.Run("no changes", func(t *testing.T) {
		ops := []DiffOp{
			{Type: OpEqual, Local: blockPtr(BlockParagraph, "Same.")},
		}
		if HasChanges(ops) {
			t.Error("HasChanges() = true, want false")
		}
	})

	t.Run("with changes", func(t *testing.T) {
		ops := []DiffOp{
			{Type: OpEqual, Local: blockPtr(BlockParagraph, "Same.")},
			{Type: OpLocalOnly, Local: blockPtr(BlockParagraph, "New.")},
		}
		if !HasChanges(ops) {
			t.Error("HasChanges() = false, want true")
		}
	})

	t.Run("empty", func(t *testing.T) {
		if HasChanges(nil) {
			t.Error("HasChanges(nil) = true, want false")
		}
	})
}

func TestSyncEndToEnd(t *testing.T) {
	local := "# My Note\n\nLocal paragraph.\n\nShared paragraph.\n"
	remote := "# My Note\n\nShared paragraph.\n\nRemote paragraph.\n"

	result := Sync(local, remote, MergeUnion)

	if !strings.Contains(result.Merged, "# My Note") {
		t.Error("merged should contain heading")
	}
	if !strings.Contains(result.Merged, "Local paragraph.") {
		t.Error("merged should contain local paragraph")
	}
	if !strings.Contains(result.Merged, "Shared paragraph.") {
		t.Error("merged should contain shared paragraph")
	}
	if !strings.Contains(result.Merged, "Remote paragraph.") {
		t.Error("merged should contain remote paragraph")
	}
	if !HasChanges(result.Changes) {
		t.Error("should report changes")
	}
	if result.Summary == "" {
		t.Error("summary should not be empty")
	}
}

func TestSyncIdentical(t *testing.T) {
	doc := "# Title\n\nSame content.\n"

	result := Sync(doc, doc, MergeUnion)

	if HasChanges(result.Changes) {
		t.Error("identical documents should have no changes")
	}
	if result.Summary != "no changes detected" {
		t.Errorf("summary = %q, want %q", result.Summary, "no changes detected")
	}
}

func TestSyncPreferLocal(t *testing.T) {
	local := "# Title\n\nLocal version.\n\nShared.\n"
	remote := "# Title\n\nRemote version.\n\nShared.\n"

	result := Sync(local, remote, MergePreferLocal)

	if !strings.Contains(result.Merged, "Local version.") {
		t.Error("prefer-local should include local blocks")
	}
	if strings.Contains(result.Merged, "Remote version.") {
		t.Error("prefer-local should not include remote-only blocks")
	}
}

func TestSyncPreferRemote(t *testing.T) {
	local := "# Title\n\nLocal version.\n\nShared.\n"
	remote := "# Title\n\nRemote version.\n\nShared.\n"

	result := Sync(local, remote, MergePreferRemote)

	if !strings.Contains(result.Merged, "Remote version.") {
		t.Error("prefer-remote should include remote blocks")
	}
	if strings.Contains(result.Merged, "Local version.") {
		t.Error("prefer-remote should not include local-only blocks")
	}
}

func TestSyncStructuralDifferences(t *testing.T) {
	local := `# My Document

Introduction paragraph.

## Section One

- Item A
- Item B

## Section Two

Conclusion.
`
	remote := `# My Document

Introduction paragraph.

## Section One

- Item A
- Item B
- Item C

## New Section

New content here.

## Section Two

Conclusion.
`

	result := Sync(local, remote, MergeUnion)

	// Should include the expanded list from remote
	if !strings.Contains(result.Merged, "Item C") {
		t.Error("merged should include Item C from remote")
	}
	// Should include the new section from remote
	if !strings.Contains(result.Merged, "## New Section") {
		t.Error("merged should include new section from remote")
	}
	if !strings.Contains(result.Merged, "New content here.") {
		t.Error("merged should include new content from remote")
	}
	// Should keep original structure
	if !strings.Contains(result.Merged, "## Section One") {
		t.Error("merged should keep original sections")
	}
	if !strings.Contains(result.Merged, "Conclusion.") {
		t.Error("merged should keep conclusion")
	}
}

// blockPtr is a test helper that creates a pointer to a Block
func blockPtr(typ BlockType, content string) *Block {
	return &Block{Type: typ, Content: content}
}
