# Type: AWS Lambda Function
# Add connectionId to DynamoDB table on connection

import os
import json
import boto3 # type: ignore

ddb = boto3.resource('dynamodb', region_name=os.environ['AWS_REGION'])
table = ddb.Table(os.environ['DDB_TABLE_NAME'])

def lambda_handler(event, context):
    put_params = {
      'connectionId': event['requestContext']['connectionId']
    }
    try:
        # Put item in DynamoDB table
        table.put_item(Item=put_params)
    except Exception as e:
        # Log error
        # print("ERROR MSG:", f"{json.dumps(str(e))}")
        return {
            'statusCode': 500,
            'body': f"Failed to connect: {json.dumps(str(e))}"
        }

    return {"statusCode":200}