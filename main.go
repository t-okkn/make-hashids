package main

import (
	"flag"
	"fmt"
	"strings"
)

const LISTEN_PORT string = ":8501"

var (
	Version  string
	Revision string
)

// <summary>: main関数
func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Error:", "引数が不足しています")
		ShowHelp()
		return
	}

	switch strings.ToLower(args[0]) {
	case "version":
		fmt.Println(Version, Revision)

	case "daemon":
		SetupRouter().Run(LISTEN_PORT)

	case "encode":
		if len(args) < 2 {
			fmt.Println("Error:", "文字列を入力してください")
			ShowHelp()
			return
		}

		for _, s := range args[1:] {
			out := fmt.Sprintf("source: %s, hashids:", s)
			fmt.Println(out, EncodeToHashids(s))
		}

	case "short":
		if len(args) < 2 {
			fmt.Println("Error:", "文字列を入力してください")
			ShowHelp()
			return
		}

		for _, s := range args[1:] {
			fmt.Println(s+":", GetShortHashids(s))
		}

	default:
		fmt.Println("Error:", "無効な引数が入力されました")
		ShowHelp()
	}
}

func ShowHelp() {
	var sb strings.Builder
	sb.Grow(500)

	sb.WriteString("\nHelp: Hashidsを生成します\n")
	sb.WriteString("Usage:\n")
	sb.WriteString("  make-hashids {encode | short | daemon | version} [...strings]\n\n")

	enc := fmt.Sprintf("%-12s文字列から可逆なHashidsを生成します\n", "  encode")
	// dec := fmt.Sprintf("%-12s可逆なHashidsを文字列にデコードします\n", "  decode")
	sht := fmt.Sprintf("%-12s文字列から不可逆なHashidsを生成します\n", "  short")
	dae := fmt.Sprintf("%-12s不可逆なHashidsを生成するWebサーバを起動します\n", "  daemon")
	ver := fmt.Sprintf("%-12sバージョン情報を表示します\n", "  version")

	sb.WriteString(enc)
	// sb.WriteString(dec)
	sb.WriteString(sht)
	sb.WriteString(dae)
	sb.WriteString(ver)

	fmt.Println(sb.String())
}
