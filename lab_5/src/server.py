import asyncio
import math
import websockets
import json

async def calculate_experiment_results(data):
    numbers = data['numbers']
    result = []
    average = sum(numbers) / len(numbers)
    result.append(f"Average: {average}")

    sum_squares = sum((num - average) ** 2 for num in numbers)
    mean_square_deviation = math.sqrt(sum_squares / len(numbers))
    result.append(f"Standard deviation: {mean_square_deviation}")

    relative_error = (mean_square_deviation / average) * 100
    result.append(f"Relative error: {relative_error}")

    return result

async def server(websocket, path):
    while True:
        try:
            data = await websocket.recv()
            print(f"Received data: {data}")

            data = json.loads(data)

            result = await calculate_experiment_results(data)

            await websocket.send(json.dumps(result, ensure_ascii=False).encode('utf-8'))
        except websockets.exceptions.ConnectionClosedError:
            break
print("IP:", end=" ")
ip = input()
print("Port:", end=" ")
port = input()
start_server = websockets.serve(server, ip, port)

asyncio.get_event_loop().run_until_complete(start_server)
asyncio.get_event_loop().run_forever()
