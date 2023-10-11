CREATE TABLE IF NOT EXISTS {prefix}webperf_rum_grant_hostnames (
    username                        LowCardinality(String),
    hostname                        LowCardinality(String),
    owner_username                  LowCardinality(String),
    INDEX index_owner owner_username TYPE bloom_filter GRANULARITY 1
)
ENGINE = MergeTree
ORDER BY (username, hostname)
PARTITION BY username
