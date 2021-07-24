/*
 * Retrieves all FRC shirt trades on the yearly ChiefDelphi thread.
 * Aggregates data into CSV file with required metadata.
 * Uses Sheets API and JSON struct. See README for more details.
 */

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Struct to represent a Google Sheet of shirts
type shirtSheet struct {
	ID          string `json: "id"`
	Seller      string `json: "seller"`
	Contact     string `json: "contact"`
	StartRow    string `json: "startRow"`
	ExcludeRows string `json: "excludeRows"`
	TeamNumCol  string `json: "teamNumCol"`
	TeamNameCol string `json: "teamNameCol"`
	SizeCol     string `json: "sizeCol"`
	YearCol     string `json: "yearCol"`
	DescCol     string `json: "descCol"`
}

const MaxInt int = int(^uint(0) >> 1)
const MinInt int = -MaxInt - 1

func main() {
	// Check for year in passed args, error if not found
	const yearErrMsg string = "Required 4-digit year parameter not found"
	if len(os.Args) < 1 {
		panic(yearErrMsg)
	}
	var year *string = nil
	for _, arg := range os.Args {
		if len(arg) == 4 {
			year = &arg
			break
		}
	}
	if year == nil {
		panic(yearErrMsg)
	}

	// Parse sheets JSON, create output CSV, auth Google Sheets
	sheets := parseJson(fmt.Sprintf("in/shirt-sheets-%s.json", *year))
	outFile, err := os.Create(fmt.Sprintf("out/shirts-%s.csv", *year))
	handleErr(err, true, "Cannot create out CSV file")
	csv := csv.NewWriter(outFile)
	srv := gutilInit()

	// Write CSV headers
	csv.Write([]string{
		"Team Number",
		"Team Name",
		"Size",
		"Year",
		"Description",
		"Seller",
		"Contact",
	})

	// Iterate through sheets in JSON
	totalListings := 0
	for _, sheet := range sheets {
		// Use start row if key/value exists in JSON
		startRow := byte(49) // Digit "1" decimal ASCII
		if len(sheet.StartRow) != 0 {
			startRow = sheet.StartRow[0]
		}

		// Find min and max columns to retrieve
		minColNum := MaxInt
		maxColNum := MinInt
		findMinMaxCol([]string{
			sheet.TeamNumCol,
			sheet.TeamNameCol,
			sheet.SizeCol,
			sheet.YearCol,
			sheet.DescCol,
		}, &minColNum, &maxColNum)
		sheetRange := fmt.Sprintf("%c%c:%c", minColNum, startRow, maxColNum)
		fmt.Printf("Retrieving data of range %s for %s\n", sheetRange, sheet.ID)

		resp, err := srv.Spreadsheets.Values.Get(sheet.ID, sheetRange).Do()
		handleErr(err, true, "Cannot retrieve data for "+sheet.ID)
		// Iterate through rows (shirts)
		for rIdx, row := range resp.Values {
			numericIdx, _ := strconv.Atoi(fmt.Sprintf("%c", startRow))
			if isExcluded(rIdx+numericIdx, sheet.ExcludeRows) {
				fmt.Printf("Row %d was marked as excluded. Skipping.\n", rIdx)
				continue
			}
			tNum := parseValues(&row, minColNum, sheet.TeamNumCol)
			tName := parseValues(&row, minColNum, sheet.TeamNameCol)
			desc := parseValues(&row, minColNum, sheet.DescCol)
			// Check for combined team number/name(s)
			if sheet.TeamNumCol == sheet.TeamNameCol && strings.Contains(tNum, "-") {
				tNum = strings.Trim(tNum[:strings.Index(tNum, "-")], " ")
				tName = strings.Trim(tName[strings.Index(tName, "-")+1:], " ")
			}

			if strings.TrimSpace(tNum) == "" && strings.TrimSpace(tName) == "" && strings.TrimSpace(desc) == "" {
				fmt.Printf("Data array for row %d was empty. Skipping.\n", rIdx)
				continue
			}

			// Write shirt row to CSV
			err := csv.Write([]string{
				tNum,
				tName,
				parseValues(&row, minColNum, sheet.SizeCol),
				parseValues(&row, minColNum, sheet.YearCol),
				desc,
				sheet.Seller,
				sheet.Contact,
			})
			handleErr(err, true, fmt.Sprintf("Could not write data for %s to file", sheet.ID))
			csv.Flush()
		}
		totalListings += len(resp.Values)
		fmt.Printf("%d listings for seller \"%s\" written to CSV\n\n", len(resp.Values), sheet.Seller)
	}
	fmt.Printf("\n%d total listings written for %d sellers", totalListings, len(sheets))
	defer outFile.Close()
}

// Parse JSON input file into shirt sheet struct
func parseJson(path string) []shirtSheet {
	data := []shirtSheet{}
	file, err := ioutil.ReadFile(path)
	handleErr(err, true, "Could not read input file at "+path)
	err = json.Unmarshal([]byte(file), &data)
	handleErr(err, true, "JSON unmarshal failed")
	return data
}

// Parse values within a row given first column letter and current column
func parseValues(row *[]interface{}, fcol int, col string) string {
	switch len(col) {
	case 1:
		idx := int(col[0])
		idx -= fcol
		if idx < len(*row) {
			return (*row)[idx].(string)
		}
		// Fallthrough intentional
	case 0:
		return ""
	default:
		return parseValues(row, fcol, string(col[0])) + " " + parseValues(row, fcol, col[2:])
	}
	return ""
}

// Find the maximum and minimum column values
func findMinMaxCol(colChars []string, min *int, max *int) {
	for _, cc := range colChars {
		for idx := range cc {
			ccNum := int(cc[idx])
			// If cc is an uppercase ASCII letter
			if ccNum > 64 && ccNum < 91 {
				if ccNum < *min {
					*min = ccNum
				} else if ccNum > *max {
					*max = ccNum
				}
			}
		}
	}
}

// Check if a row index is excluded by a row-exclude range
func isExcluded(idx int, exclRows string) bool {
	const errStr string = "Could not convert string to int"
	// Split by comma
	terms := strings.Split(exclRows, ",")
	for _, term := range terms {
		// Check for range else single value
		if strings.Contains(term, ":") {
			var err error
			colonIdx := strings.Index(term, ":")

			// Lower bound of range
			lowTerm := term[:colonIdx]
			var lower int
			if len(lowTerm) == 0 {
				lower = 0
			} else {
				lower, err = strconv.Atoi(lowTerm)
				handleErr(err, true, errStr)
			}

			// Upper bound of range
			highTerm := term[colonIdx+1:]
			var upper int
			if len(highTerm) == 0 {
				upper = MaxInt
			} else {
				upper, err = strconv.Atoi(highTerm)
				handleErr(err, true, errStr)
			}

			// Value is between lower and upper bound, rtn
			if idx >= lower && idx <= upper {
				return true
			}
		} else if len(term) != 0 {
			// Value is equal to term, rtn
			value, err := strconv.Atoi(term)
			handleErr(err, true, errStr)
			if idx == value {
				return true
			}
		}
	}
	return false
}

// Handle errors with panic or without
func handleErr(err error, doPanic bool, msg string) {
	if err != nil {
		if doPanic {
			panic(msg)
		} else {
			fmt.Println("ERROR: " + msg)
		}
	}
}
