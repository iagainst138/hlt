package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    //"os/signal"
    "regexp"
    "strings"
    "time"
)


const (
    RESET   = "\x1b[0m"
    BRIGHT  = "\x1b[1m"
)


var colors map[string]string = map[string]string {
    "bright":   BRIGHT,
    "blue":     "\x1b[34m",
    "cyan":     "\x1b[36m",
    "green":    "\x1b[32m",
    "magenta":  "\x1b[35m",
    "red":      "\x1b[31m",
    "yellow":   "\x1b[33m",
}


type matches []string

func (m *matches) String() string {
    return ""
}

func (m *matches) Set(value string) error {
    *m = append(*m, value)
    return nil
}


func highlight(in io.Reader, m matches, color string, logchan chan string, ignorecase, bright bool) {
    if bright {
        color = BRIGHT + color
    }
    pattern := ""
    for i,p := range m {
        pattern += p
        if((i+1) < len(m)) {
            pattern += "|"
        }
    }
    if ignorecase {
        pattern = "(?i)" + pattern
    }
    reader := bufio.NewReader(in)
    for {
        // TODO validate this
        _,err := reader.Peek(1)
        for err != nil {
            time.Sleep(500 * time.Millisecond)
            _,err = reader.Peek(1)
        }
        if text,err := reader.ReadString('\n'); err != nil {
            // TODO throw error here
            fmt.Println("Error:", err)
            return
        } else {
            text = strings.TrimSpace(text)
            re := regexp.MustCompile(pattern)
            if(len(re.FindString(text)) > 0) {
                fmt.Printf("%s%s%s\n", color, text, RESET)
            } else {
                fmt.Println(text)
            }
            if logchan != nil {
                logchan<-text
            }
        }
    }
}


func logtofile(logfile string, logchan chan string, doappend bool) {
    // TODO handle error
    file_opts := os.O_RDWR|os.O_CREATE
    if doappend {
        file_opts = file_opts|os.O_APPEND
    } else {
        file_opts = file_opts|os.O_TRUNC
    }
    l,e := os.OpenFile(logfile, file_opts, 0644)
    if e != nil {
        fmt.Println("Error:", e)
        os.Exit(1)
    }
    logger := log.New(l, "", log.Ldate|log.Ltime)
    logger.Print("--- started ---")
    for {
        logger.Print(<-logchan)
    }
}


func show_help() {
    // TODO improve this 
    fmt.Println("Error: Color does not exist.")
    fmt.Printf("Valid colors:")
    i := 1
    for c := range colors {
        fmt.Printf(" %s", c)
        if(i < len(colors)) {
            fmt.Printf(",")
        } else {
            fmt.Printf(".")
        }
        i++
    }
    fmt.Println()
}


func main() {
    var bright bool
    var color string
    var doappend bool
    var ignorecase bool
    var logchan = make(chan string)
    var logfile string
    var m matches

    flag.BoolVar(&doappend, "a", false, "Append tp log specified by -l.")
    flag.BoolVar(&bright, "b", false, "Make output bright.")
    flag.StringVar(&color, "c", "bright", "Change color of matched output.")
    flag.BoolVar(&ignorecase, "i", false, "Ignore case.")
    flag.StringVar(&logfile, "l", "", "A regular expression to find.")
    flag.Var(&m, "m", "A regular expression to find.")
    flag.Parse()

    if logfile != "" {
        go logtofile(logfile, logchan, doappend)
        // TODO overkill?
        /*c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        go func(){
            for sig := range c {
                // TODO is the if needed
                if sig.String() == "interrupt" {
                    logchan<-"--- finished ---"
                    os.Exit(0)
                }
            }
        }()*/
    } else {
        logchan = nil
    }

    if _,ok := colors[color]; ok == false {
        show_help()
        os.Exit(1)
    }

    if len(m) == 0 {
        fmt.Println("Error no regular expressions specified.")
        os.Exit(1)
    }

    stat,_ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) == 0 {
        highlight(os.Stdin, m, colors[color], logchan, ignorecase, bright)
    } else {
        if flag.NArg() == 0 {
            fmt.Println("Error: no input file specified.")
            os.Exit(1)
        } else {
            // NOTE only handles one file for now
            f,err := os.Open(flag.Arg(0))
            if err != nil {
                fmt.Println("Error opening file:", err)
                os.Exit(2)
            }
            highlight(f, m, colors[color], logchan, ignorecase, bright)
        }
    }
}
