#!/bin/python
# -*- coding: utf-8 -*-
# author:  gagral.x@gmail.com
# time: 2024/1/11 21:45
# description: This script is used to receive alert struct from AlertManager

from flask import Flask, request, jsonify
import json

app = Flask(__name__)

@app.route("/api/receive", methods=["POST"])
def receive_handler():

    """
     Receive alert struct from AlertManager
     @return: Alert data struct
     """

    data = request.get_json()
    print("Alert data struct:\n", data)

    resp_data = {
        "token": "3d7f1f8e5a2c4468b9f4661290c6c615",
    }

    print(json.dumps(resp_data, indent=4, sort_keys=True))

    return jsonify(resp_data)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=18080)
