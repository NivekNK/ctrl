package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	commands := []string{
		"pwsh",
		"-c",
		"$wingetApp = \"APP_ID\";",
		"$wingetList = winget list $wingetApp;",
		"if ($wingetList -match \"No installed package found matching input criteria.\") { \"\" } else {",
		"    $array = $wingetList[$wingetList.Length - 1] -split \"\\s+\";",
		"    [array]::Reverse($array);",
		"    $count = ($wingetList[$wingetList.Length - 3] -split \"\\s+\").Length;",
		"    if ($count -ge 5) { $array[2] + \",\" + $array[1] } else { $array[1] + \",_\" }",
		"}",
		// "$output = winget list $package | Select-Object -First 1;",
		// "$output",
		// "if ($output -match \"No installed package found matching input criteria.\") { \"\" } else {",
		// "    $array = ($output -split \"\\s+\" | [array]::Reverse($_); $_)",
		// "    if ($array.Length -ge 3) { $array[2] } else { \"\" }",
		// "}",
	}

	commandName := commands[0]
	commands = commands[1:]

	app_id := "GitHub.cli"
	for i := len(commands) - 1; i >= 0; i-- {
		if strings.Contains(commands[i], "APP_ID") {
			commands[i] = strings.Replace(commands[i], "APP_ID", app_id, 1)
			break
		}
	}

	command := exec.Command(commandName, commands...)
	stdout, err := command.Output()
	fmt.Println(string(stdout))
	if err != nil {
		fmt.Println(err)
	}
}

// import (
// 	"ctrl/cmd"
// 	"ctrl/util"
// 	_ "embed"
// 	"fmt"
// )

// //go:embed .sql/schema.sql
// var schema string

// func main() {
// 	cmd.Execute()
// }

// func init() {
// 	util.Schema = &schema

// 	instance, err := util.InitializeInstance()
// 	if err != nil {
// 		panic(err)
// 	}

// 	config, err := util.LoadConfig()
// 	if err != nil {
// 		panic(err)
// 	}
// 	cmd.Config = config

// 	data, err := util.LoadData()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	data.Refresh(config, instance)
// }
