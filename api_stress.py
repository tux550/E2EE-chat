import time
import requests
import json

CONNECTION_ID = "placeholder"
TARGET_USER_ID = "u3"
API_URL = "http://my-docker-app-lb-1325993742.us-east-2.elb.amazonaws.com/"

# Configuration
wait_interval = 0.5  # Interval between requests in seconds
request_count = 100  # Total number of requests to send
json_payload = {
    "connectionId": CONNECTION_ID,
    "body": {"method":"get_bundle","params":{"user_id":TARGET_USER_ID}}
}

def send_request():
    try:
        response = requests.post(API_URL, json=json_payload)
        print(f"Status Code: {response.status_code}")
        print(f"Response: {response.text}")
    except requests.exceptions.RequestException as e:
        print(f"Request failed: {e}")

def main():
    print(f"Starting requests to {API_URL} with an interval of {wait_interval} seconds...")
    print(f"Total requests to send: {request_count}")
    for i in range(request_count):
        print(f"Sending request {i + 1} of {request_count}...")
        send_request()
        if i < request_count - 1:  # Avoid waiting after the last request
            time.sleep(wait_interval)
    print("All requests completed.")

if __name__ == "__main__":
    main()
