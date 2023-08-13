package e

import "fmt"

// Function for wrapping error and leading to the form: 'message: error'
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
// Finction for proccessing error and wrapping its if err isn`t nil
func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}