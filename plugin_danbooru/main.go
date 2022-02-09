// Package deepdanbooru 二次元图片标签识别
package deepdanbooru

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/FloatTech/AnimeAPI/danbooru"
	"github.com/FloatTech/AnimeAPI/saucenao"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img/writer"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/zbputils/control/order"
)

const cachefile = "data/danbooru/"

func init() { // 插件主体
	_ = os.RemoveAll(cachefile)
	err := os.MkdirAll(cachefile, 0755)
	if err != nil {
		panic(err)
	}
	engine := control.Register("danbooru", order.AcquirePrio(), &control.Options{
		DisableOnDefault: false,
		Help: "二次元图片标签识别\n" +
			"- 鉴赏图片[图片]",
	})
	// 上传一张图进行评价
	engine.OnKeywordGroup([]string{"鉴赏图片"}, zero.OnlyGroup, ctxext.MustProvidePicture).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("少女祈祷中..."))
			for _, url := range ctx.State["image_url"].([]string) {
				name := ""
				r, err := saucenao.SauceNAO(url)
				if err != nil {
					name = "未知图片"
				} else {
					name = r[0].Title
				}
				t, err := danbooru.TagURL(name, url)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				digest := md5.Sum(helper.StringToBytes(url))
				f := cachefile + hex.EncodeToString(digest[:])
				if file.IsNotExist(f) {
					_ = writer.SavePNG2Path(f, t)
				}
				ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + f))
			}
		})
}
