import requests

# Get the total number of tips
total_tips_response = requests.get("https://neovim.tips/api/total")
if total_tips_response.status_code != 200:
    print("Failed to retrieve the total number of tips")
    exit(1)

# Convert the response text to an integer
try:
    total_tips = int(total_tips_response.text)
except ValueError:
    print("Invalid response for the total number of tips")
    exit(1)

# Loop through tips 1 to total_tips
for i in range(1, total_tips + 1):
    print(f"{i}) ", end="")

    # Run request for each tip and output the result
    tip_response = requests.get(f"https://neovim.tips/api/{i}")
    if tip_response.status_code == 200:
        print(tip_response.text)
    else:
        print("Error retrieving tip")
