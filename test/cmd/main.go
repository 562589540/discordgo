package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := "XDINTg1NjU0MDMzNDE2MTky.GYxlGL"
	dg, err := discordgo.New(token)
	if err != nil {
		return
	}
	defer dg.Close()

	dg.Debug = true
	dg.LogLevel = 10

	adapter := discordgo.NewMirrorZlibAdapter()
	dg.SetMirrorZlibAdapter(adapter)
	dg.SetBaseURL("https://discord-d-com-s-mjc.aiwentu.net/")
	dg.SetCookie("___AID_mjc=1; ___AID_mjg2=1; ___MJ4=%7B%22mj%22%3A%7B%22s%22%3A%22%22%2C%22style%22%3A%22%22%2C%22v%22%3A%226.1%22%2C%22w%22%3A%22%22%2C%22params%22%3A%22%22%2C%22desc%22%3A%22%22%7D%2C%22niji%22%3A%7B%22s%22%3A%22%22%2C%22style%22%3A%22%22%2C%22v%22%3A%226%22%2C%22w%22%3A%22%22%2C%22params%22%3A%22%22%2C%22desc%22%3A%22%22%7D%2C%22_v%22%3A%222%22%2C%22toEn%22%3Afalse%2C%22relax%22%3Afalse%7D; ___SID=K0ACFShqFE9-9x9jHMo9v; ___FP2=4vz9xeyJsYW5ndWFnZXMiOiJ6aC1DTix6aCIsInVzZXJBZ2VudCI6Ik1vemlsbGEvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE1XzcpIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS8xMzUuMC4wLjAgU2FmYXJpLzUzNy4zNiIsInRpbWVab25lIjoiQXNpYS9TaGFuZ2hhaSIsInRpbWVab25lT2Zmc2V0Ijo4LCJkZXZpY2VNZW1vcnkiOjgsIm1heFRvdWNoUG9pbnRzIjowLCJoYXJkd2FyZUNvbmN1cnJlbmN5Ijo4LCJwbGF0Zm9ybSI6Im1hY09TIiwic2NyZWVuIjoiMTQ0MCw4MDMiLCJ3ZWJkcml2ZXIiOmZhbHNlLCJ2aWRlb0NhcmQiOiJBTkdMRSAoQXBwbGUsIEFOR0xFIE1ldGFsIFJlbmRlcmVyOiBBcHBsZSBNMSwgVW5zcGVjaWZpZWQgVmVyc2lvbikiLCJjYW52YXNJZCI6MjExNzM2MDQ5LCJkZXNrdG9wIjpmYWxzZSwiaXNQcml2YXRlIjpmYWxzZSwiZnAiOjIxMjc2MzgwNzksImhhc0dmdyI6dHJ1ZSwiaXBDbiI6IjEwNi4xMTcuMjI2LjMxIiwiaXBGcm9tIjoicnRjIiwiY2RwIjp0cnVlfQ..izmv")
	err = dg.Open()
	if err != nil {
		return
	}
	fmt.Println("打开 discord session 成功")
	select {}

}
