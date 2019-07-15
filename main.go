package main

// import "nginx/unit"

const LISTEN_PORT string = ":8501"

var (
	Version string
	Revision string
)

// summary => main関数（サーバを開始します）
/////////////////////////////////////////
func main() {
	// unit.ListenAndServe(LISTEN_PORT, SetupRouter())
	SetupRouter().Run(LISTEN_PORT)
}

