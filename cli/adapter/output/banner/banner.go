package banner

import (
	"bufio"
	"github.com/happierall/l"
	"log"
	"os"
)

func LoadBanner(bannerId string) []string {
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		log.Fatal(dirErr)
	}
	file, err := os.Open(dir + "/adapter/output/banner/" + bannerId + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var data []string
	s := bufio.NewScanner(file)
	for s.Scan() {
		data = append(data, s.Text())

		// other code what work with parsed line...
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}

func customColorize(text string, color string) string {
	var resetColor = "\033[39m"
	return color + text + resetColor
}

func PrintBanner(bannerDraw []string, color string) {

	if color == "yellow" {
		for i := 0; i < len(bannerDraw); i++ {
			if i <= 1 {
				l.Printf("%s", customColorize(bannerDraw[i], "\033[38;5;156m"))
			} else if i <= 3 {
				l.Printf("%s", customColorize(bannerDraw[i], "\033[38;5;155m"))
			} else {
				l.Printf("%s", customColorize(bannerDraw[i], "\033[38;5;046m"))
			}
		}
	}

	if color == "blue" {
		for i := 0; i < len(bannerDraw); i++ {
			if i <= 1 {
				l.Printf("%s", customColorize(bannerDraw[i], "\033[38;5;164m"))
			} else if i <= 3 {
				l.Printf("%s", customColorize(bannerDraw[i], "\033[38;5;091m"))
			} else {
				l.Printf("%s", customColorize(bannerDraw[i], "\033[38;5;054m"))
			}
		}
	}

	if color == "red" {
		for i := 0; i < len(bannerDraw); i++ {
			if i <= 1 {
				l.Printf("%s", customColorize(bannerDraw[i], "\033[38;5;196m"))
			} else if i <= 3 {
				l.Printf("%s", customColorize(bannerDraw[i], "\033[38;5;160m"))
			} else {
				l.Printf("%s", customColorize(bannerDraw[i], "\033[38;5;124m"))
			}
		}
	}
}
