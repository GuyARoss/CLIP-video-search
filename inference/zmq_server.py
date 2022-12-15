import zmq
import threading

from zmq_ops import registry

def build_worker(port: int):
    def worker():
        print('started worker', port)
        context = zmq.Context()
        socket = context.socket(zmq.REP)
        socket.bind(f"tcp://*:{str(port)}")

        while True:
            message = socket.recv()        
            message:str = message.decode('utf-8')

            split_message = message.split(',')
            type = split_message[0]
            
            response = registry[type](*split_message[1:])
            socket.send(bytes(response,'utf-8'))

    return worker


def main() -> int:    
    free_ports = [5550, 5551, 5552]

    threads = [threading.Thread(target=build_worker(free_ports[i])) for i in range(len(free_ports))]
    for thread in threads:
        thread.start()

if __name__ == '__main__':
    SystemExit(main())