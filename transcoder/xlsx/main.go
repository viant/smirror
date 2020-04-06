package xlsx

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main()  {

	data, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), "Desktop", "test.xlsx"))
	if err != nil {
		log.Fatal(err)
	}
	file, err  := xlsx.OpenBinary(data)
	if err != nil {
		log.Fatal(err)
	}

	if len(file.Sheets) == 0 {
		return
	}
	sheet := file.Sheets[0]

	for i, row := range sheet.Rows {
		fmt.Printf("[%v] -> %+v\n", i, len(row.Cells))

		for j, _ := range row.Cells {
			fmt.Printf("%v (%v) |  ", row.Cells[j].Value,  row.Cells[j].Type())
		}
	}

	//CellTypeStringFormula
	//	CellTypeNumeric
	//	CellTypeBool
	//	// CellTypeInline is not respected on save, all inline string cells will be saved as SharedStrings
	//	// when saving to an XLSX file. This the same behavior as that found in Excel.
	//	CellTypeInline
	//	CellTypeError
	//	// d (Date): Cell contains a date in the ISO 8601 format.
	//	// That is the only mention of this format in the XLSX spec.
	//	// Date seems to be unused by the current version of Excel, it stores dates as Numeric cells with a date format string.
	//	// For now these cells will have their value output directly. It is unclear if the value is supposed to be parsed
	//	// into a number and then formatted using the formatting or not.
	//	CellTypeDate

}


