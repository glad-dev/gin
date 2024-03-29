package configuration

import (
	"fmt"
	"regexp"

	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/style/color"

	"github.com/charmbracelet/lipgloss"
)

// Colors contains hex strings of colors for the different styles.
type Colors struct {
	Blurred string
	Border  string
	Focused string
}

// UpdateColors sets the colors in the configuration to the provided values and writes the updated configuration back.
func UpdateColors(c Colors, config *Config) error {
	err := c.CheckValidity()
	if err != nil {
		return err
	}

	config.Colors = c

	return write(config)
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

	style.UpdateColors()

	return nil
}

// CheckValidity checks if the colors in the struct are valid hex strings. Empty strings are ignored.
func (c *Colors) CheckValidity() error {
	if len(c.Blurred) > 0 {
		err := checkColor(c.Blurred)
		if err != nil {
			log.Error("color 'Blurred' is invalid", "color", c.Blurred, "error", err)

			return fmt.Errorf("color 'Blurred' (%s) is invalid: %w", c.Blurred, err)
		}
	}

	if len(c.Border) > 0 {
		err := checkColor(c.Border)
		if err != nil {
			log.Error("color 'Border' is invalid", "color", c.Border, "error", err)

			return fmt.Errorf("color 'Border' (%s) is invalid: %w", c.Border, err)
		}
	}

	if len(c.Focused) > 0 {
		err := checkColor(c.Focused)
		if err != nil {
			log.Error("color 'Focused' is invalid", "color", c.Focused, "error", err)

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
		return fmt.Errorf("color contains invalid chars. Expected hex code, got %s", colorStr)
	}

	return nil
}
