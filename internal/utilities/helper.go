package utilities

import (
	"os"
	"path/filepath"
	"strings"
)

func CheckForFileExtension(name string, extension string) bool {
	return strings.HasSuffix(name, extension)
}

func DeleteAll(path string) error {
	d, err := os.Open(path)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckIfLangIsSupported(lang string) bool {
	// ckeck if lang is supported
	supported := []string{"go", "cpp", "csharp", "java", "python", "ruby", "pyi", "php", "objc", "kotlin", "node", "dart"}
	for _, l := range supported {
		return l == lang
	}
	return false
}

// compare two strings array and return the difference
func Difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []string{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// compare two strings array and return the intersection
func Intersection(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []string{}
	for _, x := range a {
		if _, ok := mb[x]; ok {
			ab = append(ab, x)
		}
	}
	return ab
}

// compare two strings array and return the union
func Union(a, b []string) []string {
	return append(a, Difference(b, a)...)
}

// compare two strings array and return the symmetric difference
func SymmetricDifference(a, b []string) []string {
	return append(Difference(a, b), Difference(b, a)...)
}

// compare two strings array and return the relative complement of a in b
func RelativeComplement(a, b []string) []string {
	return Difference(a, Intersection(a, b))
}

// check if a string array contains a string
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
