package main

import (
	"flag"
	"fmt"
	"gostream/cmd"
)

func main() {
	// We define the args of the CLI
	var logName = flag.String("logs", "", "AWS Lambda Name")

	// TODO : Improve the error handling message

	// flag.Usage = func() {
	// 	fmt.Fprintf(os.Stderr, `Do something with Thing one and Thing two
	// 	Usage of:
	// 		[flags] [positional arg1] [positional arg2]

	// 	Positional arg1 must acquiesce the adjunct.  Positional arg2, if provided,
	// 	must reticulate the splines.
	// 	Flags:
	// 	`)
	// }

	flag.Parse()

	// if len(*logName) == 0 {
	// 	panic("Set a Lambda name")
	// }

	if len(*logName) > 0 {
		fmt.Println(*logName)
		cmd.BeginLogStream(*logName)
	}

	// switch os.Args[1] {
	// case "logs":
	// 	cmd.BeginLogStream(*logName)
	// default:
	// 	fmt.Printf("Unknown command %s\n", os.Args[1])
	// }
}
