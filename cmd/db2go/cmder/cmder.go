package cmder

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readStdin(end byte) (strInputs []string) {
	var err error
	var v string
	reader := bufio.NewReader(os.Stdin)
	if v, err = reader.ReadString(end); err != nil {
		fmt.Printf("ReadString error %s", err.Error())
		return nil
	}

	ss := strings.Split(v, "\n")
	for _, v := range ss {
		s := strings.TrimSpace(v)
		if s == "" {
			continue //remove blank line
		}
		strInputs = append(strInputs, s)
	}
	return
}

//strPrompt 命令行提示符(prompt content)
//end 结束字符串(delimiter character)
func Prompt(strPrompt string, end byte) (strInputs []string) {
	fmt.Printf("%s", strPrompt)
	return readStdin(end)
}
