/*
Copyright © 2022 extract

*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	linkhref "github.com/secgo/vuln/LinkHref"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Extract",
	Short: "Extract Value attribute [From tag HTML]",
	Long: `Extract Value attribute [From tag HTML]
	Example: ./vuln -u https://example.com -t a -a href
	-u short url
	-t short tag
	-a short attribut
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if u != "" && t != "" && a != "" {
			ch := make(chan []byte)
			go VisitUrl(u, ch)
			docHtml := strings.NewReader(string(<-ch))
			lh, err := linkhref.Parse(docHtml, t, a)
			if err != nil {
				panic(err)
			}
			for _, v := range lh {
				fmt.Println(v.Href)
			}
		} else {
			cmd.Help()
		}
	},
	Version: "v1.0.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var u, t, a string

func init() {
	rootCmd.PersistentFlags().StringVarP(&u, "url", "u", "", "Add link https://example.com")
	rootCmd.PersistentFlags().StringVarP(&t, "tag", "t", "", "Add tag [a/div/img/script/...more]")
	rootCmd.PersistentFlags().StringVarP(&a, "attribute", "a", "", "Add attribute [value/src/href/type/...more]")

}

// visit url
func VisitUrl(u string, chh chan []byte) {
	var trans = &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: time.Second,
		}).DialContext,
	}
	res, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		chh <- []byte(err.Error())

	}
	res.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	res.Header.Set("Connection", "close")

	c := &http.Client{Transport: trans, Timeout: 3 * time.Second}

	resp, err := c.Do(res)
	if err != nil {
		chh <- []byte(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		chh <- []byte(err.Error())
	}
	chh <- body
	close(chh)

}
