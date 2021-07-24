# frc-shirt-aggregator

Retrieves all FRC shirt trades on the yearly ChiefDelphi thread.

Aggregates data into a CSV file with required metadata.

## Usage

1. Install `go` and set proper `GOPATH`, then clone in this repo
2. Copy `shirt-sheets-template.json` into `in/` and rename to `shirt-sheets-YEAR.json`
3. Modify data as needed, look to the included JSONs for examples
4. Run `go run github.com/boomaa23/frc-shirt-aggregator YEAR` (or `go build` as desired)

> NOTA BENE: Data may not be fully accurate. Check base spreadsheet before DMing sellers to ensure items are available and tradable.

## Data Format

### Input

- All values are string types
- Uses JSON formatting
- One copy of each entry is required per shirt spreadsheet
- State definitions
    - Required = program will fail without entry
    - Recommended = inclusion will improve accuracy
    - Optional - inclusion adds more data
- A copy of this as a template can be found [here](https://github.com/Boomaa23/frc-shirt-aggregator/blob/master/shirt-sheets-template.json)

| Name | Key | State | Value |
|------|-----|----------|-------|
| ID | id | required | Google sheet ID
| Seller | seller | recommended | Name of seller (First Last)
| Contact | contact | recommended | Contact info for seller
| Start Row | startRow | optional | Number of row to start parsing at
| Excluded Rows | excludeRows | optional | Rows to exclude from parsing
| Team Number Column | teamNumCol | recommended* | Letter of column containing team numbers
| Team Name Column | teamNameCol | recommended* | Letter of column containing team names
| Size Column | sizeCol | optional | Letter of column containing shirt sizes
| Year Column | yearCol | optional | Letter of column containing shirt years
| Description Column | descCol | recommended* | Letter of column containing item description

 - Excluded rows must be comma-separated inclusive ranges (ex `1:3`, `1:`, or `:3`)
 - At least one of: Team Number, Team Name, or Description must be included

### Output

| Team Number | Team Name | Size | Year | Description | Seller | Contact |
|-------------|-----------|------|------|-------------|--------|---------|
| Ex: 5818    | Riviera Robotics | M | 2019 | Logo Shirt | Boomaa23 | /u/Boomaa23 |


### Currently Included Data
| Year | CD Thread | JSON input | CSV output | Last Update |
|------|-----------|------------|------------|-------------|
| 2021 | [/t/390455](https://www.chiefdelphi.com/t/2021-shirt-trading-thread/390455/) | [in/shirt-sheets-2021.json](https://github.com/Boomaa23/frc-shirt-aggregator/blob/master/in/shirts-sheets-2021.json) | [out/shirts-2021.csv](https://github.com/Boomaa23/frc-shirt-aggregator/blob/master/out/shirts-2021.csv) | 07/24/2021
| 2020 | [/t/371821](https://www.chiefdelphi.com/t/2020-shirt-trading-thread/371821/) | [in/shirt-sheets-2020.json](https://github.com/Boomaa23/frc-shirt-aggregator/blob/master/in/shirts-sheets-2020.json) | [out/shirts-2020.csv](https://github.com/Boomaa23/frc-shirt-aggregator/blob/master/out/shirts-2020.csv) | 07/24/2021
| 2019 | [/t/335501](https://www.chiefdelphi.com/t/2019-shirt-trading-thread/335501/) | [in/shirt-sheets-2019.json](https://github.com/Boomaa23/frc-shirt-aggregator/blob/master/in/shirts-sheets-2019.json) | [out/shirts-2019.csv](https://github.com/Boomaa23/frc-shirt-aggregator/blob/master/out/shirts-2019.csv) | 07/24/2021