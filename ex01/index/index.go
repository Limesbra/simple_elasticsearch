package index

import (
	"ex01/model"
	"fmt"
	"strings"
)

const body = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Places</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>

<body>
<h5>Total: %d</h5>
<h5>Current: %d</h5>
<ul>
%s
</ul>
%s
</body>
</html>`

func createButton(page int, name string) string {
	return fmt.Sprintf(`<a href="/?page=%d">%s</a>`, page, name)
}

func createPlaceEntry(place model.Place) string {
	return fmt.Sprintf(
		`
		<li>
			<div>%s</div>
			<div>%s</div>
			<div>%s</div>
		</li>`,
		place.Name,
		place.Address,
		place.Phone,
	)
}

func BuildHTML(total int, pageSize int, page int, places []model.Place) string {
	restaurants := make([]string, 10)
	if page != total {
		for i := 0; i < 10; i++ {
			restaurants[i] = createPlaceEntry(places[i])
		}
	} else {
		for i := 0; i < 9; i++ {
			restaurants[i] = createPlaceEntry(places[i])
		}
	}

	listOfPlaces := strings.Join(restaurants, "\n")

	var firstButton, prevButton, nextButton, lastButton string
	if page != 1 && total != 1 {
		firstButton = createButton(1, "First")
		prevButton = createButton(page-1, "Previous")
	}

	if page != total {
		nextButton = createButton(page+1, "Next")
		lastButton = createButton(total, "Last")
	}

	stringButtons := strings.Join(
		[]string{
			firstButton,
			prevButton,
			nextButton,
			lastButton,
		},
		"\n",
	)
	return fmt.Sprintf(body, total, page, listOfPlaces, stringButtons)
}
