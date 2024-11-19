import asyncio
import httpx
import json
from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from fastapi.responses import JSONResponse
from fastapi import HTTPException

#Backend server ports
HTTP_SERVER_PORT = 8080
#Gateway server ports
API_SERVER_PORT = 8082

HTTP_SERVER_URL = f"http://local-e2ee-app:{HTTP_SERVER_PORT}"

# Data structure to keep track of WebSocket connections and associated usernames
connections_username = {}
connections_objs = {}

# FastAPI application to serve HTTP API
app = FastAPI()

async def forward_to_http(message: dict):
    """Forward a WebSocket message to the HTTP server."""
    async with httpx.AsyncClient() as client:
        try:
            print(f"Forwarding to HTTP: {message}")
            response = await client.post(HTTP_SERVER_URL, json=message)
            print(f"HTTP Response: {response.status_code}")
        except Exception as e:
            print(f"Error forwarding to HTTP: {e}")

@app.websocket("/ws")
async def websocket_handler(websocket: WebSocket):
    """Handle incoming WebSocket connections."""
    await websocket.accept()
    username = None
    connection_id = str(id(websocket))  # Connection ID is a string

    # Extract the username header if available
    for header, value in websocket.headers.items():
        if header.lower() == "username":
            username = value
            break

    if username:
        # Save the connection username and connection
        connections_username[connection_id] = username  
        connections_objs[connection_id] = websocket

    print(f"Connection opened with ID: {connection_id}, Username: {username}")

    try:
        while True:
            message = await websocket.receive_json()
            print(f"Received from WebSocket: {message}")
            # Create the template message
            payload = {
                "connectionId": connection_id,
                "body": message,
            }
            # Forward to HTTP
            await forward_to_http(payload)
    except WebSocketDisconnect:
        print(f"Connection closed for ID: {connection_id}")
        if connection_id in connections_username:
            del connections_username[connection_id]  # Remove connection info when closed

@app.post("/send/{connectionId}")
async def send_message(connectionId: str, body: dict):
    """Receive a message body and forward it to the corresponding WebSocket."""
    if connectionId in connections_username:
        # Get the WebSocket connection object
        conn = connections_objs.get(connectionId)
        # Send message to the WebSocket
        try:
            # Send the message as JSON to the WebSocket
            await conn.send_json(body)
            return JSONResponse(content={"status": "success"})
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Failed to send message: {str(e)}")
    else:
        raise HTTPException(status_code=404, detail="Connection ID not found")

@app.get("/connections/{connectionId}")
async def get_connection(connectionId: str):
    """Return the username for a given connection ID."""
    if connectionId in connections_username:
        return JSONResponse(content={"connectionId": connectionId, "username": connections_username[connectionId]})
    else:
        raise HTTPException(status_code=404, detail="Connection ID not found")

@app.get("/username/{username}")
async def get_connections_by_username(username: str):
    """Return all connections for a given username."""
    user_connections = [
        {"connectionId": cid, "username": username}
        for cid, uname in connections_username.items()
        if uname == username
    ]
    if user_connections:
        return JSONResponse(content=user_connections)
    else:
        raise HTTPException(status_code=404, detail="No connections found for this username")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=API_SERVER_PORT)
