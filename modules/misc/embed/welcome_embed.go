package embed

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	v2 "github.com/S42yt/tuubaa-bot/utils/embed"
	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

const pinkAccent = 0xE8629A

func loadFont(dc *gg.Context, path string, size float64) error {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return err
	}
	face := truetype.NewFace(f, &truetype.Options{Size: size})
	dc.SetFontFace(face)
	return nil
}

func fallbackFont(dc *gg.Context, size float64) {
	f, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(f, &truetype.Options{Size: size})
	dc.SetFontFace(face)
}

func setFont(dc *gg.Context, size float64) {
	if err := loadFont(dc, "assets/fonts/JetBrainsMono-Bold.ttf", size); err != nil {
		fallbackFont(dc, size)
	}
}

func drawShadowText(dc *gg.Context, text string, x, y float64) {
	dc.SetRGBA(0, 0, 0, 0.7)
	dc.DrawStringAnchored(text, x+2, y+2, 0.5, 0.5)
}

func scaleImageToFit(src image.Image, targetSize int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

func BuildWelcomeImage(avatarURL, displayName string, memberCount int) (*bytes.Buffer, error) {
	rand.Seed(time.Now().UnixNano())

	files, err := filepath.Glob("assets/welcome/*")
	if err != nil || len(files) == 0 {
		return nil, err
	}
	bgPath := files[rand.Intn(len(files))]

	bgFile, err := os.Open(bgPath)
	if err != nil {
		return nil, err
	}
	defer bgFile.Close()

	bgImg, _, err := image.Decode(bgFile)
	if err != nil {
		return nil, err
	}

	dc := gg.NewContextForImage(bgImg)
	w := float64(dc.Width())
	h := float64(dc.Height())

	dc.SetRGBA(0, 0, 0, 0.5)
	dc.DrawRectangle(0, 0, w, h)
	dc.Fill()

	avatarSize := int(w * 0.22)
	cx := int(w / 2)
	cy := int(h / 2)
	cornerRadius := float64(avatarSize) * 0.18

	titleY := float64(cy) - float64(avatarSize)/2 - 50
	setFont(dc, h*0.065)
	drawShadowText(dc, "Wilkommen zum goldenen van von tuubaa !!!", w/2, titleY)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored("Wilkommen zum goldenen van von tuubaa !!!", w/2, titleY, 0.5, 0.5)

	resp, err := http.Get(avatarURL)
	if err == nil {
		defer resp.Body.Close()
		avImg, _, err := image.Decode(resp.Body)
		if err == nil {
			scaled := scaleImageToFit(avImg, avatarSize)

			avDC := gg.NewContext(avatarSize, avatarSize)
			avDC.DrawRoundedRectangle(0, 0, float64(avatarSize), float64(avatarSize), cornerRadius)
			avDC.Clip()
			avDC.DrawImage(scaled, 0, 0)

			dc.DrawImageAnchored(avDC.Image(), cx, cy, 0.5, 0.5)
		}
	}

	nameY := float64(cy) + float64(avatarSize)/2 + 60
	setFont(dc, h*0.065)
	drawShadowText(dc, displayName, w/2, nameY)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(displayName, w/2, nameY, 0.5, 0.5)

	memberY := nameY + h*0.1
	setFont(dc, h*0.05)
	drawShadowText(dc, fmt.Sprintf("Member #%d", memberCount), w/2, memberY)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(fmt.Sprintf("Member #%d", memberCount), w/2, memberY, 0.5, 0.5)

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, dc.Image()); err != nil {
		return nil, err
	}
	return buf, nil
}

func BuildWelcomeComponents(avatarURL, mainChannelID, displayName string, memberCount int) ([]discordgo.MessageComponent, error) {
	content := v2.NewTextDisplayBuilder().SetContent(fmt.Sprintf(
		"# Willkommen 👋\nDu bist jetzt teil der gefangenen im Van!\nOb du entkommst? Natürlich nicht :3\n\nViel Spaß in der Haupthalle <#%s>",
		mainChannelID,
	)).Build()

	section := v2.NewSectionBuilder()
	section.AddComponent(content)
	section.SetAccessory(v2.NewThumbnailBuilder().SetURL(avatarURL).Build())

	mg := v2.NewMediaGalleryBuilder()
	mg.AddImageURL("attachment://welcome.png")

	container := v2.NewContainerBuilder().
		SetAccentColor(pinkAccent).
		AddComponent(section.Build()).
		AddComponent(mg.Build())

	return []discordgo.MessageComponent{container.Build()}, nil
}