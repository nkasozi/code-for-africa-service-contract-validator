package main

func main() {
	debug_enabled := false
	app := NewApplication(debug_enabled)
	app.Run()
}
