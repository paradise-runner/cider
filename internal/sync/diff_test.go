package sync

import (
	"testing"
)

func TestDiffIdentical(t *testing.T) {
	blocks := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "Paragraph one."},
		{Type: BlockParagraph, Content: "Paragraph two."},
	}

	ops := Diff(blocks, blocks)

	if len(ops) != 3 {
		t.Fatalf("Diff() returned %d ops, want 3", len(ops))
	}
	for i, op := range ops {
		if op.Type != OpEqual {
			t.Errorf("ops[%d].Type = %d, want OpEqual", i, op.Type)
		}
	}
}

func TestDiffLocalAddition(t *testing.T) {
	local := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "New paragraph."},
		{Type: BlockParagraph, Content: "Shared paragraph."},
	}
	remote := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "Shared paragraph."},
	}

	ops := Diff(local, remote)

	expected := []DiffOpType{OpEqual, OpLocalOnly, OpEqual}
	if len(ops) != len(expected) {
		t.Fatalf("Diff() returned %d ops, want %d", len(ops), len(expected))
	}
	for i, op := range ops {
		if op.Type != expected[i] {
			t.Errorf("ops[%d].Type = %d, want %d", i, op.Type, expected[i])
		}
	}
	if ops[1].Local.Content != "New paragraph." {
		t.Errorf("ops[1].Local.Content = %q, want %q", ops[1].Local.Content, "New paragraph.")
	}
}

func TestDiffRemoteAddition(t *testing.T) {
	local := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "Shared paragraph."},
	}
	remote := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "Remote paragraph."},
		{Type: BlockParagraph, Content: "Shared paragraph."},
	}

	ops := Diff(local, remote)

	expected := []DiffOpType{OpEqual, OpRemoteOnly, OpEqual}
	if len(ops) != len(expected) {
		t.Fatalf("Diff() returned %d ops, want %d", len(ops), len(expected))
	}
	for i, op := range ops {
		if op.Type != expected[i] {
			t.Errorf("ops[%d].Type = %d, want %d", i, op.Type, expected[i])
		}
	}
	if ops[1].Remote.Content != "Remote paragraph." {
		t.Errorf("ops[1].Remote.Content = %q, want %q", ops[1].Remote.Content, "Remote paragraph.")
	}
}

func TestDiffDeletion(t *testing.T) {
	local := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "Remaining."},
	}
	remote := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "Deleted paragraph."},
		{Type: BlockParagraph, Content: "Remaining."},
	}

	ops := Diff(local, remote)

	// The "Deleted paragraph." is remote-only
	expected := []DiffOpType{OpEqual, OpRemoteOnly, OpEqual}
	if len(ops) != len(expected) {
		t.Fatalf("Diff() returned %d ops, want %d", len(ops), len(expected))
	}
	for i, op := range ops {
		if op.Type != expected[i] {
			t.Errorf("ops[%d].Type = %d, want %d", i, op.Type, expected[i])
		}
	}
}

func TestDiffModification(t *testing.T) {
	local := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "Updated version of paragraph."},
		{Type: BlockParagraph, Content: "Shared end."},
	}
	remote := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "Original version of paragraph."},
		{Type: BlockParagraph, Content: "Shared end."},
	}

	ops := Diff(local, remote)

	// The modified paragraph appears as local-only + remote-only
	expected := []DiffOpType{OpEqual, OpLocalOnly, OpRemoteOnly, OpEqual}
	if len(ops) != len(expected) {
		t.Fatalf("Diff() returned %d ops, want %d", len(ops), len(expected))
	}
	for i, op := range ops {
		if op.Type != expected[i] {
			t.Errorf("ops[%d].Type = %d, want %d", i, op.Type, expected[i])
		}
	}
}

func TestDiffBothEmpty(t *testing.T) {
	ops := Diff(nil, nil)
	if len(ops) != 0 {
		t.Errorf("Diff(nil, nil) returned %d ops, want 0", len(ops))
	}
}

func TestDiffOneEmpty(t *testing.T) {
	blocks := []Block{
		{Type: BlockHeading, Content: "# Title"},
		{Type: BlockParagraph, Content: "Content."},
	}

	t.Run("local empty", func(t *testing.T) {
		ops := Diff(nil, blocks)
		if len(ops) != 2 {
			t.Fatalf("Diff() returned %d ops, want 2", len(ops))
		}
		for _, op := range ops {
			if op.Type != OpRemoteOnly {
				t.Errorf("expected OpRemoteOnly, got %d", op.Type)
			}
		}
	})

	t.Run("remote empty", func(t *testing.T) {
		ops := Diff(blocks, nil)
		if len(ops) != 2 {
			t.Fatalf("Diff() returned %d ops, want 2", len(ops))
		}
		for _, op := range ops {
			if op.Type != OpLocalOnly {
				t.Errorf("expected OpLocalOnly, got %d", op.Type)
			}
		}
	})
}

func TestDiffCompletelyDifferent(t *testing.T) {
	local := []Block{
		{Type: BlockParagraph, Content: "Local only A."},
		{Type: BlockParagraph, Content: "Local only B."},
	}
	remote := []Block{
		{Type: BlockParagraph, Content: "Remote only X."},
		{Type: BlockParagraph, Content: "Remote only Y."},
	}

	ops := Diff(local, remote)

	// All blocks are unique to their side
	if len(ops) != 4 {
		t.Fatalf("Diff() returned %d ops, want 4", len(ops))
	}

	localCount, remoteCount := 0, 0
	for _, op := range ops {
		switch op.Type {
		case OpLocalOnly:
			localCount++
		case OpRemoteOnly:
			remoteCount++
		case OpEqual:
			t.Error("unexpected OpEqual in completely different documents")
		}
	}
	if localCount != 2 {
		t.Errorf("got %d local-only ops, want 2", localCount)
	}
	if remoteCount != 2 {
		t.Errorf("got %d remote-only ops, want 2", remoteCount)
	}
}

func TestDiffReorderedBlocks(t *testing.T) {
	local := []Block{
		{Type: BlockParagraph, Content: "A"},
		{Type: BlockParagraph, Content: "B"},
		{Type: BlockParagraph, Content: "C"},
	}
	remote := []Block{
		{Type: BlockParagraph, Content: "C"},
		{Type: BlockParagraph, Content: "A"},
		{Type: BlockParagraph, Content: "B"},
	}

	ops := Diff(local, remote)

	// LCS is either {A, B} or {C} depending on algorithm.
	// With our LCS, A and B form a longer subsequence.
	equalCount := 0
	for _, op := range ops {
		if op.Type == OpEqual {
			equalCount++
		}
	}
	if equalCount < 2 {
		t.Errorf("expected at least 2 equal ops for reordered blocks, got %d", equalCount)
	}
}
