-- name: GetFromUrl :one
SELECT code FROM links
WHERE long_url = ?;

-- name: GetFromCode :one
SELECT long_url FROM links
WHERE code = ?;

-- name: CreateLink :exec
INSERT INTO links (
    code, long_url
) VALUES (
    ?, ?
);