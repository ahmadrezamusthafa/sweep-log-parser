package main

import (
	"bufio"
	"fmt"
	. "github.com/tokopedia/sweep-log/core"
	"github.com/tokopedia/sweep-log/core/enum"
	"log"
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

	if err := scanLogDirectory(); err != nil {
		return
	}

	if err := getNotifySuccessData(); err != nil {
		return
	}

	if err := getValidateUseData(); err != nil {
		return
	}

	fmt.Println("\n=== FINISH YOOOOOOOOOOOO CAKKK !!! ===")
}

func deleteOutputFile(){
	DeleteOutputFile(GENERATED_VALIDATE_USE_FILENAME)
	DeleteOutputFile(GENERATED_NOTIFY_SUCCESS_FILENAME)
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

func scanLogDirectory() error {

	fmt.Println("\n=== Scanning log directory ===")

	err := filepath.Walk(LOG_PARENT_DIR, VisitDirectory(&LogPaths))
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
		GenerateCommand(logPath, filters, NOTIFY_SUCCESS)
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
		GenerateCommand(logPath, filters, VALIDATE_USE)
	}

	return nil
}


