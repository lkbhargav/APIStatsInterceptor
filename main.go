package main

import (
	"APIStatsInterceptor/types"
	"APIStatsInterceptor/util"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/lkbhargav/requests"
)

func main() {
	var url string
	flag.StringVar(&url, "url", "accounts.bgalytics.com", "a string var")

	var path string
	flag.StringVar(&path, "path", "", "Status^status,something2,anything3^COMMA|Nodes^number_of_nodes^P<any text to be prefixed>")

	var headersStr string
	flag.StringVar(&headersStr, "headers", "", "Authorization:abdhkxxx,Origin:bgalytics.com")

	var frequency *int
	frequency = flag.Int("freq", 1000, "time in milliseconds")

	var jsonBody string
	flag.StringVar(&jsonBody, "jsonBody", "", "Eg. {\"name\":\"Jameson\"}")

	var requestMethod string
	flag.StringVar(&requestMethod, "requestMethod", "GET", "Supported values => GET,POST,PUT")

	var color string
	flag.StringVar(&color, "color", "black", "Supported values => white, black, yellow, red, blue, green")

	flag.Parse()

	if !strings.Contains(url, "://") {
		url = "http://" + url
	}

	sets, err := util.ParseSets(path)

	if err != nil {
		fmt.Println("Invalid path passed, please check the format and retry again | Format: Status^status^COMMA|Nodes^number_of_nodes^PERCENT")
		fmt.Println(err)
		os.Exit(1)
	}

	var reqBody map[string]interface{}

	if jsonBody != "" {
		json.Unmarshal([]byte(jsonBody), &reqBody)
	}

	headers := make(map[string]string)

	if headersStr != "" {
		headers, err = util.ParseHeaders(headersStr)

		if err != nil {
			fmt.Println("Invalid headers passed, please check the format and retry again | Format: Authorization:abdhkxxx,Origin:bgalytics.com")
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if len(sets) == 0 {
		fmt.Println("Not enough path variables to continue. Try again later")
		os.Exit(1)
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()

	uiEvents := ui.PollEvents()

	var txt string

	go func() { // running in a goroutine to make the keyboard events bind and work
		for {
			var req requests.Request

			req.URL = url
			req.Headers = headers

			if requestMethod == "POST" {
				req.Method = requests.POST

				if len(reqBody) > 0 {
					req.JSONBody = reqBody
				}
			} else if requestMethod == "PUT" {
				req.Method = requests.PUT

				if len(reqBody) > 0 {
					req.JSONBody = reqBody
				}
			}

			resp := req.Do()

			txt = ""

			if resp.Error != nil {
				txt = fmt.Sprintf("Error: %v | Please check the URL and try again later.", err)
			}

			var tmp []string

			p.SetRect(0, 0, 1000, len(sets)+1*100)

			p.Border = false

			switch color {
			case "white":
				p.TextStyle.Fg = ui.ColorWhite
			case "yellow":
				p.TextStyle.Fg = ui.ColorYellow
			case "blue":
				p.TextStyle.Fg = ui.ColorBlue
			case "red":
				p.TextStyle.Fg = ui.ColorRed
			case "green":
				p.TextStyle.Fg = ui.ColorGreen
			default:
				p.TextStyle.Fg = ui.ColorBlack
			}

			if txt == "" {
				txt = fmt.Sprintf("Date: %v | URL: %v | Frequency: %vs\n\n", time.Now().Format(time.RFC1123), url, float64(*frequency)/1000)

				for _, path := range sets {
					tmp = []string{"response"}
					tmp = append(tmp, path.Path...)

					switch path.Option {
					case types.Comma:
						txt = txt + fmt.Sprintf("%v: %v\n\n", path.Name, humanize.Comma(int64(util.GetValNestedMap(resp.Response, tmp).(float64))))
					case types.Percent:
						txt = txt + fmt.Sprintf("%v: %v%%\n\n", path.Name, util.GetValNestedMap(resp.Response, tmp))
					case types.Data:
						txt = txt + fmt.Sprintf("%v: %v\n\n", path.Name, humanize.Bytes(uint64(util.GetValNestedMap(resp.Response, tmp).(float64))))
					case types.Prefix:
						txt = txt + fmt.Sprintf("%v: %v %v\n\n", path.Name, path.OptionalVal, util.GetValNestedMap(resp.Response, tmp))
					case types.Suffix:
						txt = txt + fmt.Sprintf("%v: %v %v\n\n", path.Name, util.GetValNestedMap(resp.Response, tmp), path.OptionalVal)
					default:
						txt = txt + fmt.Sprintf("%v: %v\n\n", path.Name, util.GetValNestedMap(resp.Response, tmp))
					}
				}
			}

			p.Text = txt

			ui.Render(p)

			time.Sleep(time.Duration(*frequency) * time.Millisecond)
		}
	}()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
