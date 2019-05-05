package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func DeleteOutputFile(filePath string) {
	var err = os.Remove(GENERATED_OUTPUT_DIR + filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func VisitDirectory(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if match, _ := regexp.MatchString(FILE_PATTERN, info.Name()); !match {
			return nil
		}

		*files = append(*files, path)
		return nil
	}
}

func IsSliceContain(slices []interface{}, itf interface{}) bool {
	for _, s := range slices {
		if s == itf {
			return true
		}
	}

	return false
}

func DelimSliceToString(data []string, delim string) string {
	str := strings.Trim(strings.Join(data, delim), "[]")
	lChar := str[len(str)-1:]
	if lChar == delim {
		str = str[:len(str)-1]
	}

	return str
}

func GenerateCommandOld(logPath string, filters []Filter, fromType int) {
	cmds := []*exec.Cmd{}
	cmds = append(cmds, exec.Command("zcat", logPath))
	for _, filter := range filters {
		cmds = append(cmds, exec.Command("grep", filter.GrepType.FormatV2(filter.Value)...))
	}

	output, err := PipeCommands(cmds...)
	if err != nil {
		// ignore
	} else {
		switch fromType {
		case NOTIFY_SUCCESS:
			AppendToFile(GENERATED_NOTIFY_SUCCESS_FILENAME, ParseJSONRequestOnly(string(output)))
		case VALIDATE_USE:
			AppendToFile(GENERATED_VALIDATE_USE_FILENAME, ParseJSONRequestOnly(string(output)))
		}
	}
}

func GenerateCommand(logPath string, filters []Filter, fromType int) {
	cmd := fmt.Sprintf("zcat '%s'", logPath)
	for _, filter := range filters {
		cmd += fmt.Sprintf(" | "+filter.GrepType.Format(), filter.Value)
	}

	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		// ignore
	} else {
		switch fromType {
		case NOTIFY_SUCCESS:
			AppendToFile(GENERATED_NOTIFY_SUCCESS_FILENAME, ParseJSONRequestOnly(string(output)))
		case VALIDATE_USE:
			AppendToFile(GENERATED_VALIDATE_USE_FILENAME, ParseJSONRequestOnly(string(output)))
		}

	}
}

func PipeCommands(commands ...*exec.Cmd) ([]byte, error) {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		return nil, err
	}
	return final, nil
}

func ParseJSONRequestOnly(str string) string {
	strb := strings.Builder{}
	regex, _ := regexp.Compile(`({\\"data\\"\:[a-zA-Z0-9\,\.\s\+\-\_\{\}\"\\\:\[\]\/\(\)\@\;\&\*\'\?\#]*)]"}(\r?\n)`)
	if regex.MatchString(str) {
		var getParsing = regex.FindAllStringSubmatch(str, -1)
		for _, group := range getParsing {
			if len(group) > 0 {
				strb.WriteString(fmt.Sprintln(strings.Replace(group[1], "\\", "", -1)))
			}
		}
	}

	return strb.String()
}

func AppendToFile(fileName string, buffer string) {

	fmt.Printf("... Append result to %s ...\n", fileName)

	f, err := os.OpenFile(GENERATED_OUTPUT_DIR+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(buffer + "\n"); err != nil {
		log.Println(err)
	}
}
