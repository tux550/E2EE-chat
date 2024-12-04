import asyncio
import websockets

# WebSocket server URL
WEBSOCKET_URL = "ws://example.com:12345"  # Replace with your WebSocket URL

MESSAGE_COUNT = 100 # Number of messages to send

async def send_messages():
    async with websockets.connect(WEBSOCKET_URL) as websocket:
        print(f"Connected to WebSocket server at {WEBSOCKET_URL}")
        
        counter = 1  # Message counter
        
        try:
            while counter <= MESSAGE_COUNT:
                # Create a message
                message = f"Message {counter}"
                print(f"Sending: {message}")
                
                # Send the message
                await websocket.send(message)
                
                # Wait for a response (optional)
                response = await websocket.recv()
                print(f"Received: {response}")
                
                # Increment the counter and wait before sending the next message
                counter += 1
                await asyncio.sleep(1)  # Send a message every second
        except KeyboardInterrupt:
            print("Stopped by user")
        except Exception as e:
            print(f"An error occurred: {e}")

# Entry point for the script
if __name__ == "__main__":
    asyncio.run(send_messages())
