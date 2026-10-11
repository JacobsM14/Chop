package render

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"chop/internal/store"

	"github.com/charmbracelet/lipgloss"
)

type Color struct {
	Name string
	Code string
}

var Colors = []Color{
	{"red", "9"}, {"green", "10"}, {"blue", "12"},
	{"purple", "13"}, {"orange", "208"}, {"white", "15"},
}

func Names() []string {
	names := make([]string, len(Colors))
	for i, c := range Colors {
		names[i] = c.Name
	}
	return names
}

func Valid(name string) bool {
	for _, c := range Colors {
		if c.Name == name {
			return true
		}
	}
	return false
}

func style(color string) lipgloss.Style {
	for _, c := range Colors {
		if c.Name == color {
			return lipgloss.NewStyle().Foreground(lipgloss.Color(c.Code))
		}
	}
	return lipgloss.NewStyle()
}

// Paint colorea un texto con el color indicado (sin color si no existe).
func Paint(color, text string) string {
	return style(color).Render(text)
}

// Title devuelve el nombre de la chop en negrita y subrayado con una línea.
func Title(name, color string) string {
	line := strings.Repeat("─", utf8.RuneCountInString(name))
	return style(color).Bold(true).Render(name) + "\n" + style(color).Render(line)
}

// Item formatea un item como "- comando: descripción" (o "1. comando: descripción" si es ordenada).
func Item(i int, ordered bool, it store.ListItem) string {
	bullet := "-"
	if ordered {
		bullet = fmt.Sprintf("%d.", i+1)
	}
	line := bullet + " " + style(it.TextColor).Bold(true).Render(it.Text)
	if it.Description != "" {
		line += ": " + Paint(it.DescriptionColor, it.Description)
	}
	return line
}

// Command formatea un item como "$ comando", tal y como se va a ejecutar.
func Command(it store.ListItem) string {
	return "$ " + style(it.TextColor).Bold(true).Render(it.Text)
}
