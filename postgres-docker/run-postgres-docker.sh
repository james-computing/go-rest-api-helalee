# Runs a docker container with PostgreSQL.
# --name some-postgres sets the container name as some-postgres
# -e POSTGRES_PASSWORD=replacewithpassword sets the password to be used by the database
# -p first:second maps the first port from outside the container to the second port inside the container
# -d detaches execution, so you can keep using the terminal while the container is running
# postgres is the container image name
docker run --name some-postgres -e POSTGRES_PASSWORD=replacewithpassword -p 5432:5432 -d postgres