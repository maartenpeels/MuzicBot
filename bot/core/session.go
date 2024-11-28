package core

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type (
	Session struct {
		Queue              *SongQueue
		discord            *discordgo.Session
		guildId, ChannelId string
		connection         *Connection
	}

	SessionManager struct {
		sessions map[string]*Session
	}

	JoinProperties struct {
		Muted    bool
		Deafened bool
	}
)

func newSession(manager *SessionManager, guildId string, channelId string, discord *discordgo.Session, connection *Connection) *Session {
	session := new(Session)
	session.Queue = NewSongQueue(manager)
	session.discord = discord
	session.guildId = guildId
	session.ChannelId = channelId
	session.connection = connection

	go session.Queue.Start(session)

	return session
}

func (sess *Session) Play(url string) error {
	return sess.connection.Play(url)
}

func (sess *Session) Skip() {
	sess.connection.Skip()
}

func (sess *Session) Stop() {
	sess.connection.Stop()
}

func (sess *Session) SendMessage(content string) {
	_, err := sess.discord.ChannelMessageSend(sess.ChannelId, content)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}
}

func NewSessionManager() *SessionManager {
	return &SessionManager{make(map[string]*Session)}
}

func (manager *SessionManager) GetByGuild(guildId string) *Session {
	for _, sess := range manager.sessions {
		if sess.guildId == guildId {
			return sess
		}
	}
	return nil
}

func (manager *SessionManager) Join(discord *discordgo.Session, guildId, channelId string,
	properties JoinProperties) (*Session, error) {
	vc, err := discord.ChannelVoiceJoin(guildId, channelId, properties.Muted, properties.Deafened)
	if err != nil {
		return nil, err
	}
	sess := newSession(manager, guildId, channelId, discord, NewConnection(vc))
	manager.sessions[channelId] = sess
	return sess, nil
}

func (manager *SessionManager) Leave(session Session) {
	session.Queue.Stop()
	session.connection.Stop()
	session.connection.Disconnect()
	delete(manager.sessions, session.ChannelId)
}

func (manager *SessionManager) LeaveAll() {
	for _, session := range manager.sessions {
		session.Queue.Stop()
		session.connection.Stop()
		session.connection.Disconnect()
	}
	manager.sessions = make(map[string]*Session)
}
