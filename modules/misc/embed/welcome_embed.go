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
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

const yellowAccent = 0xFAD900

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
	dc.SetRGBA(0, 0, 0, 0.6)
	dc.DrawStringAnchored(text, x+2, y+2, 0.5, 0.5)
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

	dc.SetRGBA(0, 0, 0, 0.45)
	dc.DrawRectangle(0, 0, w, h)
	dc.Fill()

	resp, err := http.Get(avatarURL)
	if err == nil {
		defer resp.Body.Close()
		avImg, _, err := image.Decode(resp.Body)
		if err == nil {
			size := int(w * 0.22)
			cx := int(w / 2)
			cy := int(h * 0.38)
			radius := float64(size) * 0.18

			borderSize := size + 8
			borderDC := gg.NewContext(borderSize, borderSize)
			borderDC.SetRGB(0.98, 0.85, 0.13)
			borderDC.DrawRoundedRectangle(0, 0, float64(borderSize), float64(borderSize), radius+2)
			borderDC.Fill()
			dc.DrawImageAnchored(borderDC.Image(), cx, cy, 0.5, 0.5)
			avDC := gg.NewContext(size, size)
			avDC.DrawRoundedRectangle(0, 0, float64(size), float64(size), radius)
			avDC.Clip()
			avDC.DrawImageAnchored(avImg, size/2, size/2, 0.5, 0.5)
			dc.DrawImageAnchored(avDC.Image(), cx, cy, 0.5, 0.5)
		}
	}

	textX := w / 2
	nameY := h * 0.64
	memberY := h * 0.76

	setFont(dc, h*0.055)
	drawShadowText(dc, fmt.Sprintf("Willkommen, %s!", displayName), textX, nameY)
	dc.SetRGB(0.98, 0.85, 0.13)
	dc.DrawStringAnchored(fmt.Sprintf("Willkommen, %s!", displayName), textX, nameY, 0.5, 0.5)

	setFont(dc, h*0.038)
	drawShadowText(dc, fmt.Sprintf("Du bist Mitglied #%d", memberCount), textX, memberY)
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(fmt.Sprintf("Du bist Mitglied #%d", memberCount), textX, memberY, 0.5, 0.5)

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, dc.Image()); err != nil {
		return nil, err
	}

	return buf, nil
}

func BuildWelcomeComponents(avatarURL, mainChannelID, displayName string, memberCount int) ([]discordgo.MessageComponent, error) {
	content := v2.NewTextDisplayBuilder().SetContent(fmt.Sprintf(
		"# Willkommen 👋\nDu bist jetzt teil der gefangenen im Van!\nOb du entkommst? Natürlich nicht :3\n\nViel Spaß in der Haupthalle <#%s>\n\n-# Willkommen %s, du bist Mitglied #%d",
		mainChannelID, displayName, memberCount,
	)).Build()

	section := v2.NewSectionBuilder()
	section.AddComponent(content)
	section.SetAccessory(v2.NewThumbnailBuilder().SetURL(avatarURL).Build())

	mg := v2.NewMediaGalleryBuilder()
	mg.AddImageURL("attachment://welcome.png")

	container := v2.NewContainerBuilder().
		SetAccentColor(yellowAccent).
		AddComponent(section.Build()).
		AddComponent(mg.Build())

	return []discordgo.MessageComponent{container.Build()}, nil
}