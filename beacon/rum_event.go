package beacon

type RumEvent struct {
	Created_At               string `json:"created_at"`
	Url                      string `json:"url"`
	Cumulative_Layout_Shift  string `json:"cumulative_layout_shift"`
	Country_Code             string `json:"country_code"`
	Device_Type              string `json:"device_type"`
	Device_Manufacturer      string `json:"device_manufacturer"`
	T_Resp                   string `json:"t_resp"`
	T_Page                   string `json:"t_page"`
	T_Done                   string `json:"t_done"`
	Connect_Duration         string `json:"connect_duration"`
	Ssl_Negotiation_Duration string `json:"ssl_negotiation_duration"`
	Next_Hop_Protocol        string `json:"next_hop_protocol"`
	Dns_Duration             string `json:"dns_duration"`
	First_Byte_Duration      string `json:"first_byte_duration"`
	Session_Id               string `json:"session_id"`
	Session_Length           string `json:"session_length"`
	Operating_System         string `json:"operating_system"`
	Browser_Name             string `json:"browser_name"`
	Browser_Version          string `json:"browser_version"`
	Event_Type               string `json:"event_type"`
	Redirect_Duration        string `json:"redirect_duration"`
	Redirects_Count          string `json:"redirects_count"`
	First_Contentful_Paint   string `json:"first_contentful_paint"`
	First_Paint              string `json:"first_paint"`
	First_Input_Delay        string `json:"first_input_delay"`
	Largest_Contentful_Paint string `json:"largest_contentful_paint"`
	User_Agent               string `json:"user_agent"`
	Visibility_State         string `json:"visibility_state"`
}
