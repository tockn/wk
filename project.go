package main

type ProjectStore interface {
	Current() (string, error)
	Set() (string, error)
}
