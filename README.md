# GovRedemption

Built a simple CLI system in Go. I have been wanting to build something in Go for a while now, and I thought that this was a simple project to start with. Took <2 hours to complete.

## Functionality

The GovRedemption CLI system allows users to manage redemption records for teams based on staff pass IDs. The main features include:

- **Load Mappings**: Load staff pass ID to team name mappings from a CSV file.
- **Check Eligibility**: Verify if a team is eligible for redemption based on their staff pass ID.
- **Add Redemption**: Record a redemption for a team if they have not already redeemed their gift.
- **Error Handling**: Inform users if a team has already redeemed their gift or if the staff pass ID is not found.

## Run Code 
To start the program, use the following command:

Start: go run main.go <mapping_csv_file>  
Example Start: go run main.go staff-id-to-team-mapping.csv

