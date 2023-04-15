package config

import (
	"fmt"
	"regexp"

	"gn/style"
	"gn/style/color"

	"github.com/charmbracelet/lipgloss"
)

type Colors struct {
	Blurred string
	Border  string
	Focused string
}

func (c *Colors) setColors() error {
	err := c.CheckValidity()
	if err != nil {
		return err
	}

	if len(c.Blurred) > 0 {
		color.Blurred = lipgloss.Color(c.Blurred)
	}

	if len(c.Border) > 0 {
		color.Border = lipgloss.Color(c.Border)
	}

	if len(c.Focused) > 0 {
		color.Focused = lipgloss.Color(c.Focused)
	}

	style.Update()

	return nil
}

func (c *Colors) CheckValidity() error {
	if len(c.Blurred) > 0 {
		err := checkColor(c.Blurred)
		if err != nil {
			return fmt.Errorf("color 'Blurred' (%s) is invalid: %w", c.Blurred, err)
		}
	}

	if len(c.Border) > 0 {
		err := checkColor(c.Border)
		if err != nil {
			return fmt.Errorf("color 'Border' (%s) is invalid: %w", c.Border, err)
		}
	}

	if len(c.Focused) > 0 {
		err := checkColor(c.Focused)
		if err != nil {
			return fmt.Errorf("color 'Focused' (%s) is invalid: %w", c.Focused, err)
		}
	}

	return nil
}

func checkColor(colorStr string) error {
	if len(colorStr) != 7 {
		return fmt.Errorf("invalid length. Expected 7, got %d", len(colorStr))
	}

	r := regexp.MustCompile("#[0-9a-fA-F]{6}")
	if !r.MatchString(colorStr) {
		return fmt.Errorf("color contains invalid chars")
	}

	return nil
}
