package worker

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/util"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	UserId int32 `json:"user_id"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Bytes("payload", task.Payload()).
		Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUserById(ctx, payload.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		UserID:     user.ID,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to Talebound!"

	verifyUrl := fmt.Sprintf("%s/verify?id=%d&secret_code=%s", processor.config.FullDomain, verifyEmail.ID, verifyEmail.SecretCode)

	emailContent := fmt.Sprintf(EmailTemplate, fmt.Sprintf(`
		Hello <b>%s</b>,<br/><br/>

                Welcome to <b>Talebound</b>, the ultimate platform for text-based role-playing games and immersive storytelling experiences! We're thrilled to have you join our community of dreamers, creators, and adventurers.<br/><br/>

                Before you can fully enjoy all the features that Talebound has to offer, we need you to verify your email address. By doing so, you help us ensure the security and authenticity of our community.<br/><br/>

                Please click the button below to verify your email address:<br/>

                <div class="button-wrapper">
                    <a href="%s" target="_blank" class="button">Verify email</a>
                </div>

                If you are unable to click the button, copy and paste the following URL into your browser:<br/>

                <a href="%s" target="_blank">%s</a>
                <br/><br/>

                Once your email address is verified, you'll be able to start creating characters, joining campaigns, and exploring the vast realms of imagination that Talebound has to offer.<br/><br/>

                If you have any questions or need assistance, please don't hesitate to reach out to our support team at <a href="mailto:support@talebound.net">support@talebound.net</a><br/><br/>

                Thank you for joining Talebound, and we look forward to seeing the incredible stories you'll create and adventures you'll embark upon.<br/><br/>

                Happy storytelling!<br/>
                The Talebound Team<br/>
	`, user.Username, verifyUrl, verifyUrl, verifyUrl))

	to := []string{user.Email}
	err = processor.mailer.SendEmail(subject, emailContent, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")

	return nil
}
