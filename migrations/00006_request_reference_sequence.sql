-- +goose Up
CREATE SEQUENCE request_reference_seq START WITH 1000;

SELECT setval(
    'request_reference_seq',
    GREATEST(
        999,
        COALESCE(
            MAX(NULLIF(substring(reference FROM '^SS-([0-9]+)$'), '')::bigint),
            0
        )
    ),
    true
)
FROM requests;

ALTER TABLE requests
ALTER COLUMN reference SET DEFAULT (
    'SS-' || lpad(nextval('request_reference_seq')::text, 4, '0')
);

-- +goose Down
ALTER TABLE requests
ALTER COLUMN reference DROP DEFAULT;

DROP SEQUENCE request_reference_seq;
