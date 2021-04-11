package main

import (
	"bufio"
	"fmt"
	"log"
	"lrn/choose-your-adventure/adventure"
	"lrn/choose-your-adventure/cli/page"
	"lrn/choose-your-adventure/cli/page/pages"
	"os"
	"strings"
)

const ADVENTURES_DIR = "ADVENTURES_DIR"

func main() {
	indexPage, err := getAdventureIndexPage()
	if err != nil {
		log.Fatal(err)
	}
	Run(indexPage)
}

func getAdventureIndexPage() (*pages.AdventuresIndexPage, error) {
	advDir := os.Getenv(ADVENTURES_DIR)
	if advDir == "" {
		return nil, fmt.Errorf("please, provide %s env var", ADVENTURES_DIR)
	}

	advLoader := adventure.NewFsAdventureLoader(advDir)
	adv, err := advLoader.LoadAdventures()
	if err != nil {
		return nil, fmt.Errorf("couldn't load adventures: %s", err)
	}
	indexPage := &pages.AdventuresIndexPage{
		Adventures: adv,
	}
	return indexPage, nil
}

func Run(indexPage page.Page) {
	inputCh := createInputChannel()

	var page page.Page = indexPage
	printStartMessage()
loop:
	for {
		printSep()
		page.Display()
		for input := range inputCh {
			switch input {
			case "!exit":
				return
			case "!index":
				page = indexPage
			case "!help":
				printHelpPage(page)

			default:
				if nextPage := page.HandleInput(input); nextPage != nil {
					page = nextPage
					continue loop
				}
			}
			printSep()
			page.Display()
		}
	}
}

func printSep() {
	fmt.Println("-------------------")
}

func printStartMessage() {
	fmt.Println("Hello! You can get help at any time by printing !help.")
}

func printHelpPage(page page.Page) {
	printSep()
	fmt.Println("!exit: exit from app")
	fmt.Println("!index: return to index page")
	fmt.Println("!help: show help")
	fmt.Println(page.GetHelpMessage())
}

// createInputChannel returns channel with user's input line.
// Input is already trimmed from spaces
func createInputChannel() <-chan string {
	ch := make(chan string)
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			text = strings.TrimSpace(text)
			ch <- text
		}
	}()
	return ch
}
