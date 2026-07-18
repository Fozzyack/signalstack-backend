# SignalStack Database Design

Reference schema for the IT request and expert workspace. The current backend
already contains `users` and `sessions` tables.

## Current Tables

### `users`

Authentication and IT expert accounts.

```text
id UUID PRIMARY KEY
email TEXT UNIQUE NOT NULL
name TEXT NOT NULL
password_hash TEXT NOT NULL
role TEXT NOT NULL                 -- admin, expert
is_active BOOLEAN NOT NULL
created_at TIMESTAMPTZ NOT NULL
updated_at TIMESTAMPTZ NOT NULL
```

### `sessions`

Authenticated login sessions linked to `users`.

```text
id UUID PRIMARY KEY
user_id UUID REFERENCES users(id) ON DELETE CASCADE
token TEXT UNIQUE NOT NULL
expires_at TIMESTAMPTZ NOT NULL
created_at TIMESTAMPTZ NOT NULL
```

## MVP Tables

### `requests`

An IT request submitted from the landing page.

```text
id UUID PRIMARY KEY
reference TEXT UNIQUE NOT NULL       -- SS-1048
title TEXT NOT NULL
nescription TEXT NOT NULL
client_name TEXT NOT NULL
client_email TEXT NOT NULL
status TEXT NOT NULL                  -- new, in_progress, waiting, resolved, closed
created_at TIMESTAMPTZ NOT NULL
updated_at TIMESTAMPTZ NOT NULL
resolved_at TIMESTAMPTZ
```

Keeping the submitted client name and email on the request is appropriate for
the MVP. It preserves the exact contact details supplied at submission time.

Recommended indexes:

- `requests(status)`
- `requests(created_at)`
- `requests(client_email)`

### `request_assignments`

Join table between requests and experts. This supports multiple experts on one
request.

```text
id UUID PRIMARY KEY
request_id UUID REFERENCES requests(id) ON DELETE CASCADE
user_id UUID REFERENCES users(id) ON DELETE CASCADE
role TEXT NOT NULL                    -- lead, contributor
assigned_at TIMESTAMPTZ NOT NULL
unassigned_at TIMESTAMPTZ
personal_deadline TIMESTAMPTZ
completed_at TIMESTAMPTZ
```

The personal deadline belongs here because two experts may have different
deadlines for the same request.

Recommended constraints and indexes:

- Prevent duplicate active assignments for the same user and request.
- `request_assignments(request_id)`
- `request_assignments(user_id)`
- `request_assignments(personal_deadline)`

### `request_notes`

Shared internal notes visible to the IT team.

```text
id UUID PRIMARY KEY
request_id UUID REFERENCES requests(id) ON DELETE CASCADE
author_id UUID REFERENCES users(id)
body TEXT NOT NULL
created_at TIMESTAMPTZ NOT NULL
updated_at TIMESTAMPTZ NOT NULL
```

### `task_notes`

Private notes belonging to one expert's assignment.

```text
id UUID PRIMARY KEY
assignment_id UUID REFERENCES request_assignments(id) ON DELETE CASCADE
author_id UUID REFERENCES users(id) ON DELETE CASCADE
body TEXT NOT NULL
created_at TIMESTAMPTZ NOT NULL
updated_at TIMESTAMPTZ NOT NULL
```

Personal notes should be linked to an assignment rather than only a request.
That prevents one expert's private notes from appearing to another expert on
the same request.

### `request_events`

Immutable audit history for meaningful request changes.

```text
id UUID PRIMARY KEY
request_id UUID REFERENCES requests(id) ON DELETE CASCADE
actor_id UUID REFERENCES users(id)
event_type TEXT NOT NULL
metadata JSONB
created_at TIMESTAMPTZ NOT NULL
```

Example event types:

```text
request_created
assignment_added
assignment_removed
lead_changed
status_changed
note_added
deadline_changed
request_resolved
request_reopened
```

Example metadata for a status change:

```json
{
  "old_status": "new",
  "new_status": "in_progress"
}
```

Events should be created by the application service whenever a mutation
succeeds. They should not be edited through normal application flows.

## Specializations

### `skills`

Available expert specializations.

```text
id UUID PRIMARY KEY
name TEXT UNIQUE NOT NULL
created_at TIMESTAMPTZ NOT NULL
```

Examples include Networking, Cloud Infrastructure, Microsoft 365, Hardware,
Security, and Support.

### `user_skills`

Many-to-many relationship between experts and their specializations.

```text
user_id UUID REFERENCES users(id) ON DELETE CASCADE
skill_id UUID REFERENCES skills(id) ON DELETE CASCADE
proficiency TEXT                       -- optional
PRIMARY KEY (user_id, skill_id)
```

### `request_skills` (optional)

Internal classification used to recommend experts for a request.

```text
request_id UUID REFERENCES requests(id) ON DELETE CASCADE
skill_id UUID REFERENCES skills(id) ON DELETE CASCADE
PRIMARY KEY (request_id, skill_id)
```

These classifications do not need to be shown as visible issue tags in the
dashboard. They can remain internal routing metadata.

## Future Tables

### `attachments`

Metadata for files attached to requests. Store the actual files in object
storage, not directly in Postgres.

```text
id UUID PRIMARY KEY
request_id UUID REFERENCES requests(id) ON DELETE CASCADE
uploaded_by UUID REFERENCES users(id)
filename TEXT NOT NULL
storage_key TEXT NOT NULL
content_type TEXT NOT NULL
size_bytes BIGINT NOT NULL
created_at TIMESTAMPTZ NOT NULL
```

### `organizations`

Use this if clients represent recurring companies rather than one-off contact
details.

```text
id UUID PRIMARY KEY
name TEXT NOT NULL
created_at TIMESTAMPTZ NOT NULL
updated_at TIMESTAMPTZ NOT NULL
```

`requests` could then include `organization_id` while retaining submitted
client name and email snapshots.

### `contacts`

Use this if multiple people at one organization can submit or discuss
requests.

```text
id UUID PRIMARY KEY
organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE
name TEXT NOT NULL
email TEXT NOT NULL
created_at TIMESTAMPTZ NOT NULL
updated_at TIMESTAMPTZ NOT NULL
```

### `request_messages`

Use this if clients need to view requests or communicate through the product.
Keep client-visible messages separate from internal and private notes.

```text
id UUID PRIMARY KEY
request_id UUID REFERENCES requests(id) ON DELETE CASCADE
author_user_id UUID REFERENCES users(id)
author_email TEXT
body TEXT NOT NULL
visibility TEXT NOT NULL               -- client, internal
created_at TIMESTAMPTZ NOT NULL
```

### `notifications`

Useful for assignment, deadline, status, and client reply notifications.

```text
id UUID PRIMARY KEY
user_id UUID REFERENCES users(id) ON DELETE CASCADE
request_id UUID REFERENCES requests(id) ON DELETE CASCADE
type TEXT NOT NULL
read_at TIMESTAMPTZ
created_at TIMESTAMPTZ NOT NULL
```

## Relationships

```text
users
  ├── sessions
  ├── user_skills
  ├── request_assignments
  ├── request_notes
  ├── task_notes
  └── request_events

requests
  ├── request_assignments
  ├── request_notes
  ├── request_events
  ├── request_skills
  └── attachments

request_assignments
  └── task_notes
```

## Recommended Build Order

1. Add request creation and request retrieval.
2. Add assignment and self-assignment.
3. Add status changes and audit events.
4. Add personal deadlines and private task notes.
5. Add shared internal notes.
6. Add skills and expert recommendations.
7. Add attachments, client accounts, and notifications.

## Decisions To Revisit

- Will clients ever log in and view their requests?
- Should every request have a lead, or only some requests?
- Can experts assign other experts, or only assign themselves?
- Is a deadline personal to an expert or shared across the request?
- Should request status be shared by everyone or tracked per assignment?
- Are skills only for recommendations, or should they restrict assignment?
- How long should client contact data and audit history be retained?
- Should closed requests be archived or remain searchable indefinitely?

## MVP Recommendation

Start with:

- `users`
- `sessions`
- `requests`
- `request_assignments`
- `task_notes`
- `request_events`
- `skills`
- `user_skills`

This supports the landing-page submission flow, the shared request queue,
multi-expert assignment, the My Tasks page, personal notes, deadlines, and an
