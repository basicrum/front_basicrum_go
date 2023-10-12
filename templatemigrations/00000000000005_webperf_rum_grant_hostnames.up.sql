CREATE TABLE IF NOT EXISTS {prefix}webperf_rum_grant_hostnames (
    username                        LowCardinality(String),
    hostname                        LowCardinality(String),
    owner_username                  LowCardinality(String),
    updated_at                      DateTime64(3) DEFAULT now(),
    INDEX index_owner owner_username TYPE bloom_filter GRANULARITY 1
)
ENGINE = ReplacingMergeTree(updated_at)
ORDER BY (username, hostname)
PARTITION BY username
PRIMARY KEY (username, hostname)
