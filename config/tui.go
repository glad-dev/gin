package config

// Bubbletea called the Update function multiple times. Thus, we need to prevent the functions called in the Update
// function to be run multiple times.
// Since we only call one config function per program run, we can safely use a pair of variables for all calls

var (
	firstCall   = true
	errPrevious error
)
