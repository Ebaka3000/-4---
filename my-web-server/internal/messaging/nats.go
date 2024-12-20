package messaging

import (
    "log"
    "github.com/nats-io/nats.go"
    stan "github.com/nats-io/stan.go"
)

type NATSClient struct {
    conn   stan.Conn
    subject string
}

func NewNATSClient(clusterID, clientID string) (stan.Conn, error) {
    sc, err := stan.Connect(clusterID, clientID)
    if err != nil {
        return nil, err
    }
    return sc, nil
}

func (n *NATSClient) Subscribe(handler func(msg []byte)) error {
    _, err := n.conn.Subscribe(n.subject, func(m *stan.Msg) {
        handler(m.Data)
    })
    return err
}

func (n *NATSClient) Close() {
    n.conn.Close()
}