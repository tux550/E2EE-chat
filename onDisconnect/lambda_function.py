# Type: AWS Lambda Function
# Remove from DynamoDB on disconnect

import os
import json
import boto3 # type: ignore

ddb = boto3.resource('dynamodb', region_name=os.environ['AWS_REGION'])
table_conn = ddb.Table(os.environ['DDB_TABLE_CONN'])

def lambda_handler(event, context):

    # Get connectionId
    connectionId = event['requestContext']['connectionId']

    # Remove connection from table
    try:
        table_conn.delete_item(
            Key={
                'connectionId': connectionId
            }
        )
    except Exception as e:
        return {
            'statusCode': 500,
            'body': f"Failed to disconnect: {json.dumps(str(e))}"
        }
    return {"statusCode":200}