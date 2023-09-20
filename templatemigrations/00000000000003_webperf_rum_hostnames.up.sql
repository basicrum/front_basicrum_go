CREATE TABLE IF NOT EXISTS {prefix}webperf_rum_hostnames (
    hostname                        LowCardinality(String),
    updated_at                      DateTime64(3) DEFAULT now()
)
ENGINE = ReplacingMergeTree
PARTITION BY hostname
ORDER BY hostname
SETTINGS index_granularity = 8192