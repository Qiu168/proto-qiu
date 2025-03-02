package main

type generator interface {
	generate() error
}
