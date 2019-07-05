package main

// summary => main関数（サーバを開始します）
/////////////////////////////////////////
func main() {
	// http.ListenAndServe(":8080", SetupRouter())
	SetupRouter().Run()
}

