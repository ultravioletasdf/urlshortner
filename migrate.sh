atlas schema apply \
--url "$TURSO_URL?authToken=$TURSO_TOKEN" \
--dev-url "sqlite://dev.db" \
--to "file://sql/schema.sql"