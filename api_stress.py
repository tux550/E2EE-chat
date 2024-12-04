import time
import requests
import json

CONNECTION_ID = "placeholder"
API_URL = "http://example.com/api"

# Configuration
wait_interval = 0.5  # Interval between requests in seconds
request_count = 100  # Total number of requests to send
json_payload = {
    "connectionId": CONNECTION_ID,
    "body": {"method":"upload_bundle","params":{"user_id":"user1","bundle":{"identity_key":{"identity_key":"aMOejY1459h1kYeWru7E92Y+a7EDVDdxsNMxA26JkTE="},"signed_pre_key":{"key":"08TKYDKu6naDU9GcH2n1ssylgL9bPgQbWRp1WRqSBQA=","signature":"HgpqxkuYe5QW1J/xGXc4k2r0HRPo1XLcLN/6HDLASYJyH2NxKyKXF2FbOP52912yA/uzYc9qU/tUscMh8L5QCA=="},"one_time_pre_keys":[{"key":"I69qnvMToo+e6R0/D/6m53gt2HcXOGCoXGoqukD8yQk=","id":0},{"key":"OtKWpXPFNu1IG2N+rqVUdmRp32Qv6i5M0omlq3avrwg=","id":1},{"key":"gkSORxJSgmoT3Azj/R8X6qaX8arEZIyJZAt3wGxpBGI=","id":2},{"key":"E2c5V52FrntpOmrvp1E5YS+d3WisCydmEqLfqz8Msg4=","id":3},{"key":"jn6vWMdI6jENG3s2KfTZwXtZ+VNHwCrV6lxKCl5UMGU=","id":4}]}}}
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
