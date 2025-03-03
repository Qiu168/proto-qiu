// Package generator is generate language files
package generator

type Generator interface {
	Generate() error
}
