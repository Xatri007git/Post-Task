package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/xuri/excelize/v2"

	"strconv"

	"os"

	"math"
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
var class int64
var js string
var e float64 // margin of error for addition

type Student struct {
	Emplid string

	quiz_marks    float64
	midsem_marks  float64
	labtest_marks float64
	weeklab_marks float64
	compre_marks  float64
	total_marks   float64
}

type Report struct {
	Errors          []string `json:"Errors"`
	Component_avgs  []string `json:"Component Averages"`
	Branch_avg_2024 []string `json:"Branchwise average 2024 batch"`
	Rank_quiz       []string `json:"Rank holders for quiz"`
	Rank_midsem     []string `json:"Rank holders for midsem"`
	Rank_labtest    []string `json:"Rank holders for labtest"`
	Rank_weekly_lab []string `json:"Rank holders for weekly labs"`
	Rank_compre     []string `json:"Rank holders for compre"`
	Rank_overall    []string `json:"Rank holders for overall total"`
}

func create_json(errors []string, comp_avs []string, branch_avs []string,
	rank_q []string, rank_m []string, rank_l []string, rank_wl []string, rank_c []string, rank_o []string, filepath string) {
	report := Report{errors, comp_avs, branch_avs, rank_q, rank_m, rank_l, rank_wl, rank_c, rank_o}
	loc := []Report{report}

	finaljson, err := json.MarshalIndent(loc, "", "\t")
	if err != nil {
		fmt.Println("Could not create json file")
		os.Exit(1)
	}

	writefile := os.WriteFile(filepath, finaljson, 0644)
	if writefile != nil {
		fmt.Println("Error in writing json file")
		os.Exit(1)
	}

	fmt.Println("The json file is displayed below")
	fmt.Printf("%s\n", finaljson)
	fmt.Printf("The above json file has been exported at %s", filepath)
}

func read_flags() {
	flag.Int64Var(&class, "class", -1, "class No. of the student")
	flag.StringVar(&js, "export", "", "export to json")
	flag.Parse()
}

func class_fliter(row []string) bool {
	if class == -1 {
		return true // no filter
	}

	return strInt(row[1]) == class
}

func export_filter() bool {
	if js != "json" && js != "" {
		fmt.Println("Invalid value for export flag")
		os.Exit(1)
	}
	return js == "json"
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
	quiz := strfloat(row[4])
	midsem := strfloat(row[5])
	labtest := strfloat(row[6])
	weeklab := strfloat(row[7])
	compre := strfloat(row[9])
	total := quiz + midsem + labtest + weeklab + compre
	*all_students = append(*all_students, Student{emplid, quiz, midsem, labtest, weeklab, compre, total})
}

func strfloat(x string) float64 {

	y, err := strconv.ParseFloat(x, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return y
}

func strInt(x string) int64 {

	y, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return y
}

func branch_wise_avg(row []string) {
	id := row[3]
	year := id[:4]

	quiz := strfloat(row[4])
	midsem := strfloat(row[5])
	labtest := strfloat(row[6])
	weeklab := strfloat(row[7])
	compre := strfloat(row[9])

	if year == "2024" {
		key := id[:6]
		total_avg_2024[key] = total_avg_2024[key] + quiz + midsem + labtest + weeklab + compre
		n_students_2024[key] = n_students_2024[key] + 1
	}
}

func calculate_avg(row []string) {

	quiz := strfloat(row[4])
	quiz_avg += quiz

	midsem := strfloat(row[5])
	midsem_avg += midsem

	labtest := strfloat(row[6])
	labtest_avg += labtest

	weeklab := strfloat(row[7])
	weeklab_avg += weeklab

	compre := strfloat(row[9])
	compre_avg += compre
}

// floating point error
func validate_row(row []string, errors *[]string) bool {

	good := true
	quiz := strfloat(row[4])

	midsem := strfloat(row[5])

	labtest := strfloat(row[6])

	weeklab := strfloat(row[7])

	pct := strfloat(row[8])

	compre := strfloat(row[9])

	total := strfloat(row[10])

	if math.Abs(quiz+midsem+labtest+weeklab-pct) > e {
		*errors = append(*errors, fmt.Sprintf("%s%s", "Incorrect pct calculation in row with Sl No. ", row[0]))
		good = false
	}

	if math.Abs(quiz+midsem+labtest+weeklab+compre-total) > e {
		*errors = append(*errors, fmt.Sprintf("%s%s", "Incorrect total calculation in row with Sl No. ", row[0]))
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

func display_row(row []string) {
	for _, col := range row {
		fmt.Print(col, "\t")
	}
	fmt.Println()
}

func main() {

	read_flags()
	e = 0.000001
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
	error_rows := []string{} // stores the row number which has an error
	errors := []string{}     // stores the errors at a row if any
	empty_rows := []string{}
	branch_avg_2024 := []string{}
	rank_quiz := []string{}
	rank_midsem := []string{}
	rank_labtest := []string{}
	rank_weekly_lab := []string{}
	rank_compre := []string{}
	rank_overall := []string{}

	var fileName string
	fmt.Println("Enter the path of the data file:")
	fmt.Scanln(&fileName)

	f, err := excelize.OpenFile(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := f.GetRows("CSF111_202425_01_GradeBook")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("The spreadsheet is diaplayed below")
	for _, row := range rows {

		if rownum == 1 {
			display_row(row) // display the description row
			rownum++
		} else {

			if class_fliter(row) {
				display_row(row)
				if is_row_empty(row) {
					empty_rows = append(empty_rows, fmt.Sprintf("Row with Sl No. %s is empty", row[0]))
					rownum++
					continue
				} // check for empty rows

				if !validate_row(row, &errors) {
					error_rows = append(error_rows, row[0])
				}

				calculate_avg(row)
				branch_wise_avg(row)
				add_student(&all_students, row)
				rownum++
			}
		}

	}

	if rownum == 2 {
		fmt.Printf("Class flag is invalid, there is no entry in class No. with value %d\n", class)
		os.Exit(1)
	}

	fmt.Println()

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

	fmt.Println("The rows with the following Sl No. contain errors")
	fmt.Println(error_rows)

	fmt.Println()
	fmt.Println("The errors are logged below")

	for _, item := range errors {
		fmt.Println(item)
	}

	fmt.Println()

	fmt.Println("The component averages are displayed")

	qa := fmt.Sprintf("The quiz average = %f", quiz_avg)
	fmt.Printf("The quiz average = %f\n", quiz_avg)

	ma := fmt.Sprintf("The midsem average = %f", midsem_avg)
	fmt.Printf("The midsem average = %f\n", midsem_avg)

	la := fmt.Sprintf("The labtest average = %f", labtest_avg)
	fmt.Printf("The labtest average = %f\n", labtest_avg)

	wa := fmt.Sprintf("The weekly lab average = %f", weeklab_avg)
	fmt.Printf("The weekly lab average = %f\n", weeklab_avg)

	pa := fmt.Sprintf("The pct average = %f", pct_avg)
	fmt.Printf("The pct average = %f\n", pct_avg)

	ca := fmt.Sprintf("The compre average = %f", compre_avg)
	fmt.Printf("The compre average = %f\n", compre_avg)

	ta := fmt.Sprintf("The total average = %f", total_avg)
	fmt.Printf("The total average = %f\n", total_avg)

	Component_avgs := []string{qa, ma, la, wa, pa, ca, ta}

	fmt.Println()
	fmt.Println("The branch-wise average for the 2024 batch are displayed")

	for key := range total_avg_2024 {
		s := fmt.Sprintf("Total average for the 2024 %s batch = %f", key[4:6], total_avg_2024[key]/float64(n_students_2024[key]))
		branch_avg_2024 = append(branch_avg_2024, s)
		fmt.Printf("The total average for the 2024 %s batch = %f\n", key[4:6], total_avg_2024[key]/float64(n_students_2024[key]))
	}

	fmt.Println()

	fmt.Println("Rank holders of each component")

	fmt.Println("Rank-holders of quiz")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].quiz_marks > all_students[j].quiz_marks })
	for i := 0; i < 3; i++ {
		rank_quiz = append(rank_quiz, fmt.Sprintf("%s%s%s%f%s%d", "Emplid = ", all_students[i].Emplid, ", Quiz Marks = ", all_students[i].quiz_marks, ", Rank = ", i+1))
		display(all_students[i], i+1, "quiz")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of midsem")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].midsem_marks > all_students[j].midsem_marks })
	for i := 0; i < 3; i++ {
		rank_midsem = append(rank_midsem, fmt.Sprintf("%s%s%s%f%s%d", "Emplid = ", all_students[i].Emplid, ", Midsem Marks = ", all_students[i].midsem_marks, ", Rank = ", i+1))
		display(all_students[i], i+1, "midsem")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of labtest")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].labtest_marks > all_students[j].labtest_marks })
	for i := 0; i < 3; i++ {
		rank_labtest = append(rank_labtest, fmt.Sprintf("%s%s%s%f%s%d", "Emplid = ", all_students[i].Emplid, ", Labtest Marks = ", all_students[i].labtest_marks, ", Rank = ", i+1))
		display(all_students[i], i+1, "labtest")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of weekly labs")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].weeklab_marks > all_students[j].weeklab_marks })
	for i := 0; i < 3; i++ {
		rank_weekly_lab = append(rank_weekly_lab, fmt.Sprintf("%s%s%s%f%s%d", "Emplid = ", all_students[i].Emplid, ", Weekly lab Marks = ", all_students[i].weeklab_marks, ", Rank = ", i+1))
		display(all_students[i], i+1, "weekly lab")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of compre")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].compre_marks > all_students[j].compre_marks })
	for i := 0; i < 3; i++ {
		rank_compre = append(rank_compre, fmt.Sprintf("%s%s%s%f%s%d", "Emplid = ", all_students[i].Emplid, ", Compre Marks = ", all_students[i].compre_marks, ", Rank = ", i+1))
		display(all_students[i], i+1, "compre")
		fmt.Println("-----------------------")
	}
	fmt.Println()

	fmt.Println("Rank-holders of overall total")
	sort.Slice(all_students, func(i, j int) bool { return all_students[i].total_marks > all_students[j].total_marks })
	for i := 0; i < 3; i++ {
		rank_overall = append(rank_overall, fmt.Sprintf("%s%s%s%f%s%d", "Emplid = ", all_students[i].Emplid, ", Total Marks = ", all_students[i].total_marks, ", Rank = ", i+1))
		display(all_students[i], i+1, "total")
		fmt.Println("-----------------------")
	}

	if export_filter() {
		fmt.Println()

		fmt.Println("Enter the path of directory where you want to save json file")

		fpath := ""
		fmt.Scanln(&fpath)
		fname := ""

		fmt.Println("Enter desired file name followed by .json extension")
		fmt.Scanln(&fname)

		osfpath := filepath.FromSlash(fpath)
		create_json(errors, Component_avgs, branch_avg_2024, rank_quiz, rank_midsem, rank_labtest, rank_weekly_lab, rank_compre, rank_overall, filepath.Join(osfpath, fname))
	}
}
