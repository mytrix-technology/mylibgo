package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/mytrix-technology/mylibgo/networking/http/commons"
	"github.com/mytrix-technology/mylibgo/utils/helper"
	"github.com/mytrix-technology/mylibgo/networking/http/client"
)

type Slack struct {
	WebhookURL string
	Message    string
	Host       string
	Version    string

	StackTrace []byte
}

type slackNotificationDTO struct {
	Text        string               `json:"text,omitempty"`
	Webhook     string               `json:"-"`
	Attachments []slackAttachmentDTO `json:"attachments,omitempty"`
}

type slackAttachmentDTO struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	Color string `json:"color,omitempty"`
}

func (s *Slack) PostInfo() {
	notifDTO := &slackNotificationDTO{
		Webhook: s.WebhookURL,
		Attachments: []slackAttachmentDTO{
			slackAttachmentDTO{
				Title: s.Message,
				Text:  s.createTags(),
				Color: "#3AA3E3",
			},
		},
	}

	notifDTO.post()
}

func (s *Slack) PostError() {
	notifDTO := &slackNotificationDTO{
		Text:    s.Message,
		Webhook: s.WebhookURL,
		Attachments: []slackAttachmentDTO{
			slackAttachmentDTO{
				Title: "Stack Trace",
				Text:  s.createTags() + "\n" + string(s.StackTrace),
				Color: "danger",
			},
		},
	}

	notifDTO.post()
}

func (s *Slack) PostNotif() {
	notifDTO := &slackNotificationDTO{
		Text:    s.Message,
		Webhook: s.WebhookURL,
	}

	notifDTO.post()
}

func (s *Slack) createTags() string {
	tags := []string{}
	tags = append(tags, fmt.Sprintf("`host:%s`", s.Host))
	tags = append(tags, fmt.Sprintf("`version:%s`", s.Version))

	tagsString := strings.Join(tags, " ")
	return tagsString
}

func (s *slackNotificationDTO) post() {
	dataRequest, _ := json.Marshal(s)
	httpCall := commons.HttpCall{
		Method:      commons.HTTP_CALL_METHOD_POST,
		URL:         s.Webhook,
		DataRequest: dataRequest,
		ContentType: commons.HTTP_CALL_CONTENT_JSON,
	}

	_, err := httpCall.SendRequest()
	if err != nil {
		fmt.Println("error when send post slack. ", err)
	}
}

type SlackIntegration struct {
	name string
	webhookURL string
	tags []string
}

func NewSlackIntegration(name, webhookURL string, tags ...string) *SlackIntegration {
	return &SlackIntegration{
		name:      name,
		webhookURL: webhookURL,
		tags: tags,
	}
}

type SlackPostDTO struct {
	Tags []string `json:"-"`
	Text string `json:"text"`
	Channel null.String `json:"channel,omitempty"`
	Username null.String `json:"username,omitempty"`
	IconURL null.String `json:"icon_url,omitempty"`
	IconEmoji null.String `json:"icon_emoji,omitempty"`
	Attachments []SlackPostAttachmentDTO `json:"attachments,omitempty"`
}

type SlackPostAttachmentOption func(dto *SlackPostAttachmentDTO)
type SlackPostAttachmentDTO struct {
	Fallback null.String `json:"fallback"`
	Text null.String `json:"text,omitempty"`
	PreText null.String `json:"pre_text,omitempty"`
	Color null.String `json:"color,omitempty"`
	Title null.String `json:"title,omitempty"`
	Fields []SlackPostAttachmentFieldDTO `json:"fields,omitempty"`
}

type SlackPostAttachmentFieldDTO struct {
	Title null.String `json:"title,omitempty"`
	Value null.String `json:"value,omitempty"`
	Short bool `json:"short,omitempty"`
}

func WithColor(color string) SlackPostAttachmentOption {
	return func(a *SlackPostAttachmentDTO) {
		a.Color = null.StringFrom(color)
	}
}

func WithTitle(title string) SlackPostAttachmentOption {
	return func(a *SlackPostAttachmentDTO) {
		a.Title = null.StringFrom(title)
	}
}

func AddField(title, value string, short bool) SlackPostAttachmentOption {
	return func(a *SlackPostAttachmentDTO) {
		field := SlackPostAttachmentFieldDTO{
			Title: null.StringFrom(title),
			Value: null.StringFrom(value),
			Short: short,
		}
		a.Fields = append(a.Fields, field)
	}
}

type SlackPostOption func(dto *SlackPostDTO)

func WithUserName(username string) SlackPostOption {
	return func(c *SlackPostDTO) {
		c.Username = null.StringFrom(username)
	}
}

func WithChannel(channel string) SlackPostOption {
	return func(c *SlackPostDTO) {
		c.Channel = null.StringFrom(channel)
	}
}

func WithIconURL(iconURL string) SlackPostOption {
	return func(c *SlackPostDTO) {
		c.IconURL = null.StringFrom(iconURL)
	}
}

func WithIconEmoji(iconEmoji string) SlackPostOption {
	return func(c *SlackPostDTO) {
		c.IconEmoji = null.StringFrom(iconEmoji)
	}
}

func AddTags(tags ...string) SlackPostOption {
	return func(c *SlackPostDTO) {
		for _, tag := range tags {
			c.Tags = append(c.Tags, tag)
		}
	}
}

func AddAttachment(text string, options ...SlackPostAttachmentOption) SlackPostOption {
	return func(c *SlackPostDTO) {
		attach := SlackPostAttachmentDTO{}
		if text != "" {
			attach.Fallback = null.StringFrom(text)
			attach.Text = null.StringFrom(text)
		}
		for _, op := range options {
			op(&attach)
		}
		c.Attachments = append(c.Attachments, attach)
	}
}

func (si *SlackIntegration) PostInfo(message string, options ...SlackPostOption) error {
	post := si.buildPost(message, "#3AA3E3", options)
	return si.post(post)
}

func (si *SlackIntegration) PostError(message string, options ...SlackPostOption) error {
	post := si.buildPost(message, "danger", options)

	return si.post(post)
}

func (si *SlackIntegration) post(req *SlackPostDTO) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("slack integration failed to encode request body. %s", err)
	}

	httpclient := client.NewWithTimeout(3 * time.Second, 3 * time.Second)
	if _, err = httpclient.Method(http.MethodPost).
		URL(si.webhookURL).
		BodyWithType(body, helper.HttpContentTypeJson).
		Call(); err != nil {
			return fmt.Errorf("failed to send notif to slack. %s", err)
	}

	return nil
}

func (si *SlackIntegration) buildPost(message string, color string, options []SlackPostOption) *SlackPostDTO {
	post := SlackPostDTO{}
	for _, op := range options {
		op(&post)
	}

	tagLine := si.createTagLine(post.Tags)
	post.Text = tagLine + "\n" + message

	for idx, _ := range post.Attachments {
		post.Attachments[idx].Color = null.StringFrom(color)
	}

	if !post.Username.Valid {
		post.Username = null.StringFrom(si.name)
	}

	return &post
}

func (si *SlackIntegration) createTagLine(tags []string) string {
	if len(tags) > 0 {
		tags = append(si.tags, tags...)
	} else {
		tags = si.tags
	}

	tagStr := ""
	for _, tg := range tags {
		if len(tagStr) > 0 {
			tagStr += " "
		}
		tagStr += "`" + tg + "`"
	}

	return tagStr
}