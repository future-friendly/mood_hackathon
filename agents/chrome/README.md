# Chrome Extension Monitoring agent

## Purpose
After installation and authentication it starts collecting browsing data and send them to the server.

**Functions include:**
- Token authorization (token from dashboard)
- Turn monitoring on, turn monitoring off
- Log out (delete the monitoring agent)

**Future imporvements:**
- User site filters

## Monitoring
Collects html page source, url and sends to the analytics endpoint with the timestamp.

**Future improvements:**
- Data aggregation on single site
- User browsing graph (track page change branches)
- Time spent metric
