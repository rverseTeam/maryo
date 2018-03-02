/*

maryo/maryo.go

written by superwhiskers, licensed under gnu gplv3.
if you want a copy, go to http://www.gnu.org/licenses/

*/

package main

import (
  // internals
  "fmt"
  "runtime"
  "flag"
  "time"
  "os"
  "encoding/json"
  "strings"
  // externals
  "github.com/shiena/ansicolor"
)

func setup(fileMap map[string]string) {

  // test data
  test := make([]string, 1)
  test[0] = "account"

  // test the 'official' pretendo servers
  testOfficial := make([]string, 1)
  testOfficial[0] = "account"

  // file status map
  fileStat := make(map[string]string)
  fileStat["ne"] = "nonexistent"
  fileStat["iv"] = "invalid"
  fileStat["va"] = "valid"
  fileStat["uk"] = "unknown"

  // setup environment
  clear()
  ttitle("maryo -> setup")

  // show setup screen
  fmt.Printf("== maryo -> setup ===========================================\n")
  fmt.Printf("                                        Steps:               \n")
  fmt.Printf(" welcome to the maryo setup wizard.     > intro              \n")
  fmt.Printf(" this program will walk you through       config creation    \n")
  fmt.Printf(" setting up your very own Pretendo        confirm prefs      \n")
  fmt.Printf(" proxy server for accessing the server.   display proxy info \n")
  fmt.Printf(" -> press enter                                              \n")
  fmt.Printf("                                                             \n")
  fmt.Printf("                                                             \n")
  fmt.Printf("=============================================================\n")
  input("")

  // show config creation screen
  var method string
  for true {
    clear()
    fmt.Printf("== maryo -> setup ===========================================\n")
    fmt.Printf("                                        Steps:               \n")
    fmt.Printf(" how would you like to configure the      intro              \n")
    fmt.Printf(" proxy?                                 > config creation    \n")
    fmt.Printf(" 1. automatic                             confirm prefs      \n")
    fmt.Printf(" 2. custom                                display proxy info \n")
    fmt.Printf(" 3. template                                                 \n")
    fmt.Printf("                                                             \n")
    fmt.Printf(" -> (1|2|3)                                                  \n")
    fmt.Printf("=============================================================\n")
    method = input(": ")

    if ( method == "1" ) || ( method == "2" ) || ( method == "3" ) {
      break
    } else {
      fmt.Printf("-> please enter 1, 2, or 3\n")
      time.Sleep(1500 * time.Millisecond)
    }
  }

  // show log when
  clear()
  fmt.Printf("== maryo -> setup ===========================================\n")
  fmt.Printf("                                                             \n")
  fmt.Printf(" configuring proxy..                                         \n")
  fmt.Printf(" current config status: %s\n", fileStat[fileMap["config"]])
  if method == "1" {
    fmt.Printf(" method: automatic..\n")
    fmt.Printf("-- beginning tests\n")
    fmt.Printf(" 1. attempting to detect endpoints running on this machine\n")

    // test for endpoints on this machine
    result := make([]bool, len(test))
    for x := 0; x < len(test); x++ {

      // test the endpoint
      fmt.Printf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", test[x]), endpointsFor("local", test[x]))

      // get the json
      res, err := get(strings.Join([]string { endpointsFor("local", test[x]), "/isthisworking" }, ""))

      // prepare a struct for the json
      var parsedRes isitworkingStruct

      // parse it
      err2 := json.Unmarshal([]byte(res), &parsedRes)
      if res != "" {
        if err2 != nil {
          if err == nil {
            fmt.Printf("\n[err] : error when parsing JSON during validating %s server\n", test[x])
            panic(err2)
          }
        }
      }

      // handle the results
      if (parsedRes.Server == serverResFor(test[x])) && (err == nil) && (res != "") {
        if isWindows() {
          Writer := ansicolor.NewAnsiColorWriter(os.Stdout)
          fmt.Fprintf(Writer, "%s\n", padStrToMatchStr(fmt.Sprintf("\r  %s%s%s%s %s -> %s", code("green"), code("bold"), utilIcons("success"), code("reset"), endpointsFor("ninty", test[x]), endpointsFor("local", test[x])), fmt.Sprintf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", test[x]), endpointsFor("local", test[x])), " "))
        } else {
          fmt.Printf("%s\n", padStrToMatchStr(fmt.Sprintf("\r  %s%s%s%s %s -> %s", code("green"), code("bold"), utilIcons("success"), code("reset"), endpointsFor("ninty", test[x]), endpointsFor("local", test[x])), fmt.Sprintf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", test[x]), endpointsFor("local", test[x])), " "))
        }
        result[x] = true
      } else {
        if isWindows() {
          Writer := ansicolor.NewAnsiColorWriter(os.Stdout)
          fmt.Fprintf(Writer, "%s\n", padStrToMatchStr(fmt.Sprintf("\r  %s%s%s%s %s -> %s", code("red"), code("bold"), utilIcons("failiure"), code("reset"), endpointsFor("ninty", test[x]), endpointsFor("local", test[x])), fmt.Sprintf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", test[x]), endpointsFor("local", test[x])), " "))
        } else {
          fmt.Printf("%s\n", padStrToMatchStr(fmt.Sprintf("\r  %s%s%s%s %s -> %s", code("red"), code("bold"), utilIcons("failiure"), code("reset"), endpointsFor("ninty", test[x]), endpointsFor("local", test[x])), fmt.Sprintf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", test[x]), endpointsFor("local", test[x])), " "))
        }
        result[x] = false
      }

    }

    // show the user that we are detecting the official server
    fmt.Printf(" 2. attempting to test endpoints on the official server\n")

    // test for endpoints on the official pretendo servers
    resultOfficial := make([]bool, len(testOfficial))
    for x := 0; x < len(testOfficial); x++ {

      // test the endpoint
      fmt.Printf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", testOfficial[x]), endpointsFor("official", testOfficial[x]))

      // get the json
      res2, err3 := get(strings.Join([]string { endpointsFor("official", testOfficial[x]), "/isthisworking" }, ""))
      // prepare a struct for the json
      var parsedRes2 isitworkingStruct

      // parse it
      err4 := json.Unmarshal([]byte(res2), &parsedRes2)
      if res2 != "" {
        if err4 != nil {
          if err3 == nil {
            fmt.Printf("\n[err] : error when parsing JSON during validating %s server\n", testOfficial[x])
            panic(err4)
          }
        }
      }

      // handle the results
      if (parsedRes2.Server == serverResFor(testOfficial[x])) && (err3 == nil) && (res2 != "") {
        if isWindows() {
          Writer := ansicolor.NewAnsiColorWriter(os.Stdout)
          fmt.Fprintf(Writer, "%s\n", padStrToMatchStr(fmt.Sprintf("\r  %s%s%s%s %s -> %s", code("green"), code("bold"), utilIcons("success"), code("reset"), endpointsFor("ninty", testOfficial[x]), endpointsFor("official", testOfficial[x])), fmt.Sprintf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", testOfficial[x]), endpointsFor("official", testOfficial[x])), " "))
        } else {
          fmt.Printf("%s\n", padStrToMatchStr(fmt.Sprintf("\r  %s%s%s%s %s -> %s", code("green"), code("bold"), utilIcons("success"), code("reset"), endpointsFor("ninty", testOfficial[x]), endpointsFor("official", testOfficial[x])), fmt.Sprintf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", testOfficial[x]), endpointsFor("official", testOfficial[x])), " "))
        }
        resultOfficial[x] = true
      } else {
        if isWindows() {
          Writer := ansicolor.NewAnsiColorWriter(os.Stdout)
          fmt.Fprintf(Writer, "%s\n", padStrToMatchStr(fmt.Sprintf("\r  %s%s%s%s %s -> %s", code("red"), code("bold"), utilIcons("failiure"), code("reset"), endpointsFor("ninty", testOfficial[x]), endpointsFor("official", testOfficial[x])), fmt.Sprintf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", testOfficial[x]), endpointsFor("official", testOfficial[x])), " "))
        } else {
          fmt.Printf("%s\n", padStrToMatchStr(fmt.Sprintf("\r  %s%s%s%s %s -> %s", code("red"), code("bold"), utilIcons("failiure"), code("reset"), endpointsFor("ninty", testOfficial[x]), endpointsFor("official", testOfficial[x])), fmt.Sprintf("  %s %s -> %s", utilIcons("uncertain"), endpointsFor("ninty", testOfficial[x]), endpointsFor("official", testOfficial[x])), " "))
        }
        resultOfficial[x] = false
      }

    }

    // print out the results
    fmt.Printf("-- printing results of tests\n")
    for x := 0; x < len(result); x++ {

      // add a header saying that these are local results
      fmt.Printf("- local\n")

      // print the results
      if result[x] == true {
        fmt.Printf(" %s: success\n", test[x])
      } else {
        fmt.Printf(" %s: failiure\n", test[x])
      }

    }
    for x := 0; x < len(resultOfficial); x++ {

      // add a header saying that these are official results
      fmt.Printf("- pretendo\n")

      // print the results
      if resultOfficial[x] == true {
        fmt.Printf(" %s: success\n", testOfficial[x])
      } else {
        fmt.Printf(" %s: failiure\n", testOfficial[x])
      }

    }

    // begin creating the config
    fmt.Printf("-- creating config file\n")

    // create cfgTest, cfgResult, and using variables
    var using     string
    var cfgTest   []string
    var cfgResult []bool

    // local servers have priority
    if len(test) != 0 {
      using     = "local"
      cfgTest   = test
      cfgResult = result
    } else {
      using     = "official"
      cfgTest   = testOfficial
      cfgResult = resultOfficial
    }

    // make a map for the config
    config := make(map[string]string)

    // apply a nice helping of all of the working endpoints to the config
    for x := 0; x < len(cfgTest); x++ {
      if cfgResult[x] == true {
        config[endpointsFor("ninty", cfgTest[x])] = endpointsFor(using, cfgTest[x])
      }
    }

    // idk what to do after this
    stringifiedConfig, err := json.Marshal(config)
    if err != nil {
      fmt.Printf("[err] : error when stringifying json")
    }

    // place it into the file
    if fileMap["config"] == "iv" {
      deleteFile("config.json")
      createFile("config.json")
      writeByteToFile("config.json", stringifiedConfig)
    } else if fileMap["config"] == "ne" {
      createFile("config.json")
      writeByteToFile("config.json", stringifiedConfig)
    } else if fileMap["config"] == "uk" {
      // detect status of config and do the
      // things to write to it.
      if doesFileExist("config.json") == true {
        deleteFile("config.json")
        createFile("config.json")
        writeByteToFile("config.json", stringifiedConfig)
      } else {
        createFile("config.json")
        writeByteToFile("config.json", stringifiedConfig)
      }
    }

  } else if method == "2" {
    fmt.Printf(" method: custom..\n")
  } else if method == "3" {
    fmt.Printf(" method: stock..\n")
  }

}

func main() {

  // parse some flags here
  config := flag.String("config", "config.json", "value for config file path")
  //logging := flag.Bool("logging", false, "if set, the proxy will log all requests and data")
  doSetup := flag.Bool("setup", false, "if set, maryo will go through setup again")
  flag.Parse()

  // set window title
  ttitle("maryo")

  if *doSetup == false {

    clear()
    fmt.Printf("-- log\n")

    // startup
    fmt.Printf("= startup")

    // i might include some os-specific code
    fmt.Printf("loading 0%%: detecting os")
    fmt.Printf("%s\n", padStrToMatchStr(fmt.Sprintf("\ros: %s", runtime.GOOS), "loading 0%%: detecting os", " "))

    // file checking
    fmt.Printf("loading 50%%: checking files.")

    // map for holding file status
    fileMap := make(map[string]string)

    // config.json -- if nonexistent, it follows the user's instruction to create one, or use a builtin copy
    fileMap["config"] = "ne"
    if doesFileExist(*config) == false {
      fmt.Printf("%s\n", padStrToMatchStr("\rconfig: nonexistent", "loading 50%%: checking files.", " "))
    } else {
      // check if the file is valid JSON
      if checkJSONValidity(*config) != true {
        fmt.Printf("%s\n", padStrToMatchStr("\rconfig: invalid", "loading 50%%: checking files.", " "))
        fileMap["config"] = "iv"
      } else {
        fmt.Printf("%s\n", padStrToMatchStr("\rconfig: valid", "loading 50%%: checking files.", " "))
        fileMap["config"] = "va"
      }
    }

    // print final info
    fmt.Printf("loaded..\n")

    // do the setup function if the file isn't completely correct
    if fileMap["config"] == "ne" {
      setup(fileMap)
    } else if fileMap["config"] == "iv" {
      fmt.Printf("your config is invalid.\n")
      fmt.Printf("you have three different options:\n")
      fmt.Printf(" 1. run this program with the --setup flag\n")
      fmt.Printf(" 2. delete the config and run this program\n")
      fmt.Printf(" 3. fix the config\n")
      os.Exit(1)
    } else {
      startProxy()
    }

  // run setup function and load proxy
  } else {

    // fileMap
    var fileMap map[string]string
    fileMap = make(map[string]string)

    // place config value in there
    fileMap["config"] = "uk"

    // just do the setup function
    setup(fileMap)

  }
}
