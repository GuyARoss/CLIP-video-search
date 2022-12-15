package zmqpool

import (
	"container/list"
	"fmt"
	"strings"
	"sync"

	zmq "github.com/pebbe/zmq4"
)

type ZMQPoolImplementation interface {
	Send(string, ...string) (string, error)
}

type zeroSocketInstance struct {
	socket *zmq.Socket
}

type ZeroMQPool struct {
	sockets     []*zeroSocketInstance
	readySocket chan *zeroSocketInstance
	m           *sync.Mutex
	queueMutex  *sync.Mutex
	socketQueue *list.List
}

func serializeOutbound(data []string) string {
	// TODO: ensure that data does not contain any ','
	return strings.Join(data, ",")
}

func (pool *ZeroMQPool) waitOpenSocket(sock chan *zeroSocketInstance) {
	for {
		if pool.socketQueue.Len() > 0 {
			e := pool.socketQueue.Front()
			if e.Value != nil {
				sock <- e.Value.(*zeroSocketInstance)
			}
			pool.queueMutex.Lock()
			pool.socketQueue.Remove(e)
			pool.queueMutex.Unlock()
		}
	}
}

func (pool *ZeroMQPool) socketDone(s *zeroSocketInstance) {
	pool.queueMutex.Lock()
	pool.socketQueue.PushBack(s)
	pool.queueMutex.Unlock()
}

func (pool *ZeroMQPool) Send(operation string, messages ...string) (string, error) {
	pool.m.Lock()
	inst := <-pool.readySocket
	pool.m.Unlock()

	_, err := inst.socket.Send(fmt.Sprintf("%s,%s", operation, serializeOutbound(messages)), 0)
	if err != nil {
		return "", err
	}

	response, err := inst.socket.Recv(0)
	if err != nil {
		return "", err
	}
	pool.socketDone(inst)

	if len(response) >= 4 && response[:4] == "error" {
		return "", fmt.Errorf(response[5:])
	}

	return response, nil
}

func New(addresses []string) ZMQPoolImplementation {
	zctx, _ := zmq.NewContext()

	sockets := make([]*zeroSocketInstance, len(addresses))
	sq := list.New()
	for idx, addr := range addresses {
		s, _ := zctx.NewSocket(zmq.REQ)
		s.Connect(addr)

		t := &zeroSocketInstance{
			socket: s,
		}
		sockets[idx] = t
		sq.PushBack(t)
	}

	socketChan := make(chan *zeroSocketInstance)

	pool := &ZeroMQPool{
		sockets:     sockets,
		readySocket: socketChan,
		m:           &sync.Mutex{},
		socketQueue: sq,
		queueMutex:  &sync.Mutex{},
	}
	go pool.waitOpenSocket(socketChan)

	return pool
}
