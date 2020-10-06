package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/wf-yamaday/cdiscord/discord"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cdiscord = kingpin.New("cdiscord", "notify to discord cli tool")
	version  = cdiscord.Command("version", "display version")

	message    = cdiscord.Command("message", "send message to discord")
	username   = message.Flag("username", "username").Short('u').String()
	alertLevel = message.Flag("level", "set alert level ('info' 'danger' 'health' 'warn')").Short('l').String()
	webhookURL = message.Flag("webhook-url", "webhook url").Required().Short('w').String()

	cdiscordVer = "v0.0.1"
)

type AlertLevel int

const (
	Danger AlertLevel = iota
	Warn
	Health
	Info
)

var colors = map[AlertLevel]string{
	Danger: "16527151",
	Warn:   "16763924",
	Health: "2611313",
	Info:   "2865650",
}

func main() {

	if _, err := cdiscord.Parse(os.Args[1:]); err != nil {
		cdiscord.FatalUsage(err.Error())
	}

	switch kingpin.MustParse(cdiscord.Parse(os.Args[1:])) {
	case version.FullCommand():
		fmt.Printf("version: %s\n", cdiscordVer)
		return
	default:
		var color string

		switch *alertLevel {
		case "info":
			color = colors[Info]
		case "warn":
			color = colors[Warn]
		case "danger":
			color = colors[Danger]
		case "health":
			color = colors[Health]
		default:
			color = colors[Info]
		}

		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("[error] %s", err)
			return
		}

		user := "cdiscord"
		if *username != "" {
			user = *username
		}

		param := discord.WebhookParam{
			Usename: user,
			Embeds: []discord.Embeds{
				{
					Description: string(in),
					Color:       color,
				},
			},
		}

		if ret, err := discord.SendMessage(param, *webhookURL); err != nil {
			log.Printf("[error] %s", err)
			log.Printf("[error] %s", string(ret))
		}
	}
}
