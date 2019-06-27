package mimic

type response struct {
	Body       interface{} `json:"body"`
	StatusCode int         `json:"status_code"`
}

type responses struct {
	called int
	list   []response
}

func (rs *responses) Get() (resp response, ok bool) {
	if rs == nil || rs.called >= len(rs.list) {
		return resp, false
	}
	resp = rs.list[rs.called]
	rs.called++
	return resp, true
}

func (rs *responses) Set(resp response) {
	// Do not check for nil; allow panic to have helpful stack trace.
	rs.list = append(rs.list, resp)
}
