# Type: AWS Lambda Function
# Remove from DynamoDB on disconnect

import os
import json
import boto3 # type: ignore

ddb = boto3.resource('dynamodb', region_name=os.environ['AWS_REGION'])
table = ddb.Table(os.environ['DDB_TABLE_NAME'])

def lambda_handler(event, context):
    connectionId = event['requestContext']['connectionId']
    try:
        # Delete item from DynamoDB table
        table.delete_item(Key={'connectionId': connectionId})
    except Exception as e:
        # Log error
        # print("ERROR MSG:", f"{json.dumps(str(e))}")
        return {
            'statusCode': 500,
            'body': f"Failed to disconnect: {json.dumps(str(e))}"
        }

    return {"statusCode":200}