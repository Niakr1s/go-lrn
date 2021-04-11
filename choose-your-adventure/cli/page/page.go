package page

type Page interface {
	Display()
	GetHelpMessage() string
	HandleInput(input string) Page
}
