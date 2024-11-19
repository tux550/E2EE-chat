from flask import Flask, request, jsonify
import logging
import boto3

app = Flask(__name__)
app_version = '1.0.0'

api_client = boto3.client('apigatewaymanagementapi', endpoint_url="https://ur6gzg7x79.execute-api.us-east-2.amazonaws.com/production")

logging.basicConfig(level=logging.INFO)
logging.info(f"Starting echo app version {app_version}")

@app.route('/echo', methods=['POST'])
def echo():
    # Get raw body
    raw_data = request.data
    logging.info(f" RAW DATA-{raw_data}")
    # Get headers
    headers = request.headers
    logging.info(f" HEADERS-{headers}")

    # Get json data

    data = request.get_json()
    if not data:
        logging.info(f" ERROR- No json data provided")
        return jsonify({'error': 'No data provided'}), 400
    # Log data
    logging.info(f" DATA-{data}")
    # Send data to the client
    if 'connectionId' in data:
        connection_id = data['connectionId']
        mymessage = f"DATA-{data}\nHEADERS-{headers}"
        try:
            logging.info(f" SENDING TO CONNECTION-{connection_id}")
            response = api_client.post_to_connection(
                ConnectionId=connection_id,
                Data=mymessage.encode('utf-8')
            )
            if response['ResponseMetadata']['HTTPStatusCode'] != 200:
                logging.info(f" ERROR- {response}")
            logging.info(f" SUCCESS SENT TO CONNECTION-{connection_id}")
        except Exception as e:
            logging.info(f" ERROR- {e}")
    else:
        logging.info(f" ERROR- No connectionId provided")

    return jsonify({'echo': data}), 200

# Health check endpoint
@app.route('/')
def health_check():
    return 'OK', 200

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=8080)