package game

import (
	"fmt"
	"strings"
)

func prompt(p *Player, msg string) {
	fmt.Printf("\n%s\n> %s:", msg, p.Name)
}

func getUserTextInput(p *Player, msg string, choices []string) string {
	var input string
	msg = fmt.Sprintf("%s. %s", msg, strings.Join(choices, ", "))
	for {
		prompt(p, msg)
		fmt.Scanln(&input)
		for _, choice := range choices {
			if strings.EqualFold(input, choice) {
				return strings.ToLower(input)
			}
		}
		fmt.Printf("Invalid choice. Enter one of: %s\n", strings.Join(choices, ", "))
	}

}

func getUserMinNumberInputOrPass(p *Player, msg string, min int) int {
	var input int
	for {
		prompt(p, msg)
		fmt.Scanln(&input)
		if input == 0 {
			return 0
		}
		if input >= min {
			return input
		}
		fmt.Println("Invalid choice, number must be equal to greater than ", min)
	}
}
