package console

import (
	"UserTransferService/src/system/l2f"
	"UserTransferService/src/userService/cliHandler"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Reads console input and loads config to channel
func ReadConsole() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Shell")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			l2f.Log.Printf("reading console error: %s", err.Error())
		}
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		cliHandler.CommandCheck(text)
	}
}