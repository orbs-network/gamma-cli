# Copyright 2019 the gamma-cli authors
# This file is part of the gamma-cli library in the Orbs project.
#
# This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
# The above notice should be included in all copies or substantial portions of the software.

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