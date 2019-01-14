import json
import sys
import socket

def get_endpoint(ip):
    return "http://" + ip + ":8080"

def get_config(ip):
    return {
        "Environments": {
            "docker": {
                "VirtualChain": 42,
                "Endpoints": [get_endpoint(ip)]
            },
            "docker-experimental": {
                "VirtualChain": 42,
                "Endpoints": [get_endpoint(ip)],
                "Experimental": True
            },
        }
    }

if __name__ == "__main__":
    ip = socket.gethostbyname(socket.gethostname())

    if sys.argv[1] == "endpoint":
        print get_endpoint(ip)
    elif sys.argv[1] == "json":
        print json.dumps(get_config(ip))