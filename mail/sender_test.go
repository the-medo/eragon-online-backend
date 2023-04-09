package mail

import (
	"github.com/stretchr/testify/require"
	"github.com/the-medo/talebound-backend/util"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
		<style>
			.blue {
				color: blue;	
			}
		</style>
		<h1>Hello world!</h1>
		<p class="blue">This is a text with class.</p>
		<p style="color: red; font-size: 15px;">This is a text with inline css.</p>
	`

	to := []string{"martinmederly@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
