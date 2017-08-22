package proxy

import (
	"time"

	"github.com/emersion/go-imap"
)

type mailbox struct {
	u    *user
	name string
	info *imap.MailboxInfo
}

func (m *mailbox) ensureSelected() error {
	if selectedMbox := m.u.c.Mailbox(); selectedMbox != nil && selectedMbox.Name == m.name {
		return nil
	}

	_, err := m.u.c.Select(m.name, false)
	return err
}

func (m *mailbox) Name() string {
	return m.name
}

func (m *mailbox) Info() (*imap.MailboxInfo, error) {
	return m.info, nil
}

func (m *mailbox) Status(items []string) (*imap.MailboxStatus, error) {
	if selectedMailbox := m.u.c.Mailbox(); selectedMailbox != nil && selectedMailbox.Name == m.name {
		return selectedMailbox, nil
	}

	return m.u.c.Status(m.name, items)
}

func (m *mailbox) SetSubscribed(subscribe bool) error {
	if subscribe {
		return m.u.c.Subscribe(m.name)
	} else {
		return m.u.c.Unsubscribe(m.name)
	}
}

func (m *mailbox) Check() error {
	if err := m.ensureSelected(); err != nil {
		return err
	}

	return m.u.c.Check()
}

func (m *mailbox) ListMessages(uid bool, seqset *imap.SeqSet, items []string, ch chan<- *imap.Message) error {
	defer close(ch)

	if err := m.ensureSelected(); err != nil {
		return err
	}

	messages := make(chan *imap.Message)
	done := make(chan error, 1)
	go func() {
		if uid {
			done <- m.u.c.UidFetch(seqset, items, messages)
		} else {
			done <- m.u.c.Fetch(seqset, items, messages)
		}
	}()

	for msg := range messages {
		ch <- msg
	}

	return <-done
}

func (m *mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	if err := m.ensureSelected(); err != nil {
		return nil, err
	}

	if uid {
		return m.u.c.UidSearch(criteria)
	} else {
		return m.u.c.Search(criteria)
	}
}

func (m *mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return m.u.c.Append(m.name, flags, date, body)
}

func (m *mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	if err := m.ensureSelected(); err != nil {
		return err
	}

	if uid {
		return m.u.c.UidStore(seqset, string(operation), flags, nil)
	} else {
		return m.u.c.Store(seqset, string(operation), flags, nil)
	}
}

func (m *mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error {
	if err := m.ensureSelected(); err != nil {
		return err
	}

	if uid {
		return m.u.c.UidCopy(seqset, dest)
	} else {
		return m.u.c.Copy(seqset, dest)
	}
}

func (m *mailbox) Expunge() error {
	if err := m.ensureSelected(); err != nil {
		return err
	}

	return m.u.c.Expunge(nil)
}
