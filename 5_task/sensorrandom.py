import random
from flask import Flask, jsonify, request

app = Flask(__name__)

@app.get("/temperature")
def temperature():
    if "location" not in request.args:
        return "location required", 400
    return jsonify(value=random.randint(-20, 40))

@app.get("/temperature/<sensor_id>")
def temperature_by_id(sensor_id):
    return jsonify(value=random.randint(-20, 40))

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8081)
