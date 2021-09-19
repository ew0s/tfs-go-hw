package main

import (
	"fmt"
)

type (
	argumentType  string
	argumentValue string
)

type (
	figure [][]string
)

type color = argumentValue

type figureProperties map[argumentType]argumentValue
type structureAlgorithm func(*figure)
type figureModifier func(*figure, figureProperties)

const (
	Size  argumentType = "size"
	Char  argumentType = "char"
	Color argumentType = "color"
)

const (
	SizeDefaultValue  argumentValue = "15"
	CharDefaultValue  argumentValue = "X"
	ColorDefaultValue argumentValue = "\033[31m"
)

const (
	Red    color = "\033[31m"
	Green  color = "\033[32m"
	Yellow color = "\033[33m"
	Blue   color = "\033[34m"
)

func (v argumentValue) toInt() (int, bool) {
	var convertedValue int
	var isConvertible = true
	if _, err := fmt.Sscan(string(v), &convertedValue); err != nil {
		isConvertible = false
	}
	return convertedValue, isConvertible
}

func (v argumentValue) toRune() (rune, bool) {
	var isConvertible = true
	runes := []rune(v)
	if len(runes) == 0 || len(runes) > 1 {
		isConvertible = false
	}
	return runes[0], isConvertible
}

func (f *figure) updateElements(elem string) {
	for rowInd, figureRow := range *f {
		for columnIdx := range figureRow {
			(*f)[rowInd][columnIdx] = elem
		}
	}
}

func (f *figure) init() figureProperties {
	properties := figureProperties{}
	properties.fillWIthDefaultValues()
	setSize(properties[Size])(f, properties)
	setChar(properties[Char])(f, properties)
	setColor(properties[Color])(f, properties)
	return properties
}

func (f figure) print() {
	for _, figureRow := range f {
		for _, element := range figureRow {
			fmt.Printf("%s ", element)
		}
		fmt.Printf("\n")
	}
}

func (f figureProperties) fillWIthDefaultValues() {
	f[Size] = SizeDefaultValue
	f[Char] = CharDefaultValue
	f[Color] = ColorDefaultValue
}

func setSize(v argumentValue) figureModifier {
	return func(fig *figure, properties figureProperties) {
		if _, ok := v.toInt(); ok {
			properties[Size] = v
		}
		var size int
		var ok bool
		if size, ok = v.toInt(); !ok {
			size, _ = properties[Size].toInt()
		}
		*fig = make(figure, size)
		for idx := range *fig {
			(*fig)[idx] = make([]string, size)
		}
		newElement := properties[Color] + properties[Char]
		fig.updateElements(string(newElement))
	}
}

func setChar(v argumentValue) figureModifier {
	return func(fig *figure, properties figureProperties) {
		if _, ok := v.toRune(); ok {
			properties[Char] = v
		}
		var char rune
		var ok bool
		color := properties[Color]
		if char, ok = v.toRune(); !ok {
			char, _ = properties[Char].toRune()
		}
		newElement := string(color) + string(char)
		fig.updateElements(newElement)
	}
}

func setColor(v argumentValue) figureModifier {
	return func(fig *figure, properties figureProperties) {
		properties[Color] = v
		color := v
		char := properties[Char]
		newElement := string(color) + string(char)
		fig.updateElements(newElement)
	}
}

func hourglass() structureAlgorithm {
	return func(fig *figure) {
		if len(*fig) <= 2 {
			return
		}

		whitespaceChar := " "
		for rowInd, figureRow := range *fig {
			if rowInd == 0 || rowInd == len(*fig)-1 {
				continue
			}
			for columnInd := range figureRow {
				if rowInd == columnInd || (columnInd == (len(*fig) - 1 - rowInd)) {
					continue
				}
				(*fig)[rowInd][columnInd] = whitespaceChar
			}
		}
	}
}

func constructFigure(algo structureAlgorithm, modifiers ...figureModifier) figure {
	var figure figure
	properties := figure.init()
	for _, modifier := range modifiers {
		modifier(&figure, properties)
	}

	algo(&figure)
	figure.print()
	return figure
}

func main() {
	constructFigure(hourglass(), setSize("5"), setChar("@"), setColor(Red))
	constructFigure(hourglass(), setColor(Green), setSize("6"), setChar("&"))
	constructFigure(hourglass(), setChar("$"), setColor(Yellow), setSize("7"))
	constructFigure(hourglass(), setColor(Blue), setSize("8"), setChar("X"))
}
