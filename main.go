package main

func main() {
	// http.ListenAndServe(":8080", SetupRouter())
	SetupRouter().Run()
}

