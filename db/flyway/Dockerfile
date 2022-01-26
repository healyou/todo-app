FROM flyway/flyway:8.4.1

COPY ./patch/drivers /flyway/drivers
COPY ./patch/sql /flyway/sql
COPY ./patch/flyway-dev.conf /flyway/conf/flyway-dev.conf
COPY ./patch/flyway-prod.conf /flyway/conf/flyway-prod.conf