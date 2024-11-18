import asyncio
import websockets
import httpx
import json
from websockets.exceptions import ConnectionClosed

# Define ports
WEBSOCKET_PORT = 8081
HTTP_FORWARD_PORT = 8080
HTTP_LISTEN_PORT = 8082

# Store active WebSocket connections
connections = {}

async def forward_to_http(message: dict):
    """Forward a WebSocket message to the HTTP server."""
    async with httpx.AsyncClient() as client:
        try:
            response = await client.post(f"http://localhost:{HTTP_FORWARD_PORT}", json=message)
            print(f"Forwarded to HTTP: {message}, HTTP Response: {response.status_code}")
        except Exception as e:
            print(f"Error forwarding to HTTP: {e}")

async def websocket_handler(websocket, path):
    """Handle incoming WebSocket connections."""
    connection_id = id(websocket)  # Use Python object ID as a unique connection ID
    connections[connection_id] = websocket
    print(f"Connection opened with ID: {connection_id}")

    try:
        async for message in websocket:
            print(f"Received from WebSocket (ID {connection_id}): {message}")
            # Create the template message
            payload = {
                "connectionId": str(connection_id),
                "body": message,
            }
            # Forward to HTTP
            await forward_to_http(payload)
    except ConnectionClosed:
        print(f"Connection closed for ID: {connection_id}")
    finally:
        connections.pop(connection_id, None)

async def handle_http_request(reader, writer):
    """Handle incoming HTTP POST requests."""
    data = await reader.read(1024)
    try:
        # Parse the received data
        request = json.loads(data.decode())
        connection_id = int(request["connectionId"])
        body = request["body"]
        print(f"Received from HTTP: {request}")

        # Find the WebSocket connection
        websocket = connections.get(connection_id)
        if websocket:
            # Send the message to the WebSocket client
            await websocket.send(body)
            response = {"status": "success", "message": "Sent to WebSocket"}
        else:
            response = {"status": "error", "message": "Connection ID not found"}
    except Exception as e:
        response = {"status": "error", "message": str(e)}

    # Respond to the HTTP client
    writer.write(f"HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{json.dumps(response)}".encode())
    await writer.drain()
    writer.close()

async def http_server():
    """Start the HTTP listener."""
    server = await asyncio.start_server(handle_http_request, "localhost", HTTP_LISTEN_PORT)
    print(f"HTTP server is running on http://localhost:{HTTP_LISTEN_PORT}")
    async with server:
        await server.serve_forever()

async def main():
    """Start the WebSocket and HTTP servers."""
    websocket_server = websockets.serve(websocket_handler, "localhost", WEBSOCKET_PORT)
    print(f"WebSocket server is running on ws://localhost:{WEBSOCKET_PORT}")

    await asyncio.gather(websocket_server, http_server())

if __name__ == "__main__":
    asyncio.run(main())
