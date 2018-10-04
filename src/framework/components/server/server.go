package server

import (
	"strings"
)

type Server struct{}

func (Server) I_PosInvP(m *string) {
	*m = strings.ToUpper(*m)
}

//func (Receiver) Loop(invR, terR, invP, terP, i_PosInvP chan string) {
func (Server) Loop(invP, terP, i_PosInvP chan string) {
	//msg := "testV1"
	//i := 0
	msgRecv := ""
	for {
		select {
		//case invR <- msg + strconv.Itoa(i):
		//fmt.Println("Receiver:: invR")
		//	i++
		case terP <- msgRecv:
		//fmt.Println("Receiver:: terP")
		//case msgRcv = <-terR:
		//fmt.Println("Receiver:: terR")
		case <-invP:
			//fmt.Println("receiver:: InvP")
		case msgRcv := <-i_PosInvP:
			//fmt.Println("Receiver:: i_PosInvP")
			Server{}.I_PosInvP(&msgRcv)
		}
	}
}