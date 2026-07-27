-- +goose Up
DELETE FROM request_assignments duplicate
USING request_assignments original
WHERE duplicate.request_id = original.request_id
  AND duplicate.user_id = original.user_id
  AND duplicate.unassigned_at IS NULL
  AND original.unassigned_at IS NULL
  AND (
      duplicate.assigned_at > original.assigned_at
      OR (
          duplicate.assigned_at = original.assigned_at
          AND duplicate.id > original.id
      )
  );

CREATE UNIQUE INDEX request_assignments_active_user_request_idx
    ON request_assignments (request_id, user_id)
    WHERE unassigned_at IS NULL;

-- +goose Down
DROP INDEX request_assignments_active_user_request_idx;
