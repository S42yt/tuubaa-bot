package embed

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

const (
	ComponentTypeContainer    discordgo.ComponentType = 17
	ComponentTypeMediaGallery discordgo.ComponentType = 12
	ComponentTypeTextDisplay  discordgo.ComponentType = 10
	ComponentTypeSection      discordgo.ComponentType = 9
	ComponentTypeThumbnail    discordgo.ComponentType = 11
)

type TextDisplay struct {
	Content string `json:"content"`
}

func (t TextDisplay) Type() discordgo.ComponentType { return ComponentTypeTextDisplay }

func (t TextDisplay) MarshalJSON() ([]byte, error) {
	type textDisplay TextDisplay
	return json.Marshal(struct {
		textDisplay
		Type discordgo.ComponentType `json:"type"`
	}{textDisplay: textDisplay(t), Type: t.Type()})
}

func NewTextDisplayBuilder() *TextDisplay { return &TextDisplay{} }

func (t *TextDisplay) SetContent(content string) *TextDisplay {
	t.Content = content
	return t
}

func (t *TextDisplay) Build() discordgo.MessageComponent { return *t }

type UnfurledMediaItem struct {
	URL string `json:"url"`
}

type MediaGalleryItem struct {
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler"`
}

type MediaGallery struct {
	ID    int                `json:"id,omitempty"`
	Items []MediaGalleryItem `json:"items"`
}

func (m MediaGallery) Type() discordgo.ComponentType { return ComponentTypeMediaGallery }

func (m MediaGallery) MarshalJSON() ([]byte, error) {
	type mediaGallery MediaGallery
	return json.Marshal(struct {
		mediaGallery
		Type discordgo.ComponentType `json:"type"`
	}{mediaGallery: mediaGallery(m), Type: m.Type()})
}

func NewMediaGalleryBuilder() *MediaGallery {
	return &MediaGallery{Items: []MediaGalleryItem{}}
}

func (m *MediaGallery) AddImageURL(url string) *MediaGallery {
	m.Items = append(m.Items, MediaGalleryItem{
		Media:   UnfurledMediaItem{URL: url},
		Spoiler: false,
	})
	return m
}

func (m *MediaGallery) Build() discordgo.MessageComponent { return *m }

type Thumbnail struct {
	Media       UnfurledMediaItem `json:"media"`
	Description *string           `json:"description,omitempty"`
	Spoiler     bool              `json:"spoiler"`
}

func (t Thumbnail) Type() discordgo.ComponentType { return ComponentTypeThumbnail }

func (t Thumbnail) MarshalJSON() ([]byte, error) {
	type thumbnail Thumbnail
	return json.Marshal(struct {
		thumbnail
		Type discordgo.ComponentType `json:"type"`
	}{thumbnail: thumbnail(t), Type: t.Type()})
}

func NewThumbnailBuilder() *Thumbnail { return &Thumbnail{} }

func (t *Thumbnail) SetURL(url string) *Thumbnail {
	t.Media = UnfurledMediaItem{URL: url}
	return t
}

func (t *Thumbnail) SetDescription(desc string) *Thumbnail {
	t.Description = &desc
	return t
}

func (t *Thumbnail) Build() discordgo.MessageComponent { return *t }

type Section struct {
	Components []discordgo.MessageComponent `json:"components"`
	Accessory  discordgo.MessageComponent   `json:"accessory"`
}

func (s Section) Type() discordgo.ComponentType { return ComponentTypeSection }

func (s Section) MarshalJSON() ([]byte, error) {
	type section Section
	return json.Marshal(struct {
		section
		Type discordgo.ComponentType `json:"type"`
	}{section: section(s), Type: s.Type()})
}

func NewSectionBuilder() *SectionBuilder {
	return &SectionBuilder{}
}

type SectionBuilder struct {
	components []discordgo.MessageComponent
	accessory  discordgo.MessageComponent
}

func (s *SectionBuilder) AddComponent(comp discordgo.MessageComponent) *SectionBuilder {
	s.components = append(s.components, comp)
	return s
}

func (s *SectionBuilder) SetAccessory(comp discordgo.MessageComponent) *SectionBuilder {
	s.accessory = comp
	return s
}

func (s *SectionBuilder) Build() discordgo.MessageComponent {
	return Section{
		Components: s.components,
		Accessory:  s.accessory,
	}
}

type Container struct {
	ID          int                          `json:"id,omitempty"`
	AccentColor *int                         `json:"accent_color,omitempty"`
	Spoiler     bool                         `json:"spoiler"`
	Components  []discordgo.MessageComponent `json:"components"`
}

func (c Container) Type() discordgo.ComponentType { return ComponentTypeContainer }

func (c Container) MarshalJSON() ([]byte, error) {
	type container Container
	return json.Marshal(struct {
		container
		Type discordgo.ComponentType `json:"type"`
	}{container: container(c), Type: c.Type()})
}

func NewContainerBuilder() *Container {
	return &Container{Components: []discordgo.MessageComponent{}}
}

func (c *Container) AddComponent(comp discordgo.MessageComponent) *Container {
	c.Components = append(c.Components, comp)
	return c
}

func (c *Container) SetAccentColor(color int) *Container {
	c.AccentColor = &color
	return c
}

func (c *Container) Build() discordgo.MessageComponent { return *c }