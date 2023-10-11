CREATE OR REPLACE VIEW {prefix}webperf_rum_view_hostnames AS 
SELECT username, hostname, 'owner' as role_name
FROM {prefix}webperf_rum_own_hostnames
UNION ALL
SELECT username, hostname, 'granted' as role_name
FROM {prefix}webperf_rum_grant_hostnames
