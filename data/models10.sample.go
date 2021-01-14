// Package models is a test
package models

import "fmt"

// Comment here
func (u *User) getName() string {
	/** Return TODO: parse */
	return u.Name
}

/** Description
  multilines
*/
func show(name string, age int64) {
	// Line 1
	// Line 2
	fmt.Println(age, name)
}

func welcome() string {
	return "Welcome"
}

// Final comment
