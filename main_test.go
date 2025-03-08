package main

import (
    "os"
    "strings"
    "testing"
    "time"
)

// TestLoadMappingFromCSV creates a temporary CSV file, loads the data, and verifies correct parsing.
func TestLoadMappingFromCSV(t *testing.T) {
    // Prepare sample CSV data.
    csvData := `staff_pass_id,team_name,created_at
STAFF_H123804820G,BASS,1623772799000
MANAGER_T999888420B,RUST,1623772799000
BOSS_T000000001P,RUST,1623872111000
`
    // Write CSV data to a temporary file.
    tmpFile, err := os.CreateTemp("", "test_mapping_*.csv")
    if err != nil {
        t.Fatalf("unable to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())
    if _, err := tmpFile.WriteString(csvData); err != nil {
        t.Fatalf("unable to write to temp file: %v", err)
    }
    tmpFile.Close()

    // Load mappings.
    mappings, err := LoadMappingFromCSV(tmpFile.Name())
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    // Expect 3 mappings.
    if len(mappings) != 3 {
        t.Errorf("expected 3 mappings, got %d", len(mappings))
    }

    // Check first mapping.
    if mappings[0].StaffPassID != "STAFF_H123804820G" || mappings[0].TeamName != "BASS" {
        t.Errorf("unexpected mapping record: %+v", mappings[0])
    }
}

// TestBuildLookupMap verifies that a lookup map is correctly built from mappings.
func TestBuildLookupMap(t *testing.T) {
    mappings := []StaffMapping{
        {StaffPassID: "ID1", TeamName: "TeamA", CreatedAt: 1000},
        {StaffPassID: "ID2", TeamName: "TeamB", CreatedAt: 2000},
    }
    lookup := BuildLookupMap(mappings)
    if team, ok := lookup["ID1"]; !ok || team != "TeamA" {
        t.Errorf("expected ID1 to map to TeamA, got %v", team)
    }
    if team, ok := lookup["ID2"]; !ok || team != "TeamB" {
        t.Errorf("expected ID2 to map to TeamB, got %v", team)
    }
}

// TestRedemptionManager tests eligibility checking and redemption addition.
func TestRedemptionManager(t *testing.T) {
    rm := NewRedemptionManager()
    teamName := "TeamA"

    // Initially, the team should be eligible.
    if !rm.IsEligible(teamName) {
        t.Errorf("expected team to be eligible for redemption")
    }

    // Add redemption and check the returned record.
    redemption, err := rm.AddRedemption(teamName)
    if err != nil {
        t.Errorf("expected redemption to succeed, got error: %v", err)
    }
    if redemption.TeamName != teamName {
        t.Errorf("expected team name %s, got %s", teamName, redemption.TeamName)
    }
    // Allow a little delay difference for timestamp check.
    if time.Now().UnixMilli()-redemption.RedeemedAt > 1000 {
        t.Errorf("redemption timestamp seems off: %d", redemption.RedeemedAt)
    }

    // Now the team should not be eligible.
    if rm.IsEligible(teamName) {
        t.Errorf("expected team to be ineligible after redemption")
    }

    // Attempting to add another redemption should fail.
    _, err = rm.AddRedemption(teamName)
    if err == nil || !strings.Contains(err.Error(), "already redeemed") {
        t.Errorf("expected error about team already redeemed, got %v", err)
    }
}
