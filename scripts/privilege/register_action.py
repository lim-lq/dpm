import json

import requests

url = "http://localhost:8890/api/common/register/action"


def main():
    headers = {
        "Content-Type": "application/json"
    }
    data = json.load(open("./actions.json"))
    resp = requests.post(url, data=json.dumps(data), headers=headers)
    print(resp.content)
    
    
if __name__ == "__main__":
    main()