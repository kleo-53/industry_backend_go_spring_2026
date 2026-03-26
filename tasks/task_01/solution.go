package main

func greet(name string) string {
	switch name {
	case "":
		return "Hello, World!"
	default:
		return "Hello, " + name + "!"
	}
}
