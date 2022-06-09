#!/usr/bin/python

import json
import requests

headers =  {"Content-Type":"application/json", "X-Request-Datacatalog-Write-Cred": "12345678"}

with open("asset.json", "r") as f:
    data = f.read()

data = json.loads(data)
response = requests.post("http://localhost:8080/createAsset", json=data, headers=headers)
print(response.text)

assetID = response.json()["assetID"]
print("Created Asset " + assetID)

print("Now let us read that asset")

data = {"assetID": assetID, "operationType": "read"}
response = requests.post("http://localhost:8080/getAssetInfo", json=data, headers=headers)
print(response.text)

data = {"assetID": assetID}
response = requests.delete("http://localhost:8080/deleteAsset", json=data, headers=headers)
print(response.text)

