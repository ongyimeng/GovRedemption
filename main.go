package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// StaffMapping represents a record mapping a staff pass ID to a team name.
type StaffMapping struct {
	StaffPassID string
	TeamName    string
	CreatedAt   int64
}

// LoadMappingFromCSV reads a CSV file at filePath and returns a slice of StaffMapping.
// The CSV file must have headers: "staff_pass_id", "team_name", "created_at" (epoch milliseconds).
func LoadMappingFromCSV(filePath string) ([]StaffMapping, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %w", err)
	}

	var mappings []StaffMapping
	// Skip header row and iterate over each record.
	for i, record := range records {
		if i == 0 {
			// Assuming first row is the header.
			continue
		}
		if len(record) < 3 {
			// Skip invalid records.
			continue
		}
		createdAt, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			// Skip record if created_at is invalid.
			continue
		}
		mapping := StaffMapping{
			StaffPassID: record[0],
			TeamName:    record[1],
			CreatedAt:   createdAt,
		}
		mappings = append(mappings, mapping)
	}

	return mappings, nil
}

// BuildLookupMap converts a slice of StaffMapping into a map for quick lookup by staff pass ID.
func BuildLookupMap(mappings []StaffMapping) map[string]string {
	lookup := make(map[string]string)
	for _, mapping := range mappings {
		lookup[mapping.StaffPassID] = mapping.TeamName
	}
	return lookup
}

// Redemption represents a redemption record with team name and the timestamp when redemption occurred.
type Redemption struct {
	TeamName   string
	RedeemedAt int64
}

// RedemptionManager manages redemption records and ensures a team can redeem only once.
type RedemptionManager struct {
	mu          sync.Mutex
	redemptions map[string]Redemption
}

// NewRedemptionManager initializes a new RedemptionManager.
func NewRedemptionManager() *RedemptionManager {
	return &RedemptionManager{
		redemptions: make(map[string]Redemption),
	}
}

// IsEligible returns true if the team has not redeemed their gift yet.
func (rm *RedemptionManager) IsEligible(teamName string) bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	_, exists := rm.redemptions[teamName]
	return !exists
}

// AddRedemption adds a redemption record for the team if they are eligible.
// If the team has already redeemed, it returns an error.
func (rm *RedemptionManager) AddRedemption(teamName string) (*Redemption, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.redemptions[teamName]; exists {
		return nil, errors.New("team has already redeemed their gift")
	}

	redemption := Redemption{
		TeamName:   teamName,
		RedeemedAt: time.Now().UnixMilli(), // Requires Go 1.17 or later.
	}
	rm.redemptions[teamName] = redemption
	return &redemption, nil
}

func main() {
	// For demonstration, the program expects one command-line argument:
	// 1. Path to the mapping CSV file.
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <mapping_csv_file>")
		return
	}

	mappingFile := os.Args[1]

	// Load mappings from the CSV file.
	mappings, err := LoadMappingFromCSV(mappingFile)
	if err != nil {
		fmt.Println("Error loading mapping file:", err)
		return
	}
	lookup := BuildLookupMap(mappings)

	// Initialize RedemptionManager.
	redemptionManager := NewRedemptionManager()

	for {
		// Prompt the user for a staff pass ID.
		fmt.Print("Enter staff pass ID (or type 'exit' to quit): ")
		var staffPassID string
		fmt.Scanln(&staffPassID)

		// Check if the user wants to exit.
		if staffPassID == "exit" {
			fmt.Println("Exiting the program.")
			break
		}

		// Lookup the team name for the given staff pass ID.
		teamName, found := lookup[staffPassID]
		if !found {
			fmt.Println("Staff pass ID not found.")
			continue
		}
		fmt.Println("Staff pass belongs to team:", teamName)

		// Check if the team is eligible for redemption.
		if redemptionManager.IsEligible(teamName) {
			redemption, err := redemptionManager.AddRedemption(teamName)
			if err != nil {
				fmt.Println("Error during redemption:", err)
			} else {
				fmt.Printf("Redemption successful for team %s at timestamp %d\n", redemption.TeamName, redemption.RedeemedAt)
			}
		} else {
			fmt.Println("Team has already redeemed their gift. Please send the representative away.")
		}
	}
}
