package artifact

import (
	"path/filepath"
	"testing"
)

func TestParseSpecWithTestField(t *testing.T) {
	path := filepath.Join("testdata", "sample-spec.md")
	spec, err := ParseSpecFile(path)
	if err != nil {
		t.Fatalf("ParseSpecFile: %v", err)
	}

	// First scenario has TEST field
	regReq := findReq(t, spec.Requirements, "User Registration")
	if len(regReq.Scenarios) < 1 {
		t.Fatal("expected at least 1 scenario in User Registration")
	}
	s := regReq.Scenarios[0]
	if s.Name != "Successful registration" {
		t.Errorf("scenario name = %q, want %q", s.Name, "Successful registration")
	}
	if s.Test != "TestUserRegistration_SuccessfulRegistration" {
		t.Errorf("Test = %q, want %q", s.Test, "TestUserRegistration_SuccessfulRegistration")
	}
}

func TestParseSpecWithoutTestField(t *testing.T) {
	path := filepath.Join("testdata", "sample-spec.md")
	spec, err := ParseSpecFile(path)
	if err != nil {
		t.Fatalf("ParseSpecFile: %v", err)
	}

	// Second scenario (Duplicate email) has no TEST field
	regReq := findReq(t, spec.Requirements, "User Registration")
	if len(regReq.Scenarios) < 2 {
		t.Fatal("expected at least 2 scenarios in User Registration")
	}
	s := regReq.Scenarios[1]
	if s.Test != "" {
		t.Errorf("Test = %q, want empty", s.Test)
	}
}

func TestParseSpecTestFieldWithBackticks(t *testing.T) {
	content := `# Test Spec

### Requirement: Login

#### Scenario: Valid credentials

- **WHEN** POST /auth/login with valid email+password
- **THEN** returns 200 with JWT token
- **TEST** ` + "`" + `TestAuthJWT_ValidCredentials_ReturnsToken` + "`" + `

`
	spec := ParseSpecContent(content)
	if len(spec.Requirements) == 0 || len(spec.Requirements[0].Scenarios) == 0 {
		t.Fatal("expected at least 1 scenario")
	}
	s := spec.Requirements[0].Scenarios[0]
	if s.Test != "TestAuthJWT_ValidCredentials_ReturnsToken" {
		t.Errorf("Test = %q, want %q", s.Test, "TestAuthJWT_ValidCredentials_ReturnsToken")
	}
}

func TestParseSpecMixedScenarios(t *testing.T) {
	content := `# Test Spec

### Requirement: Mixed

#### Scenario: With test

- **WHEN** action
- **THEN** result
- **TEST** TestWithTest

#### Scenario: Without test

- **WHEN** action2
- **THEN** result2

`
	spec := ParseSpecContent(content)
	req := spec.Requirements[0]
	if len(req.Scenarios) != 2 {
		t.Fatalf("expected 2 scenarios, got %d", len(req.Scenarios))
	}
	if req.Scenarios[0].Test != "TestWithTest" {
		t.Errorf("scenario 0 Test = %q, want %q", req.Scenarios[0].Test, "TestWithTest")
	}
	if req.Scenarios[1].Test != "" {
		t.Errorf("scenario 1 Test = %q, want empty", req.Scenarios[1].Test)
	}
}

func findReq(t *testing.T, reqs []Requirement, name string) Requirement {
	t.Helper()
	for _, r := range reqs {
		if r.Name == name {
			return r
		}
	}
	t.Fatalf("requirement %q not found", name)
	return Requirement{}
}

func TestParseSpecOldFormatBackwardCompat(t *testing.T) {
	content := `# Old Format Spec

## Context

Some context here.

## Goals

- Goal one

## Non-Goals

- Not doing this

### Requirement: Login

#### Scenario: Valid

- **WHEN** credentials submitted
- **THEN** token returned

## Decisions

- **Decision:** Use JWT — **Rationale:** Stateless

## Risks / Trade-offs

- Token expiry
`
	spec := ParseSpecContent(content)
	// Old sections are tracked in Sections but not populated in fields
	if !spec.HasSection("Context") || !spec.HasSection("Goals") || !spec.HasSection("Non-Goals") {
		t.Error("old sections should be tracked in Sections list")
	}
	// Requirements, Decisions, Risks still work
	if len(spec.Requirements) != 1 {
		t.Errorf("expected 1 requirement, got %d", len(spec.Requirements))
	}
	if len(spec.Decisions) != 1 {
		t.Errorf("expected 1 decision, got %d", len(spec.Decisions))
	}
	if spec.Risks == "" {
		t.Error("Risks should be populated")
	}
}
