package server

import (
	"strings"
)

type Server struct{}

func (Server) I_PosInvP(m string) string {
	return strings.ToUpper(m)
}

func (s Server) Loop(invP, terP, i_PosInvP chan string) {
	var msgReq, msgRep string
	for {
		select {
		case msgReq = <-invP:
		case <-i_PosInvP:
			msgRep = s.I_PosInvP(msgReq)
		case terP <- msgRep:
		}
	}
}
