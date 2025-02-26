package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"

	"strconv"

	"os"

	"sort"
)

var quiz_avg float64
var labtest_avg float64
var midsem_avg float64
var weeklab_avg float64
var pct_avg float64
var compre_avg float64
var total_avg float64
var total_avg_2024 map[string]float64
var n_students_2024 map[string]int

type Student struct {
	Emplid string

	quiz_marks    float64
	midsem_marks  float64
	labtest_marks float64
	weeklab_marks float64
	compre_marks  float64
	total_marks   float64
}

func display(x Student, rank int, component string) {
	fmt.Printf("Emplid = %s\n", x.Emplid)
	switch component {
	case "quiz":
		fmt.Printf("Marks in %s = %f\n", component, x.quiz_marks)
	case "midsem":
		fmt.Printf("Marks in %s = %f\n", component, x.midsem_marks)
	case "labtest":
		fmt.Printf("Marks in %s = %f\n", component, x.labtest_marks)
	case "weekly lab":
		fmt.Printf("Marks in %s = %f\n", component, x.weeklab_marks)
	case "compre":
		fmt.Printf("Marks in %s = %f\n", component, x.compre_marks)
	case "total":
		fmt.Printf("Marks in %s = %f\n", component, x.total_marks)
	}

	fmt.Printf("Rank in %s = %d\n", component, rank)
}

func add_student(all_students *[]Student, row []string) {
	emplid := row[2]
	quiz := strint(row[4])
	midsem := strint(row[5])
	labtest := strint(row[6])
	weeklab := strint(row[7])
	compre := strint(row[9])
	total := quiz + midsem + labtest + weeklab + compre
	*all_students = append(*all_students, Student{emplid, quiz, midsem, labtest, weeklab, compre, total})
}

func strint(x string) float64 {

	y, err := strconv.ParseFloat(x, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return y
}

func branch_wise_avg(row []string) {
	id := row[3]
	year := id[:4]

	quiz := strint(row[4])
	midsem := strint(row[5])
	labtest := strint(row[6])
	weeklab := strint(row[7])
	compre := strint(row[9])

	if year == "2024" {
		key := id[:6]
		total_avg_2024[key] = total_avg_2024[key] + quiz + midsem + labtest + weeklab + compre
		n_students_2024[key] = n_students_2024[key] + 1
	}
}

func calculate_avg(row []string) {

	quiz := strint(row[4])
	quiz_avg += quiz

	midsem := strint(row[5])
	midsem_avg += midsem

	labtest := strint(row[6])
	labtest_avg += labtest

	weeklab := strint(row[7])
	weeklab_avg += weeklab

	compre := strint(row[9])
	compre_avg += compre
}

// floating point error
func validate_row(row []string, rownum int, errors *[]string) bool {

	good := true
	quiz := strint(row[4])

	midsem := strint(row[5])

	labtest := strint(row[6])

	weeklab := strint(row[7])

	pct := strint(row[8])

	compre := strint(row[9])

	total := strint(row[10])

	if quiz+midsem+labtest+weeklab != pct {
		*errors = append(*errors, fmt.Sprintf("%s%d%s%s", "Incorrect pct calculation in row ", rownum, " with Sl No. ", row[0]))
		good = false
	}

	if quiz+midsem+labtest+weeklab+compre != total {
		*errors = append(*errors, fmt.Sprintf("%s%d%s%s", "Incorrect pct calculation in row ", rownum, " with Sl No. ", row[0]))
		good = false
	}

	return good
}
func is_row_empty(row []string) bool {

	for _, cell := range row {
		if cell != "" {
			return false
		}
	}

	return true
}
func main() {

	n_students_2024 = make(map[string]int)
	total_avg_2024 = make(map[string]float64)
	all_students := []Student{}

	quiz_avg = 0
	labtest_avg = 0
	midsem_avg = 0
	weeklab_avg = 0
	pct_avg = 0
	compre_avg = 0
	total_avg = 0

	rownum := 1
	error_rows := []int{} // stores the row number which has an error
	errors := []string{}  // stores the errors at a row if any
	empty_rows := []string{}

	var fileName string
	fmt.Println("Enter the path of the data file:")
	fmt.Scanln(&fileName)

	f, err := excelize.OpenFile(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := f.GetRows("CSF111_202425_01_GradeBook")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("The spreadsheet is diaplayed below")
	for _, row := range rows {

		for _, col := range row {
			fmt.Print(col, "\t")
		}
		fmt.Println()

		if rownum != 1 {

			if is_row_empty(row) {
				empty_rows = append(empty_rows, fmt.Sprintf("Row number %d is empty", rownum))
				rownum++
				continue
			} // check for empty rows

			if !validate_row(row, rownum, &errors) {
				error_rows = append(error_rows, rownum)
			}

			calculate_avg(row)
			branch_wise_avg(row)
			add_student(&all_students, row)
		}

		rownum++
	}

	quiz_avg /= float64(rownum - 2)
	midsem_avg /= float64(rownum - 2)
	labtest_avg /= float64(rownum - 2)
	weeklab_avg /= float64(rownum - 2)
	compre_avg /= float64(rownum - 2)

	pct_avg = quiz_avg + midsem_avg + labtest_avg + weeklab_avg
	total_avg = pct_avg + compre_avg

	if len(empty_rows) != 0 {

		fmt.Println("The following rows are empty")

		for _, item := range empty_rows {
			fmt.Println(item)
		}

	} else {
		fmt.Println("There are no empty rows")
	}

	fmt.Println()

	fmt.Println("The following rows contain errors")
	fmt.Println(error_rows)

	fmt.Println()
	fmt.Println("The errors are logged below")

	for _, item := range errors {
		fmt.Println(item)
	}

	fmt.Println()

	fmt.Println("The component averages are displayed")
	fmt.Printf("The quiz average = %f\n", quiz_avg)
	fmt.Printf("The midsem average = %f\n", midsem_avg)
	fmt.Printf("The labtest average = %f\n", labtest_avg)
	fmt.Printf("The weekly lab average = %f\n", weeklab_avg)
	fmt.Printf("The pct average = %f\n", pct_avg)
	fmt.Printf("The compre average = %f\n", compre_avg)
	fmt.Printf("The total average = %f\n", total_avg)

	fmt.Println()
	fmt.Println("The branch-wise average for the 2024 batch are displayed")

	for key := range total_avg_2024 {
		fmt.Printf("The total average for the 2024 %s batch = %f\n", key[4:6], total_avg_2024[key]/float64(n_students_2024[key]))
	}

	fmt.Println()

	fmt.Println("Rank holders of each component")

	fmt.Println("Rank-holders of quiz")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].quiz_marks > all_students[j].quiz_marks })
	for i := 0; i < 3; i++ {
		display(all_students[i], i+1, "quiz")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of midsem")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].midsem_marks > all_students[j].midsem_marks })
	for i := 0; i < 3; i++ {
		display(all_students[i], i+1, "midsem")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of labtest")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].labtest_marks > all_students[j].labtest_marks })
	for i := 0; i < 3; i++ {
		display(all_students[i], i+1, "labtest")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of weekly labs")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].weeklab_marks > all_students[j].weeklab_marks })
	for i := 0; i < 3; i++ {
		display(all_students[i], i+1, "weekly lab")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of compre")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].compre_marks > all_students[j].compre_marks })
	for i := 0; i < 3; i++ {
		display(all_students[i], i+1, "compre")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of overall total")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].total_marks > all_students[j].total_marks })
	for i := 0; i < 3; i++ {
		display(all_students[i], i+1, "total")
		fmt.Println("-----------------------")
	}
}
