package beacon

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/basicrum/front_basicrum_go/geoip"
	"github.com/basicrum/front_basicrum_go/types"
	ua "github.com/mileusna/useragent"
	"github.com/ua-parser/uap-go/uaparser"
)

// Beacon contains the performance statistics from request
type Beacon struct {
	// Mobile
	Mob_Etype string
	Mob_Dl    string
	Mob_Rtt   string

	// Paint Timing
	Pt_Lcp string
	Pt_Fp  string
	Pt_Fcp string

	// Continuity
	C_E      string
	C_L      string
	C_Lb     string
	C_Tti_M  string
	C_Tti    string
	C_Tti_Vr string
	C_T_Fps  string
	C_F      string
	C_F_D    string
	C_F_M    string
	C_F_L    string
	C_F_S    string
	C_Cls    string
	C_Fid    string

	// Event Timing
	Et_Fid string
	Et_E   string

	// Roundtrip
	Rt_Start   string
	Rt_Bmr     string
	Rt_Tstart  string
	Rt_Bstart  string
	Rt_Blstart string
	Rt_End     string
	Rt_Tt      string
	Rt_Obo     string
	Rt_Si      string
	Rt_Ss      string
	Rt_Sl      string
	Rt_Quit    bool

	// Misc
	U              string
	T_Resp         string
	T_Page         string
	T_Done         string
	T_Other        string
	V              string
	Restiming      string
	CreatedAt      string
	Sv             string
	Sm             string
	Vis_St         string
	Ua_Plt         string
	Ua_Vnd         string
	Pid            string
	N              string
	Http_Initiator string

	// Navigation Timing
	Nt_Nav_St            string
	Nt_Fet_St            string
	Nt_Dns_St            string
	Nt_Dns_End           string
	Nt_Con_St            string
	Nt_Con_End           string
	Nt_Req_St            string
	Nt_Res_St            string
	Nt_Res_End           string
	Nt_Domloading        string
	Nt_Domint            string
	Nt_Domcontloaded_St  string
	Nt_Domcontloaded_End string
	Nt_Domcomp           string
	Nt_Load_St           string
	Nt_Load_End          string
	Nt_Unload_St         string
	Nt_Unload_End        string
	Nt_Ssl_St            string
	Nt_Enc_Size          string
	Nt_Dec_Size          string
	Nt_Trn_Size          string
	Nt_Protocol          string
	Nt_First_Paint       string
	Nt_Red_Cnt           string
	Nt_Nav_Type          string

	// Memory
	Dom_Res        string
	Dom_Doms       string
	Mem_Total      string
	Mem_Limit      string
	Mem_Used       string
	Mem_Lsln       string
	Mem_Ssln       string
	Mem_Lssz       string
	Mem_Sssz       string
	Scr_Xy         string
	Scr_Bpp        string
	Scr_Orn        string
	Cpu_Cnc        string
	Dom_Ln         string
	Dom_Sz         string
	Dom_Ck         string
	Dom_Img        string
	Dom_Img_Uniq   string
	Dom_Script     string
	Dom_Script_Ext string
	Dom_Iframe     string
	Dom_Link       string
	Dom_Link_Css   string
	Net_Sd         string
	Sb             string
}

// FromEvent creates Beacon request from http request parameters
// nolint: funlen
func FromEvent(event *types.Event) Beacon {
	values := event.RequestParameters
	return Beacon{
		// Used constructing event date
		CreatedAt: values.Get("created_at"),

		// Mobile
		Mob_Etype: values.Get("mob.etype"),
		Mob_Dl:    values.Get("mob.dl"),
		Mob_Rtt:   values.Get("mob.rtt"),

		// Paint Timing
		Pt_Lcp: values.Get("pt.lcp"),
		Pt_Fp:  values.Get("pt.fp"),
		Pt_Fcp: values.Get("pt.fcp"),

		// Continuity
		C_E:      values.Get("c.e"),
		C_Tti_M:  values.Get("c.tti.m"),
		C_T_Fps:  values.Get("c.t.fps"),
		C_Tti_Vr: values.Get("c.tti.vr"),
		C_Tti:    values.Get("c.tti"),
		C_F:      values.Get("c.f"),
		C_F_D:    values.Get("c.f.d"),
		C_F_M:    values.Get("c.f.m"),
		C_F_S:    values.Get("c.f.s"),
		C_Fid:    values.Get("c.fid"),
		C_Cls:    values.Get("c.cls"),

		// Event Timing
		Et_Fid: values.Get("et.fid"),
		Et_E:   values.Get("et.e"),

		// Misc
		U:              values.Get("u"),
		Restiming:      values.Get("restiming"),
		T_Resp:         values.Get("t_resp"),
		T_Page:         values.Get("t_page"),
		T_Done:         values.Get("t_done"),
		T_Other:        values.Get("t_other"),
		V:              values.Get("v"),
		Sv:             values.Get("sv"),
		Sm:             values.Get("sm"),
		Vis_St:         values.Get("vis.st"),
		Ua_Plt:         values.Get("ua.plt"),
		Ua_Vnd:         values.Get("ua.vnd"),
		Pid:            values.Get("pid"),
		N:              values.Get("n"),
		Http_Initiator: values.Get("http_initiator"),

		// Navigation Timing
		Nt_Nav_St:            values.Get("nt_nav_st"),
		Nt_Fet_St:            values.Get("nt_fet_st"),
		Nt_Dns_St:            values.Get("nt_dns_st"),
		Nt_Dns_End:           values.Get("nt_dns_end"),
		Nt_Con_St:            values.Get("nt_con_st"),
		Nt_Con_End:           values.Get("nt_con_end"),
		Nt_Req_St:            values.Get("nt_req_st"),
		Nt_Res_St:            values.Get("nt_res_st"),
		Nt_Res_End:           values.Get("nt_res_end"),
		Nt_Domloading:        values.Get("nt_domloading"),
		Nt_Domint:            values.Get("nt_domint"),
		Nt_Domcontloaded_St:  values.Get("nt_domcontloaded_st"),
		Nt_Domcontloaded_End: values.Get("nt_domcontloaded_end"),
		Nt_Domcomp:           values.Get("nt_domcomp"),
		Nt_Load_St:           values.Get("nt_load_st"),
		Nt_Load_End:          values.Get("nt_load_end"),
		Nt_Unload_St:         values.Get("nt_unload_st"),
		Nt_Unload_End:        values.Get("nt_unload_end"),
		Nt_Ssl_St:            values.Get("nt_ssl_st"),
		Nt_Enc_Size:          values.Get("nt_enc_size"),
		Nt_Dec_Size:          values.Get("nt_dec_size"),
		Nt_Trn_Size:          values.Get("nt_trn_size"),
		Nt_Protocol:          values.Get("nt_protocol"),
		Nt_First_Paint:       values.Get("nt_first_paint"),
		Nt_Red_Cnt:           values.Get("nt_red_cnt"),
		Nt_Nav_Type:          values.Get("nt_nav_type"),

		Rt_Start:   values.Get("navigation"),
		Rt_Bmr:     values.Get("rt.bmr"),
		Rt_Tstart:  values.Get("rt.tstart"),
		Rt_Bstart:  values.Get("rt.bstart"),
		Rt_Blstart: values.Get("rt.blstart"),
		Rt_End:     values.Get("rt.end"),
		Rt_Tt:      values.Get("rt.tt"),
		Rt_Obo:     values.Get("rt.obo"),
		Rt_Si:      values.Get("rt.si"),
		Rt_Ss:      values.Get("rt.ss"),
		Rt_Sl:      values.Get("rt.sl"),
		Rt_Quit:    values.Has("rt.quit"),

		// Memory
		Dom_Res:        values.Get("dom.res"),
		Dom_Doms:       values.Get("dom.doms"),
		Mem_Total:      values.Get("mem.total"),
		Mem_Limit:      values.Get("mem.limit"),
		Mem_Used:       values.Get("mem.used"),
		Mem_Lsln:       values.Get("mem.lsln"),
		Mem_Ssln:       values.Get("mem.ssln"),
		Mem_Lssz:       values.Get("mem.lssz"),
		Mem_Sssz:       values.Get("mem.sssz"),
		Scr_Xy:         values.Get("scr.xy"),
		Scr_Bpp:        values.Get("scr.bpp"),
		Scr_Orn:        values.Get("scr.orn"),
		Cpu_Cnc:        values.Get("cpu.cnc"),
		Dom_Ln:         values.Get("dom.ln"),
		Dom_Sz:         values.Get("dom.sz"),
		Dom_Ck:         values.Get("dom.ck"),
		Dom_Img:        values.Get("dom.img"),
		Dom_Img_Uniq:   values.Get("dom.img.uniq"),
		Dom_Script:     values.Get("dom.script"),
		Dom_Script_Ext: values.Get("dom.script.ext"),
		Dom_Iframe:     values.Get("dom.iframe"),
		Dom_Link:       values.Get("dom.link"),
		Dom_Link_Css:   values.Get("dom.link.css"),
		Net_Sd:         values.Get("net.sd"),
		Sb:             values.Get("sb"),
	}
}

// ConvertToRumEvent convert Beacon request to Rum Event
func ConvertToRumEvent(b Beacon, event *types.Event, userAgentParser *uaparser.Parser, geoIPService geoip.Service) RumEvent {
	userAgent := event.UserAgent

	userAgentClient := userAgentParser.Parse(userAgent)

	deviceType := getDeviceType(userAgent)

	screenWidth, screenHeight := getScreenSize(b.Scr_Xy)

	urlValue, err := url.Parse(b.U)
	if err != nil {
		log.Println(err)
	}

	hostname := urlValue.Hostname()

	var country, city string
	if geoIPService != nil {
		country, city, _ = geoIPService.CountryAndCity(event.Headers, event.RemoteAddr)
	}

	return RumEvent{
		Created_At:               b.CreatedAt,
		Hostname:                 hostname,
		Url:                      b.U,
		Cumulative_Layout_Shift:  json.Number(b.C_Cls),
		Device_Type:              deviceType,
		Device_Manufacturer:      userAgentClient.Device.Brand,
		Operating_System:         userAgentClient.Os.Family,
		Operating_System_Version: userAgentClient.Os.ToVersionString(),
		Browser_Name:             userAgentClient.UserAgent.Family,
		Browser_Version:          userAgentClient.UserAgent.ToVersionString(),
		Connect_Duration:         calculateDelta(b.Nt_Con_St, b.Nt_Con_End),
		Dns_Duration:             calculateDelta(b.Nt_Dns_St, b.Nt_Dns_End),
		First_Byte_Duration:      calculateDelta(b.Nt_Nav_St, b.Nt_Res_St),
		Redirect_Duration:        "0", // @todo: Calculate later
		Redirects_Count:          "0", // @todo: Calculate later
		First_Contentful_Paint:   b.Pt_Fcp,
		First_Paint:              b.Pt_Fp,
		First_Input_Delay:        json.Number(b.Et_Fid),
		Largest_Contentful_Paint: b.Pt_Lcp,
		Event_Type:               getEventType(b.Rt_Quit, b.Http_Initiator),
		Session_Id:               b.Rt_Si,
		Session_Length:           b.Rt_Sl,
		Geo_Country_Code:         country,
		Geo_City_Name:            city,
		Next_Hop_Protocol:        b.Nt_Protocol,
		User_Agent:               userAgent,
		Visibility_State:         b.Vis_St,
		Boomerang_Version:        b.V,
		Screen_Width:             screenWidth,
		Screen_Height:            screenHeight,
		Dom_Res:                  b.Dom_Res,
		Dom_Doms:                 b.Dom_Doms,
		Mem_Total:                b.Mem_Total,
		Mem_Limit:                b.Mem_Limit,
		Mem_Used:                 b.Mem_Used,
		Mem_Lsln:                 b.Mem_Lsln,
		Mem_Ssln:                 b.Mem_Ssln,
		Mem_Lssz:                 b.Mem_Lssz,
		Scr_Bpp:                  b.Scr_Bpp,
		Scr_Orn:                  b.Scr_Orn,
		Cpu_Cnc:                  b.Cpu_Cnc,
		Dom_Ln:                   b.Dom_Ln,
		Dom_Sz:                   b.Dom_Sz,
		Dom_Ck:                   b.Dom_Ck,
		Dom_Img:                  b.Dom_Img,
		Dom_Img_Uniq:             b.Dom_Img_Uniq,
		Dom_Script:               b.Dom_Script,
		Dom_Script_Ext:           b.Dom_Script_Ext,
		Dom_Iframe:               b.Dom_Iframe,
		Dom_Link:                 b.Dom_Link,
		Dom_Link_Css:             b.Dom_Link_Css,
		Page_Id:                  b.Pid,
		Ua_Vnd:                   b.Ua_Vnd,
		Ua_Plt:                   b.Ua_Plt,
		Data_Saver_On:            json.Number(b.Net_Sd),
		Mob_Etype:                b.Mob_Etype,
		Mob_Dl:                   json.Number(b.Mob_Dl),
		Mob_Rtt:                  json.Number(b.Mob_Rtt),
	}
}

func calculateDelta(p1 string, p2 string) string {
	if p2 == "" || p1 == "" {
		return ""
	}

	// Values from Navigation Timings plugin
	end, _ := strconv.Atoi(p2)
	start, _ := strconv.Atoi(p1)

	v := end - start

	if v < 0 {
		v = 0
	}

	if v > 65535 {
		v = 65535
	}

	return strconv.Itoa(v)
}

func getDeviceType(uagent string) string {
	mileusnaUaRes := ua.Parse(uagent)

	dT := "unknown"

	if mileusnaUaRes.Mobile {
		dT = "mobile"
	}
	if mileusnaUaRes.Tablet {
		dT = "tablet"
	}
	if mileusnaUaRes.Desktop {
		dT = "desktop"
	}
	if mileusnaUaRes.Bot {
		dT = "bot"
	}

	return dT
}

// nolint: revive
func getScreenSize(scr_X_Y string) (string, string) {
	s := strings.Split(scr_X_Y, "x")

	if len(s) == 2 {
		return s[0], s[1]
	}

	return "", ""
}

// nolint: revive
func getEventType(isQuit bool, httpInitiator string) string {
	if len(httpInitiator) > 0 {
		return httpInitiator
	}

	if isQuit {
		return "quit_page"
	}

	return "visit_page"
}
