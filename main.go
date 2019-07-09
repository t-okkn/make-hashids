package main

// import "nginx/unit"

var (
	Version string
	Revision string
)

// summary => main関数（サーバを開始します）
/////////////////////////////////////////
func main() {
	// unit.ListenAndServe(":8080", SetupRouter())
	SetupRouter().Run()
}

