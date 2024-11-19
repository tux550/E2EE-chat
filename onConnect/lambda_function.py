# Type: AWS Lambda Function
# Add connectionId to DynamoDB table on connection

import os
import json
import time
import boto3 # type: ignore

ddb = boto3.resource('dynamodb', region_name=os.environ['AWS_REGION'])
table_conn = ddb.Table(os.environ['DDB_TABLE_CONN'])

def lambda_handler(event, context):

    # Validate Username in header
    if "Username" not in event['headers']:
        return {
            'statusCode': 400,
            'body': "Username not found in headers"
        }
    
    # Get connectionId and username
    connectionId = event['requestContext']['connectionId']
    username = event['headers']['Username']
    print(f"ConnectionId: {connectionId}, Username: {username}")
    # Register new connection with lifetime of 2 hours
    try:
        table_conn.put_item(Item={
            'connectionId': connectionId,
            'username': username,
            'ttl': int(time.time()) + 7200
        })
    except Exception as e:
        return {
            'statusCode': 500,
            'body': f"Failed to connect: {json.dumps(str(e))}"
        }

    return {"statusCode":200}