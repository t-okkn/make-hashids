package main

import "nginx/unit"

// summary => main関数（サーバを開始します）
/////////////////////////////////////////
func main() {
	unit.ListenAndServe(":8501", SetupRouter())
	// SetupRouter().Run()
}

