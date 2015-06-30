package raft


import (
    "net"
    "fmt"
    "log"
    "encoding/gob"
)


type Radio struct {
    NodeId int
    Port int
    Comm chan *Msg
}


func NewRadio(id int) *Radio {
    return &Radio{
        NodeId: id,
        Port: 2000 + id,
        Comm: make(chan *Msg, 15),
    }
}


func (self *Radio) Listen() {  
    l, err := net.Listen("tcp", fmt.Sprintf(":%v", self.Port))
    if err != nil {
        log.Fatal("Listen Error: ", err)
    }
    defer l.Close()

    if debug {
        log.Printf("Listening on port %v as node %v...\n", self.Port, self.NodeId)
    }

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal("Accept Error: ", err)
        }

        go self.Receive(conn)
    }
}


func (self *Radio) Receive(c net.Conn) error {
    dec := gob.NewDecoder(c)
    
    var msg *Msg
    err := dec.Decode(&msg)
    if err != nil {
        log.Println("Decode Error: Failed to decode message.", err)
    }

    self.Comm <- msg

    return err
}


func (self *Radio) Send(msg *Msg, targetId int) {
    c, err := net.Dial("tcp", fmt.Sprintf(":%v", 2000 + targetId))
    if err != nil {
        if debug {
            log.Printf("Send Error %v: Unable to connect to node %v\n", self.NodeId, targetId)
        }
        
        return
    }
    defer c.Close()

    enc := gob.NewEncoder(c)
    err = enc.Encode(msg)
    if err != nil {
        log.Println("Node %v Send Error: Failed to encode/send to node", self.NodeId, targetId)
    }
}
