package beacon

import "encoding/json"

// RumEvent contains the Rum event data
type RumEvent struct {
	Created_At               string      `json:"created_at"`
	Hostname                 string      `json:"hostname"`
	Url                      string      `json:"url"`
	Cumulative_Layout_Shift  json.Number `json:"cumulative_layout_shift,omitempty"`
	Geo_Country_Code         string      `json:"geo_country_code"`
	Geo_City_Name            string      `json:"geo_city_name"`
	Device_Type              string      `json:"device_type"`
	Device_Manufacturer      string      `json:"device_manufacturer,omitempty"`
	T_Resp                   string      `json:"t_resp"`
	T_Page                   string      `json:"t_page"`
	T_Done                   string      `json:"t_done"`
	Connect_Duration         string      `json:"connect_duration"`
	Ssl_Negotiation_Duration string      `json:"ssl_negotiation_duration"`
	Next_Hop_Protocol        string      `json:"next_hop_protocol"`
	Dns_Duration             string      `json:"dns_duration"`
	First_Byte_Duration      string      `json:"first_byte_duration"`
	Session_Id               string      `json:"session_id"`
	Session_Length           string      `json:"session_length"`
	Operating_System         string      `json:"operating_system"`
	Operating_System_Version string      `json:"operating_system_version,omitempty"`
	Browser_Name             string      `json:"browser_name"`
	Browser_Version          string      `json:"browser_version,omitempty"`
	Event_Type               string      `json:"event_type"`
	Redirect_Duration        string      `json:"redirect_duration"`
	Redirects_Count          string      `json:"redirects_count"`
	First_Contentful_Paint   string      `json:"first_contentful_paint"`
	First_Paint              string      `json:"first_paint"`
	First_Input_Delay        json.Number `json:"first_input_delay,omitempty"`
	Largest_Contentful_Paint string      `json:"largest_contentful_paint"`
	User_Agent               string      `json:"user_agent,omitempty"`
	Visibility_State         string      `json:"visibility_state"`
	Boomerang_Version        string      `json:"boomerang_version"`
	Screen_Width             string      `json:"screen_width"`
	Screen_Height            string      `json:"screen_height"`
	Dom_Res                  string      `json:"dom_res"`
	Dom_Doms                 string      `json:"dom_doms"`
	Mem_Total                string      `json:"mem_total,omitempty"`
	Mem_Limit                string      `json:"mem_limit,omitempty"`
	Mem_Used                 string      `json:"mem_used,omitempty"`
	Mem_Lsln                 string      `json:"mem_lsln"`
	Mem_Ssln                 string      `json:"mem_ssln"`
	Mem_Lssz                 string      `json:"mem_lssz"`
	Scr_Bpp                  string      `json:"scr_bpp"`
	Scr_Orn                  string      `json:"scr_orn"`
	Cpu_Cnc                  string      `json:"cpu_cnc"`
	Dom_Ln                   string      `json:"dom_ln"`
	Dom_Sz                   string      `json:"dom_sz"`
	Dom_Ck                   string      `json:"dom_ck"`
	Dom_Img                  string      `json:"dom_img"`
	Dom_Img_Uniq             string      `json:"dom_img_uniq"`
	Dom_Script               string      `json:"dom_script"`
	Dom_Script_Ext           string      `json:"dom_script_ext"`
	Dom_Iframe               string      `json:"dom_iframe"`
	Dom_Link                 string      `json:"dom_link"`
	Dom_Link_Css             string      `json:"dom_link_css"`
	Page_Id                  string      `json:"page_id"`
	Ua_Vnd                   string      `json:"ua_vnd,omitempty"`
	Ua_Plt                   string      `json:"ua_plt,omitempty"`
	Data_Saver_On            json.Number `json:"data_saver_on,omitempty"`
	Mob_Etype                string      `json:"mob_etype,omitempty"`
	Mob_Dl                   json.Number `json:"mob_dl,omitempty"`
	Mob_Rtt                  json.Number `json:"mob_rtt,omitempty"`
}
