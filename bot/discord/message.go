package discord

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/member-gentei/member-gentei/pkg/common"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
)

// MessagingBot allows sending ad-hoc messages.
type MessagingBot struct {
	ctx       context.Context
	dgSession *discordgo.Session
	fs        *firestore.Client

	templates    map[string]*template.Template
	userChannels map[string]*discordgo.Channel
}

// Message sends a templated message to a single user.
func (m *MessagingBot) Message(templateName, uid string, mustBeRegistered bool) error {
	// attempt to get a user
	var user common.DiscordIdentity
	doc, err := m.fs.Collection(common.UsersCollection).Doc(uid).Get(m.ctx)
	if c := status.Code(err); c == codes.NotFound {
		if mustBeRegistered {
			return fmt.Errorf("user not registered: %s", uid)
		}
		user.UserID = uid
		err = nil
	} else if err != nil {
		return err
	} else {
		err = doc.DataTo(&user)
		if err != nil {
			return err
		}
	}
	// load or create template and DM channel
	dmTemplate, err := m.getOrParseTemplate(templateName)
	if err != nil {
		return err
	}
	userChannel, err := m.getOrCreateUserChannel(uid)
	if err != nil {
		return err
	}
	// fill in template data
	templateData := common.DMTemplateData{
		User: &user,
	}
	var buf bytes.Buffer
	err = dmTemplate.Execute(&buf, templateData)
	if err != nil {
		return err
	}
	_, err = m.dgSession.ChannelMessageSend(userChannel.ID, buf.String())
	return err
}

func (m *MessagingBot) getOrCreateUserChannel(uid string) (*discordgo.Channel, error) {
	if userChannel := m.userChannels[uid]; userChannel != nil {
		return userChannel, nil
	}
	userChannel, err := m.dgSession.UserChannelCreate(uid)
	if err != nil {
		return nil, err
	}
	m.userChannels[uid] = userChannel
	return userChannel, nil
}

func (m *MessagingBot) getOrParseTemplate(templateName string) (*template.Template, error) {
	if parsed := m.templates[templateName]; parsed != nil {
		return parsed, nil
	}
	doc, err := m.fs.Collection(common.DMTemplateCollection).Doc(templateName).Get(m.ctx)
	if err != nil {
		return nil, err
	}
	var dmTemplate common.DMTemplate
	err = doc.DataTo(&dmTemplate)
	if err != nil {
		return nil, err
	}
	parsed, err := template.New(templateName).Parse(dmTemplate.Body)
	if err != nil {
		return nil, err
	}
	m.templates[templateName] = parsed
	return parsed, nil
}

// NewMessagingBot initializes a new MessagingBot, which can do nothing but send messages.
func NewMessagingBot(ctx context.Context, token string, fs *firestore.Client) (*MessagingBot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &MessagingBot{
		ctx:          ctx,
		dgSession:    dg,
		fs:           fs,
		templates:    map[string]*template.Template{},
		userChannels: map[string]*discordgo.Channel{},
	}, nil
}
