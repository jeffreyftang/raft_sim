package raft

import (
    "fmt"
)

type MsgCode string


const (

)


type Msg struct {
    SenderId int
    Code MsgCode 
}


func (self *Msg) String() string {
    return fmt.Sprintf("%v from %v", self.Code, self.SenderId)
}