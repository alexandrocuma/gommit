package helpers

import (
	"fmt"
	"strconv"
)

// Helper functions
func ValidateTemperature(input string) error {
	val, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return fmt.Errorf("invalid number")
	}
	if val < 0.0 || val > 1.0 {
		return fmt.Errorf("temperature must be between 0.0 and 1.0")
	}
	return nil
}
