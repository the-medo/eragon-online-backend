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

const TaskSendResetPasswordEmail = "task:send_reset_password_email"

type PayloadSendResetPasswordEmail struct {
	Email string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendResetPasswordEmail(
	ctx context.Context,
	payload *PayloadSendResetPasswordEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendResetPasswordEmail, jsonPayload, opts...)

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

func (processor *RedisTaskProcessor) ProcessTaskSendResetPasswordEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendResetPasswordEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	resetPassword, err := processor.store.AddUserPasswordReset(ctx, db.AddUserPasswordResetParams{
		UserID: user.ID,
		Code:   util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create password reset: %w", err)
	}

	subject := "Password reset"

	resetUrl := fmt.Sprintf("%s/reset?code=%s", processor.config.FullDomain, resetPassword.Code)

	emailContent := fmt.Sprintf(EmailTemplate, fmt.Sprintf(`
		Hello,<br/>
		We received a request to reset the password for your account. If you initiated this request, please click the button below to create a new password.<br/><br/>
		<div class="button-wrapper">
			<a href="%s" class="button">Reset Your Password</a>
		</div>
		<br/><br/>

		If the button above does not work, please copy and paste the following URL into your browser:<br/>

		<a href="%s">%s</a>
		<br/><br/>

		Please note that this link will expire in 15 minutes. If you did not request a password reset, you can safely ignore this email. If you believe someone else is attempting to access your account, please contact our support team immediately.<br/><br/>

		Thank you for using Talebound!<br/><br/>

		Best regards,<br/>
		The Talebound Team<br/>
	`, resetUrl, resetUrl, resetUrl))

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
