package utils

import "fmt"

func RenderSettingPage(title string, options []string) int {
	// Clear console screen and print title
	fmt.Println("\033[H\033[2J")
	fmt.Println(title)
	fmt.Println()

	// Print options
	for i, option := range options {
		fmt.Printf("[%d] %s\n", i+1, option)
	}

	// Print the exit option
	fmt.Println("[0] Exit")

	// Get user input, if 0 is returned, the user wants to exit, 1-9 are valid options
	var input int
	fmt.Scan(&input)

	// Clear console screen
	fmt.Println("\033[H\033[2J")

	if input == 0 {
		return -1
	}

	if input < 1 || input > len(options) {
		return -2
	}

	return input - 1
}

func GetSetupValue(text string) string {
	fmt.Printf("%s\n Enter value: ", text)

	var input string
	fmt.Scan(&input)

	return input
}
