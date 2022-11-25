FROM grafana/grafana:9.0.5-ubuntu

# Disable Login form or not
ENV GF_AUTH_DISABLE_LOGIN_FORM "true"
# Allow anonymous authentication or not
ENV GF_AUTH_ANONYMOUS_ENABLED "true"
# Role of anonymous user
ENV GF_AUTH_ANONYMOUS_ORG_ROLE "Admin"
# Install plugins here our in your own config file
ENV GF_INSTALL_PLUGINS="vertamedia-clickhouse-datasource,ae3e-plotly-panel"