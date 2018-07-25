CREATE DATABASE forex;
CREATE USER forex WITH PASSWORD 'forex';
GRANT ALL PRIVILEGES ON DATABASE forex TO forex;
ALTER USER forex WITH SUPERUSER;

CREATE TABLE exchange_rate (
    id BIGSERIAL NOT NULL,
    from_cur VARCHAR(10) NOT NULL,
    to_cur VARCHAR(10) NOT NULL,

	PRIMARY KEY (id)
);

CREATE TABLE exchange_rate_data (
	id BIGSERIAL NOT NULL,
	exchange_rate_id BIGINT NOT NULL,
	date DATE NOT NULL,
	rate DOUBLE PRECISION NOT NULL,

	PRIMARY KEY (id),
	CONSTRAINT fk_exchange_rate_data_idx_exchange_rate_id
                    FOREIGN KEY (exchange_rate_id)
                    REFERENCES exchange_rate(id)
				ON DELETE CASCADE
);
