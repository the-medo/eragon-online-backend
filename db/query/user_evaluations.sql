-- name: GetEvaluationById :one
SELECT * FROM evaluations WHERE id = @evaluation_id;

-- name: GetEvaluationsByType :many
SELECT * FROM evaluations WHERE evaluation_type = @evaluation_type;

-- name: CreateEvaluationVote :one
INSERT INTO evaluation_votes
(
    evaluation_id,
    user_id,
    user_id_voter,
    value
)
VALUES
    (@evaluation_id, @user_id, @user_id_voter, @value)
RETURNING *;

-- name: GetEvaluationVotesByUserId :many
SELECT * FROM evaluation_votes WHERE user_id = @user_id;

-- name: GetEvaluationVoteByUserIdAndVoter :one
SELECT * FROM evaluation_votes WHERE user_id = @user_id AND user_id_voter = @user_id_voter;

-- name: GetEvaluationVoteByEvaluationIdUserIdAndVoter :one
SELECT * FROM evaluation_votes WHERE evaluation_id = @evaluation_id AND user_id = @user_id AND user_id_voter = @user_id_voter;

-- name: UpdateEvaluationVote :one
UPDATE evaluation_votes
SET
    value = COALESCE(sqlc.arg(value), value),
    created_at = NOW()
WHERE
    evaluation_id = sqlc.arg(evaluation_id) AND
    user_id = sqlc.arg(user_id) AND
    user_id_voter = sqlc.arg(user_id_voter)
RETURNING *;

-- name: DeleteEvaluationVote :exec
DELETE FROM evaluation_votes WHERE evaluation_id = @evaluation_id AND user_id = @user_id AND user_id_voter = @user_id_voter;

-- name: GetAverageUserEvaluationsByType :many
SELECT
    ev.evaluation_id,
    ev.user_id,
    e.name,
    e.description,
    e.evaluation_type,
    AVG(ev.value) AS avg_value
FROM
    evaluation_votes ev
    JOIN evaluations e ON e.id = ev.evaluation_id
WHERE
    ev.user_id = @user_id AND
    e.evaluation_type = @evaluation_type
GROUP BY
    ev.evaluation_id,
    ev.user_id,
    e.name,
    e.description,
    e.evaluation_type;