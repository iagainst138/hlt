package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "regexp"
    "strings"
    )

type matches []string

func (m *matches) String() string {
    return "A slice of regular expressions"
}

func (m *matches) Set(value string) error {
    *m = append(*m, value)
    return nil
}

func highlight(m matches) {
    reader := bufio.NewReader(os.Stdin)
    for {
        if text, e := reader.ReadString('\n'); e != nil {
            // TODO throw error here
            fmt.Println("Error:", e)
        } else {
            text = strings.TrimSpace(text)
            match := ""
            for _,r := range m {
                re := regexp.MustCompile(r)
                if match = re.FindString(text); len(match) > 0 {
                    break
                }
            }
            if(len(match) > 0) {
                fmt.Printf("\x1b[1m%s\x1b[0m\n", text)
            } else {
                fmt.Println(text)
            }
        }
    }
}

func main() {
    var m matches
    flag.Var(&m, "m", "A regular expression to find.")
    flag.Parse()
    stat, _ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) == 0 {
        highlight(m)
    } else {
        // TODO raise proper error
        fmt.Println("Error: stdin is from a terminal")
    }
}
