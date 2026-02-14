package sync

// DiffOpType describes what happened to a block between two document versions
type DiffOpType int

const (
	// OpEqual means the block exists in both versions with identical content
	OpEqual DiffOpType = iota
	// OpLocalOnly means the block exists only in the local version
	OpLocalOnly
	// OpRemoteOnly means the block exists only in the remote version
	OpRemoteOnly
)

// DiffOp represents a single operation in a block-level diff
type DiffOp struct {
	Type   DiffOpType
	Local  *Block // set for OpEqual and OpLocalOnly
	Remote *Block // set for OpEqual and OpRemoteOnly
}

// Diff computes a block-level diff between local and remote block sequences.
// It uses the Longest Common Subsequence (LCS) algorithm to align matching
// blocks, then produces a sequence of DiffOp values describing which blocks
// are shared and which are unique to each side.
func Diff(local, remote []Block) []DiffOp {
	lcs := computeLCS(local, remote)
	return buildDiffOps(local, remote, lcs)
}

// computeLCS returns the longest common subsequence of two block slices.
// Blocks are compared by content equality.
func computeLCS(a, b []Block) []Block {
	m, n := len(a), len(b)
	// Build DP table
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1].Content == b[j-1].Content {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}

	// Backtrack to find the actual subsequence
	length := dp[m][n]
	result := make([]Block, length)
	i, j := m, n
	k := length - 1
	for k >= 0 {
		if a[i-1].Content == b[j-1].Content {
			result[k] = a[i-1]
			i--
			j--
			k--
		} else if dp[i-1][j] >= dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	return result
}

// buildDiffOps walks local and remote block sequences aligned by their LCS
// and emits DiffOp values in document order.
func buildDiffOps(local, remote []Block, lcs []Block) []DiffOp {
	var ops []DiffOp
	li, ri, lci := 0, 0, 0

	for lci < len(lcs) {
		target := lcs[lci].Content

		// Emit local-only blocks before the next LCS match
		for li < len(local) && local[li].Content != target {
			b := local[li]
			ops = append(ops, DiffOp{Type: OpLocalOnly, Local: &b})
			li++
		}

		// Emit remote-only blocks before the next LCS match
		for ri < len(remote) && remote[ri].Content != target {
			b := remote[ri]
			ops = append(ops, DiffOp{Type: OpRemoteOnly, Remote: &b})
			ri++
		}

		// Emit the matched block
		lb := local[li]
		rb := remote[ri]
		ops = append(ops, DiffOp{Type: OpEqual, Local: &lb, Remote: &rb})
		li++
		ri++
		lci++
	}

	// Remaining local-only blocks after last LCS element
	for li < len(local) {
		b := local[li]
		ops = append(ops, DiffOp{Type: OpLocalOnly, Local: &b})
		li++
	}

	// Remaining remote-only blocks after last LCS element
	for ri < len(remote) {
		b := remote[ri]
		ops = append(ops, DiffOp{Type: OpRemoteOnly, Remote: &b})
		ri++
	}

	return ops
}
