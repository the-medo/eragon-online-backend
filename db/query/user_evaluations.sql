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

-- name: GetEvaluationVotesByUserIdAndVoter :many
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
    e.id AS evaluation_id,
    e.name,
    e.description,
    e.evaluation_type,
    AVG(COALESCE(ev.value, 0)) AS avg_value
FROM
    evaluations e
    LEFT JOIN evaluation_votes ev ON e.id = ev.evaluation_id AND ev.user_id = @user_id
WHERE
    e.evaluation_type = @evaluation_type
GROUP BY
    e.id,
    e.name,
    e.description,
    e.evaluation_type;