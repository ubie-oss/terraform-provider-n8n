package controllers

import "fmt"

func errDeleteProtected(kind, id string) error {
	return fmt.Errorf("cannot delete %s %s: delete protection is enabled", kind, id)
}
