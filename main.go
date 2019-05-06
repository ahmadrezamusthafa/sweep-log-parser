package main

import (
	"bufio"
	"fmt"
	"github.com/json-iterator/go"
	. "github.com/tokopedia/sweep-log/core"
	"github.com/tokopedia/sweep-log/core/enum"
	"log"
	"marisinau.com/KitaUndangAPI/util/errors"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	PaymentIDs []string
	PromoCodes []string
	LogPaths   []string
)

func main() {

	deleteOutputFile() // delete first

	if err := readInputOrderFile(); err != nil {
		return
	}

	if err := scanLogDirectory(MODE_ALL_IN_ONE); err != nil {
		return
	}

	if err := getNotifySuccessData(); err != nil {
		return
	}

	if err := getValidateUseData(); err != nil {
		return
	}

	if err := processOutput(); err != nil {
		return
	}

	fmt.Println("\n=== FINISH YOOOOOOOOOOOO CAKKK !!! ===")
}

func deleteOutputFile() {
	DeleteOutputFile(GENERATED_VALIDATE_USE_FILENAME)
	DeleteOutputFile(GENERATED_NOTIFY_SUCCESS_FILENAME)
	DeleteOutputFile(GENERATED_CSV_FILENAME)
}

func readInputOrderFile() error {

	fmt.Println("\n=== Getting code from order ===")

	file, err := os.Open(LOG_PARENT_DIR + ORDER_INPUT_FILENAME)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	regex, _ := regexp.Compile(`global_level_orders[^0-9A-Za-z]*([A-Za-z0-9]*)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rowStr := scanner.Text()

		if rowStr == "" {
			continue
		}

		rArr := strings.Split(rowStr, "|")
		if len(rArr) > 0 {
			paymentID := rArr[0]
			PaymentIDs = append(PaymentIDs, paymentID)
		}

		if regex.MatchString(rowStr) {
			var getParsing = regex.FindAllStringSubmatch(rowStr, 1)
			for _, group := range getParsing {
				if len(group) > 0 {
					code := group[1]

					s := make([]interface{}, len(PromoCodes))
					for i, v := range PromoCodes {
						s[i] = v
					}

					if !IsSliceContain(s, code) {
						PromoCodes = append(PromoCodes, code)
					}

					break
				}
			}
		}
	}

	fmt.Println(".:. TOTAL PROMO CODE FROM ORDER :", len(PromoCodes))

	return nil
}

func scanLogDirectory(mode int) error {

	fmt.Println("\n=== Scanning log directory ===")

	err := filepath.Walk(LOG_PARENT_DIR, VisitDirectory(&LogPaths, mode))
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println(".:. TOTAL SCANNED LOG FILE :", len(LogPaths))

	return nil
}

func getNotifySuccessData() error {

	fmt.Println("\n=== Get notify success data ===")

	filters := []Filter{
		{GrepType: enum.GrepStandard, Value: "NOTIFY USE SUCCESS"},
		{GrepType: enum.GrepCombined, Value: DelimSliceToString(PaymentIDs, "|")},
		//{GrepType: enum.GrepStandard, Value: "2019-04-29"},
	}

	for i, logPath := range LogPaths {
		fmt.Printf("=== Processing NS [%d:%d] %s\n", (i + 1), len(LogPaths), logPath)
		GenerateCommand(logPath, filters, NOTIFY_SUCCESS, MODE_ALL_IN_ONE)
	}

	return nil
}

func getValidateUseData() error {

	fmt.Println("\n=== Get validate use data ===")

	filters := []Filter{
		{GrepType: enum.GrepStandard, Value: "VALIDATE USE"},
		{GrepType: enum.GrepCombined, Value: DelimSliceToString(PromoCodes, "|")},
		{GrepType: enum.GrepStandard, Value: "\\\\\"book\\\\\":true"},
		//{GrepType: enum.GrepStandard, Value: "2019-04-29"},
	}

	for i, logPath := range LogPaths {
		fmt.Printf("=== Processing VU [%d:%d] %s\n", (i + 1), len(LogPaths), logPath)
		GenerateCommand(logPath, filters, VALIDATE_USE, MODE_ALL_IN_ONE)
	}

	return nil
}

func processOutput() error {

	if !IsOutputFileExist(GENERATED_NOTIFY_SUCCESS_FILENAME) && !IsOutputFileExist(GENERATED_VALIDATE_USE_FILENAME) {
		return errors.New(fmt.Sprintf("Failed, make sure %s and %s is exist"), GENERATED_NOTIFY_SUCCESS_FILENAME, GENERATED_VALIDATE_USE_FILENAME)
	}

	fmt.Println("\n=== Process final output ===")

	file, err := os.Open(GENERATED_OUTPUT_DIR + GENERATED_NOTIFY_SUCCESS_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()

		if row == "" {
			continue
		}

		var nsRow UseCodeInput
		err := jsoniter.Unmarshal([]byte(row), &nsRow)
		if err != nil {
			fmt.Println("PROBLEM WITH DATA NS:", row)
			continue
		}

		filters := []Filter{
			{GrepType: enum.GrepStandard, Value: fmt.Sprintf("\"user_id\":%d", nsRow.Data.UserData.UserID)},
			{GrepType: enum.GrepStandard, Value: fmt.Sprintf("\"grand_total\":%.f", math.Round(nsRow.Data.PaymentAmount))},
		}

		vuOutput := GenerateCommandFinalProcess(GENERATED_OUTPUT_DIR+GENERATED_VALIDATE_USE_FILENAME, filters)
		if vuOutput != "" {
			vuArr := strings.Split(vuOutput, "\n")
			for _, vu := range vuArr {

				if vu == "" {
					continue
				}

				var vuRow UseCodeInput
				err := jsoniter.Unmarshal([]byte(vu), &vuRow)
				if err != nil {
					fmt.Println("PROBLEM WITH DATA VU:", vu)
					continue
				}

				if nsRow.Data.UserData.UserID == vuRow.Data.UserData.UserID && nsRow.Data.PaymentAmount == vuRow.Data.GrandTotal {
					AppendToFile(GENERATED_CSV_FILENAME, fmt.Sprintf("%s~~~~~%s\n", row, vu))
					break
				}
			}
		}
	}

	return nil
}
