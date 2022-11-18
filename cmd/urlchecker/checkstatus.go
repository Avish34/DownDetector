package urlchecker

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-ping/ping"
	"github.com/gookit/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CheckStatus cmd
var checkStatusCmd = &cobra.Command{
	Use:   "check-status",
	Short: "check-status",
	Long:  `check-status`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("no arguments passed")
		} else if len(args) > 1 {
			log.Fatal("Only 1 URL is expected")
		} else {
			if _, err := url.ParseRequestURI(args[0]); err != nil {
				log.Fatalf("wrong url type [%s]", err)
			}
			checkStatus(args[0])
			if ok, _ := cmd.Flags().GetBool("statistics"); ok {
				getStatitics(args[0])
			}
		}
	},
}

// Check Status of the website
func checkStatus(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err.Error())
	} else if resp.StatusCode != 200 {
		color.Red.Println("DOWN!!")
	} else {
		color.Green.Println("UP!")
	}
}

// Get ping details
func getStatitics(url string) {
	hostname := extractDomain(url)
	if len(hostname) == 0 {
		log.Fatal("Invalid url")
		return
	}
	pinger, err := ping.NewPinger(hostname)
	if err != nil {
		panic(err)
	}

	pinger.Count = 1
	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}
	pinger.Run()
}

// parse the domain from the url
func extractDomain(urlLikeString string) string {

	urlLikeString = strings.TrimSpace(urlLikeString)

	if regexp.MustCompile(`^https?`).MatchString(urlLikeString) {
		read, _ := url.Parse(urlLikeString)
		urlLikeString = read.Host
	}

	if regexp.MustCompile(`^www\.`).MatchString(urlLikeString) {
		urlLikeString = regexp.MustCompile(`^www\.`).ReplaceAllString(urlLikeString, "")
	}

	return regexp.MustCompile(`([a-z0-9\-]+\.)+[a-z0-9\-]+`).FindString(urlLikeString)
}
