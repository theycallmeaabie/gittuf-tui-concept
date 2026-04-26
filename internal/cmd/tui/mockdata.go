// Copyright The gittuf Authors
// SPDX-License-Identifier: Apache-2.0

package tui

// Hardcoded sample data for the demo screens. Lets the TUI run without a real
// gittuf repo so reviewers don't have to set one up just to look at the UX.

type mockKey struct {
	name        string
	keyType     string
	fingerprint string
	addedDate   string
	expiresDate string
	expired     bool
	isRoot      bool
}

type mockBranch struct {
	pattern           string
	threshold         int
	totalApprovers    int
	requiredApprovers []string
	forcePushBlocked  bool
}

type mockVerifyCase struct {
	refName       string
	commit        string
	passed        bool
	signers       []string
	threshold     int
	checks        []mockVerifyCheck
	failureReason string
}

type mockVerifyCheck struct {
	name   string
	passed bool
	detail string
}

var mockRootKey = mockKey{
	name:        "gittuf Root of Trust",
	keyType:     "ed25519",
	fingerprint: "SHA256:rT9pHkN8qLm2vYxZ4wRcF6jK1bE3sA7uG9oVnDhCpMw",
	addedDate:   "2024-01-15",
	expiresDate: "Never",
	expired:     false,
	isRoot:      true,
}

var mockTrustedKeys = []mockKey{
	{
		name:        "Aabid",
		keyType:     "ed25519",
		fingerprint: "SHA256:8fK2pNmQ7rT9vXcW3jL5bH1aE4yR6zG8oVdU0nBhCpM",
		addedDate:   "2024-03-22",
		expiresDate: "2027-04-26",
		expired:     false,
		isRoot:      false,
	},
	{
		name:        "Bobby",
		keyType:     "rsa-4096",
		fingerprint: "SHA256:Lm9oP2qN5tR7vY3wX8kJ4hG1bA6cE0sU2dF7iVnHpKr",
		addedDate:   "2024-07-08",
		expiresDate: "2026-08-15",
		expired:     false,
		isRoot:      false,
	},
	{
		name:        "Christine",
		keyType:     "ed25519",
		fingerprint: "SHA256:Qx4vR7yT2wK9pL6bM3jH8aE5sN1cF0oU6dG7iVrZpHm",
		addedDate:   "2023-11-03",
		expiresDate: "2026-01-01",
		expired:     true,
		isRoot:      false,
	},
}

var mockVerifyCases = []mockVerifyCase{
	{
		refName:   "refs/heads/main",
		commit:    "9a4f2c1e8b3d",
		passed:    true,
		signers:   []string{"Aabid", "Bobby"},
		threshold: 2,
		checks: []mockVerifyCheck{
			{name: "RSL entry exists", passed: true, detail: "found at offset 247"},
			{name: "Signatures valid", passed: true, detail: "2 of 2 signers verified"},
			{name: "Threshold met", passed: true, detail: "2-of-2 required for refs/heads/main"},
			{name: "Force push blocked", passed: true, detail: "ancestor relationship preserved"},
		},
	},
	{
		refName:   "refs/heads/release/v1.0",
		commit:    "3e7b1d9a6c2f",
		passed:    true,
		signers:   []string{"Aabid", "Bobby"},
		threshold: 2,
		checks: []mockVerifyCheck{
			{name: "RSL entry exists", passed: true, detail: "found at offset 312"},
			{name: "Signatures valid", passed: true, detail: "2 of 2 signers verified"},
			{name: "Threshold met", passed: true, detail: "2-of-2 required for release/*"},
			{name: "Force push blocked", passed: true, detail: "ancestor relationship preserved"},
		},
	},
	{
		refName:       "refs/heads/feature/quick-fix",
		commit:        "f1c8a2e5d3b7",
		passed:        false,
		signers:       []string{"Christine"},
		threshold:     2,
		failureReason: "Threshold not met: 1 valid signature, 2 required",
		checks: []mockVerifyCheck{
			{name: "RSL entry exists", passed: true, detail: "found at offset 401"},
			{name: "Signatures valid", passed: false, detail: "Christine's key expired 2026-01-01"},
			{name: "Threshold met", passed: false, detail: "0-of-2 valid signatures (1 expired)"},
			{name: "Force push blocked", passed: true, detail: "ancestor relationship preserved"},
		},
	},
}

var mockProtectedBranches = []mockBranch{
	{
		pattern:           "main",
		threshold:         2,
		totalApprovers:    3,
		requiredApprovers: []string{"Aabid", "Bobby", "Christine"},
		forcePushBlocked:  true,
	},
	{
		pattern:           "release/*",
		threshold:         2,
		totalApprovers:    2,
		requiredApprovers: []string{"Aabid", "Bobby"},
		forcePushBlocked:  true,
	},
	{
		pattern:           "hotfix/*",
		threshold:         1,
		totalApprovers:    2,
		requiredApprovers: []string{"Aabid", "Bobby"},
		forcePushBlocked:  false,
	},
}
