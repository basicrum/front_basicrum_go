# Front Basic RUM GO

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Introduction

This is a simple application server written on Golang. The server receives payloads (beacons) sent from visitors who are browsing websites. The payload data that is being sent is collected and packaged by a JavaScript library called Boomerang JS. The main purpose of Boomerang JS is to collect performance related data from web pages.

Below is a simplified diagram of how the data flows for our use case:

![Data flow](./doc/data%20flow.png)

This project focuses on the component Basic RUM GO but this is how the data flows in the completed system that contains various components for data persistency and visualization:

1. GET/POST request is sent from a browser to our application server. The application server does some transformation of the received data.

2. The transformed data is being persisted to a database called ClickHouse.

3. Charts and other visualizations are being created while Grafana is querying ClickHouse.

## How to start dev environment

### 1. Start Front Basic RUM GO

**1.1** Set environment variables

```
# server
export SERVER_HOST="localhost"
export SERVER_PORT=8087

# database
export DATABASE_HOST="localhost"
export DATABASE_PORT=9000
export DATABASE_NAME="default"
export DATABASE_USERNAME="default"
export DATABASE_PASSWORD=""
## optional
export DATABASE_TABLE_PREFIX=""

# persistance
## optional
PERSISTANCE_DATABASE_STRATEGY=""
## optional
PERSISTANCE_TABLE_STRATEGY=""

# backup
## optional - default false
BACKUP_ENABLED=true
## optional. required if BACKUP_ENABLED=true
BACKUP_DIRECTORY=""
## optional
BACKUP_INTERVAL_SECONDS=0
```

**1.2** Start Front Basic RUM GO

```
go run main.go
```

### 2. Start ClickHouse

Run:

```
make up
```

### 3. Create a table in ClickHouse

**3.1** Load http://localhost:8143/ in your browser

**3.2** Run the following SQL:

```sql
CREATE TABLE IF NOT EXISTS integration_test_webperf_rum_events (
  event_date                      Date DEFAULT toDate(created_at),
  hostname                        LowCardinality(String),
  created_at                      DateTime,
  event_type                      LowCardinality(String),
  browser_name                    LowCardinality(String),
  browser_version                 Nullable(String),
  ua_vnd                          LowCardinality(Nullable(String)),
  ua_plt                          LowCardinality(Nullable(String)),
  device_type                     LowCardinality(String),
  device_manufacturer             LowCardinality(Nullable(String)),
  operating_system                LowCardinality(String),
  operating_system_version        Nullable(String),
  user_agent                      Nullable(String),
  next_hop_protocol               LowCardinality(String),
  visibility_state                LowCardinality(String),

  session_id                      FixedString(43),
  session_length                  UInt8,
  url                             String,
  connect_duration                Nullable(UInt16),
  dns_duration                    Nullable(UInt16),
  first_byte_duration             Nullable(UInt16),
  redirect_duration               Nullable(UInt16),
  redirects_count                 UInt8,

  first_contentful_paint          Nullable(UInt16),
  first_paint                     Nullable(UInt16),

  cumulative_layout_shift         Nullable(Float32),
  first_input_delay               Nullable(UInt16),
  largest_contentful_paint        Nullable(UInt16),

  geo_country_code                FixedString(2),
  geo_city_name                   Nullable(String),
  page_id                         FixedString(8),

  data_saver_on                   Nullable(UInt8),

  boomerang_version               LowCardinality(String),
  screen_width                    Nullable(UInt16),
  screen_height                   Nullable(UInt16),

  dom_res                         Nullable(UInt16),
  dom_doms                        Nullable(UInt16),
  mem_total                       Nullable(UInt32),
  mem_limit                       Nullable(UInt32),
  mem_used                        Nullable(UInt32),
  mem_lsln                        Nullable(UInt32),
  mem_ssln                        Nullable(UInt32),
  mem_lssz                        Nullable(UInt32),
  scr_bpp                         Nullable(String),
  scr_orn                         Nullable(String),
  cpu_cnc                         Nullable(UInt8),
  dom_ln                          Nullable(UInt16),
  dom_sz                          Nullable(UInt16),
  dom_ck                          Nullable(UInt16),
  dom_img                         Nullable(UInt16),
  dom_img_uniq                    Nullable(UInt16),
  dom_script                      Nullable(UInt16),
  dom_iframe                      Nullable(UInt16),
  dom_link                        Nullable(UInt16),
  dom_link_css                    Nullable(UInt16)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(event_date)
ORDER BY (device_type, event_date)
SETTINGS index_granularity = 8192
```

### 4. Send some test data

```
curl 'http://127.0.0.1:8087/beacon/catcher' \
  -H 'authority: beacon.basicrum.com' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -H 'origin: https://calendar.perfplanet.com' \
  -H 'pragma: no-cache' \
  -H 'referer: https://calendar.perfplanet.com/' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: no-cors' \
  -H 'sec-fetch-site: cross-site' \
  -H 'user-agent: Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1' \
  --data-raw 'mob.etype=4g&mob.dl=10&mob.rtt=50&c.e=l544045p&c.tti.m=lt&pt.lcp=1234&rt.start=navigation&rt.bmr=1981%2C226%2C221%2C6&rt.tstart=1656779946253&rt.bstart=1656779948466&rt.blstart=1656779948233&rt.end=1656779948292&t_resp=941&t_page=1098&t_done=2039&t_other=boomerang%7C25%2Cboomr_fb%7C2213%2Cboomr_ld%7C1980%2Cboomr_lat%7C233&rt.tt=2039&rt.obo=0&pt.fp=1234&pt.fcp=1234&nt_nav_st=1656779946253&nt_red_st=1656779946262&nt_red_end=1656779946977&nt_fet_st=1656779946977&nt_dns_st=1656779946977&nt_dns_end=1656779946977&nt_con_st=1656779946977&nt_con_end=1656779946977&nt_req_st=1656779946981&nt_res_st=1656779947194&nt_res_end=1656779947207&nt_domloading=1656779947229&nt_domint=1656779948234&nt_domcontloaded_st=1656779948234&nt_domcontloaded_end=1656779948243&nt_domcomp=1656779948292&nt_load_st=1656779948292&nt_load_end=1656779948292&nt_ssl_st=1656779946977&nt_enc_size=12680&nt_dec_size=47273&nt_trn_size=12980&nt_protocol=h2&nt_first_paint=1656779947487&nt_red_cnt=1&nt_nav_type=0&restiming=%7B%22https%3A%2F%2Fcalendar.perfplanet.com%2F%22%3A%7B%22wp-includes%2F%22%3A%7B%22js%2F%22%3A%7B%22wp-embed.min.js%3Fver%3D5.8.4%22%3A%223ri%2Cre%2Cq1%2C4s*1l9%2C_%2Cid*24%22%2C%22wp-emoji-release.min.js%3Fver%3D5.8.4%22%3A%223ws%2Ch7%2Cgv%2Cw*13uu%2C8c%2Ca67*23%22%7D%2C%22css%2Fdist%2Fblock-library%2Fstyle.min.css%3Fver%3D5.8.4%22%3A%222re%2C17%2C12%2Cy*18gp%2C_%2C1hph*44%22%7D%2C%22photos%2F%22%3A%7B%22alexrussell-70tr.jpg%22%3A%22*01y%2C1y%2Ca7%2C103%2C5u%2C5u%7C1rf%2C8e%2C8c%2C5n*17oz%2C_%22%2C%22yoav2016-70tr.jpg%22%3A%22*02d%2C1y%2Cij%2C103%2C2e%2C1y%7C1rf%2Cb1%2Cav%2C5s*13e7%2C_%22%2C%22paulcalvano-70tr.jpg%22%3A%22*01y%2C1y%2Cra%2C103%2C5u%2C5u%7C1rf%2Cbe%2Cb8%2C5t*15el%2C_%22%2C%22nic-70tr.jpg%22%3A%22*01y%2C1y%2Cyt%2C103%7C1rf%2Cgd%2Ccm%2C5t*14kr%2C_%22%2C%22annie-70tr.jpg%22%3A%22*01y%2C1y%2C17y%2C103%2C5u%2C5u%7C1rf%2Cii%2Ccy%2C5w*18r7%2C_%22%2C%22jana-70tr.jpg%22%3A%22*01y%2C1y%2C1ga%2C103%2C5u%2C5u%7C1rf%2Cko%2Ccy%2C5w*195r%2C_%22%2C%22andrea2021-70tr.jpg%22%3A%22*01y%2C1y%2C1nt%2C103%2C5u%2C5u%7C1rf%2Ckr%2Ccy%2C5y*13wr%2C_%22%2C%22richcecko-70tr.jpg%22%3A%22*011%2C1y%2C1w5%2C103%2C2x%2C5u%7C1rg%2Cms%2Ccx%2C5x*1492%2C_%22%2C%22alex-70tr.jpg%22%3A%22*02k%2C1y%2C24h%2C103%2C2l%2C1y%7C1rg%2Cms%2Ccx%2C5x*12p1%2C_%22%2C%22meiert-70tr.jpg%22%3A%22*01y%2C1y%2C2df%2C103%2C5u%2C5u%7C1rg%2Cmt%2Cm7%2C5x*12uo%2C_%22%2C%22artem2021-70tr.jpg%22%3A%22*01y%2C1y%2C2lr%2C103%2C5u%2C5u%7C1rg%2Cmv%2Cms%2C5x*17ae%2C_%22%2C%22Tim-Vereecke-70tr.jpg%22%3A%22*02j%2C1y%2C2ta%2C103%2C2k%2C1y%7C1rg%2Cmv%2Cm7%2C5x*113v%2C_%22%2C%22peter2015-70tr.jpg%22%3A%22*01x%2C1y%2C327%2C103%2C1x%2C1y%7C1rg%2Cmv%2Cmt%2C5x*11k0%2C_%22%2C%22robin2021-70tr.jpg%22%3A%22*01z%2C1y%2C3b0%2C103%2C8w%2C8r%7C1rg%2Cn4%2Cmj%2C5x*1d9m%2C_%22%2C%22amiya-70tr.jpg%22%3A%22*025%2C1y%2C3jd%2C103%2C4b%2C3w%7C1rg%2Cn6%2Cmj%2C5x*122y%2C_%22%2C%22kanmi-70tr.jpg%22%3A%22*01y%2C1y%2C3sz%2C103%2C79%2C78%7C1rh%2Cn5%2Cmi%2C5w*16kz%2C_%22%2C%22brunosabot-70tr.jpg%22%3A%22*01y%2C1y%2C42l%2C103%2C8w%2C8w%7C1rh%2Cng%2Cmi%2C5x*16ch%2C_%22%2C%22stoyan2015-70tr.jpg%22%3A%22*01s%2C1y%2C4c7%2C103%2Cc6%2Cdh%7C1rh%2Cnv%2Cmi%2C5x*1g01%2C_%22%2C%22hongbo-70tr.jpg%22%3A%22*01z%2C1y%2C4kj%2C103%2C5w%2C5u%7C1rh%2Co6%2Cmi%2C5x*14kx%2C_%22%2C%22amit-70tr.jpg%22%3A%22*02c%2C1y%2C53r%2C103%2C74%2C5u%7C1rh%2Cq3%2Cmi%2C5x*15j9%2C_%22%2C%22robertb-70tr.jpg%22%3A%22*027%2C1y%2C5mz%2C103%2C4g%2C3w%7C1rh%2Cqa%2Cmi%2C5x*14a9%2C_%22%2C%22leon2021-70tr.jpg%22%3A%22*022%2C1y%2C5vk%2C103%2C67%2C5u%7C1rh%2Cqa%2Cmi%2C5x*14df%2C_%22%2C%22krasimir-70tr.jpg%22%3A%22*01y%2C1y%2C6t8%2C103%2C5u%2C5u%7C1ri%2Cr9%2Cmh%2C5w*14z8%2C_%22%2C%22leonb-70tr.jpg%22%3A%22*01y%2C1y%2C72u%2C103%2C50%2C50%7C1ri%2Cr9%2Cns%2C5w*12xl%2C_%22%2C%22barry-70tr.jpg%22%3A%22*026%2C1y%2C7bz%2C103%2C4e%2C3w%7C1ri%2Crc%2Cns%2C5w*13dm%2C_%22%2C%22tanner-70tr.jpg%22%3A%22*01y%2C1y%2C7sw%2C103%2C5u%2C5u%7C1ri%2Crd%2Co4%2C5w*180h%2C_%22%2C%22erwin-70tr.jpg%22%3A%22*01y%2C1y%2C82i%2C103%2C5u%2C5u%7C1ri%2Crd%2Cq1%2C5w*16m5%2C_%22%2C%22matt-70tr.jpg%22%3A%22*025%2C1y%2C8au%2C103%2C7y%2C78%7C1ri%2Cre%2Cq1%2C5x*194e%2C_%22%7D%2C%22favicon.ico%22%3A%2201kq%2C5o%2C5n%2C1*1%2C8c%22%2C%222021%2F%22%3A%226%2Cqi%2Cq6%2Ck8%2Ck4%2Ck4%2Ck4%2Ck4%2Ck4%2Ck4%2C9*19s8%2C8c%2Cqox%22%2C%22wp-content%2Fthemes%2Fwpc2%2F%22%3A%7B%22style.css%22%3A%222re%2C16%2C12%2Cy*1182%2C_%2C29p*44%22%2C%22wpclogo.png%22%3A%22*050%2C50%2Ci%2Cbj%2Clv%2Clv%7C1re%2Cgd%2Cd6%2C4w*1ujn%2C_%22%7D%2C%22js%2Fbasicrum%2Fboomerang-1.737.60.cutting-edge.min.js%22%3A%2221j1%2C6a%2C65%2C6*1plb%2C8c%2C1ma1*25*42%22%7D%7D&u=https%3A%2F%2Fcalendar.perfplanet.com%2F2021%2F&v=1.737.60&sv=14&sm=p&rt.si=f7504dae-b2e1-464b-a28e-1311ab4de139-reejl6&rt.ss=1656779946253&rt.sl=1&vis.st=visible&ua.plt=Linux%20x86_64&ua.vnd=Google%20Inc.&pid=ifvwnpyy&n=1&c.l=22hq%2Cxb3%2Cyka%7C52wz%2Ca0%7C22x0%2Cxaw%2Cyh4&c.t.fps=07*d*615*y*62&c.t.mouse=0*9*0.n4.00.9m..6k.&c.t.mousepct=2*a*004000001710&c.t.orn=0*f*01&c.tti.vr=1990&c.tti=2128&c.f=59&c.f.d=5027&c.f.m=1&c.f.l=1&c.f.s=l54405vm&c.m.p=67&c.m.n=1414&dom.res=36&dom.doms=1&mem.total=16164031&mem.limit=2172649472&mem.used=11839111&mem.lsln=0&mem.ssln=0&mem.lssz=2&mem.sssz=2&scr.xy=414x896&scr.bpp=24%2F24&scr.orn=0%2Fportrait-primary&scr.dpx=2&cpu.cnc=8&scr.mtp=1&dom.ln=515&dom.sz=45903&dom.ck=158&dom.img=35&dom.img.uniq=29&dom.script=7&dom.script.ext=3&dom.iframe=0&dom.link=9&dom.link.css=2&sb=1' \
  --compressed
```

### 5. Query the test data

**5.1** Load the ClickHouse dev user interface

Load http://localhost:8143/ in your browser

**5.2** Query the data

Run the follwoing query:

```sql
SELECT * from integration_test_webperf_rum_events
```

### 6. Running the server

Start:

```
./main > app.log 2>&1 &
```

Periodically check the app.log file for any interesting messages

Check if the server is running:

```
netstat -ltnp | grep -w ':8087'
```

Stop the server:

```
kill -9 $(lsof -i:8087 -t)
```
