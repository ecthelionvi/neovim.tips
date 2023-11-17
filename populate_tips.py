import requests
from dotenv import load_dotenv
import os

load_dotenv()

login_username = os.getenv("SUPERUSER_NAME")
login_password = os.getenv("SUPERUSER_PASSWORD")

login_url = "https://neovim.tips/api/login"

response = requests.post(login_url, json={"username": login_username, "password": login_password})

if response.status_code != 200:
    print("Failed to retrieve token")
    exit(1)

token = response.text

add_tip_url = "https://neovim.tips/api/add"

tips = []

for tip in tips:
    tip_json_data = {"content": tip}

    response = requests.post(add_tip_url, json=tip_json_data, headers={"Authorization": token})

    print(f"Response for adding tip: {tip}")
    print(response.text)
    print("")
