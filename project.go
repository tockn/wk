package main

type ProjectStore interface {
	Current() string
	Set() string
}
