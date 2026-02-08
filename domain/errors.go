package domain

import "fmt"

type FreqError struct {
	frequencyString string
	message         string
}

func (f *FreqError) Error() string {
	return fmt.Sprintf("Error converting %s: %s", f.frequencyString, f.message)
}

type CreateError struct {
	object  string
	message string
}

func (c *CreateError) Error() string {
	return fmt.Sprintf("Error creating %s: %s", c.object, c.message)
}
