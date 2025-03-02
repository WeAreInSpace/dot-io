package connection

// Client Side

// Connection Request

/*
	Connection Field {
		ClientConnectionHeader: json 1
		Status: json 2
	}
*/

type ClientConnectionHeader struct {
	ProtocolVersion int8                 `json:"proto_ver"` //1.0
	Authentication  ClientAuthentication `json:"auth"`

	PublicKey string `json:"pub_key"`
}

type ClientAuthentication struct {
	JWT    string `json:"jwt"`
	Bearer string `json:"bearer"`
}

// Server Side

// Connection Response

/*
	Connection Field {
		ServerConnectionHeader: json 1
		Status: json 2
	}
*/

type ServerConnectionHeader struct {
	ConnectionUUID      string `json:"uuid"`
	ConnectionPublicKey string `json:"pub_key"`

	PipeAddress string `json:"pipe_addr"`
}

type Status struct {
	// 1 Ok
	//
	// 2 Client Error
	//
	// 3 Server Error
	Code int8 `json:"code"`

	// 1 Client
	//
	// 2 Server
	About int8 `json:"about"`

	// 1 Pass
	//
	// 2 Unauthorized
	//
	// 2 Bad connection
	Info int8 `json:"info"`
}
