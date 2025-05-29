FROM migrate/migrate:latest

# Copy migration files to the container
COPY migrations/ /migrations/

# Set working directory
WORKDIR /migrations

