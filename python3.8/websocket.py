#!/usr/bin/python3
# *_* coding= utf-8 *_*

import asyncio
import websockets
import ssl
import certifi
import json
import logging
import time
import hashlib

WS_URL = "wss://socket.coinex.com/"
access_id = "XXXXX"
secret_key = "XXXXX"

ssl_context = ssl.create_default_context()
ssl_context.load_verify_locations(certifi.where())


async def ping(conn):
    param = {
        "id": 1,
        "method": "server.ping",
        "params": []
    }
    while True:
        await conn.send(json.dumps(param))
        await asyncio.sleep(3)


async def auth(conn):
    current_time = int(time.time()*1000)
    sign_str = f"access_id={access_id}&tonce={current_time}&secret_key={secret_key}"
    md5 = hashlib.md5(sign_str.encode())
    param = {
        "id": 1,
        "method": "server.sign",
        "params": [access_id, md5.hexdigest().upper(), current_time]
    }
    await conn.send(json.dumps(param))   
    res = await conn.recv()
    print(json.loads(res)) 


async def subscribe_depth(conn):
    param = {
        "id": 1,
        "method": "depth.subscribe",
        "params": ['BTCUSDT', 10, "1"]        
    }
    await conn.send(json.dumps(param))
    res = await conn.recv()
    print(json.loads(res)) 


async def subscribe_asset(conn):
    param = {
        "id": 1,
        "method": "asset.subscribe",
        "params": ['USDT', 'CET']        
    }
    await conn.send(json.dumps(param))
    res = await conn.recv()
    print(json.loads(res)) 


async def main():
    # 非标准websocket服务需要关闭客户端的自动ping，设置ping_interval=None
    conn = await websockets.connect(uri=WS_URL, compression=None, ping_interval=None)
    # auth
    await auth(conn)
    # subscribe depth
    await subscribe_depth(conn)
    # subscribe asset
    await subscribe_asset(conn)

    # note: important to keepalive 
    asyncio.create_task(ping(conn))

    # loop to process update data
    while True:
        res = await conn.recv()
        res = json.loads(res)
        print(res)


if __name__ == "__main__":
    # logging.getLogger("asyncio").setLevel(logging.DEBUG)
    loop = asyncio.get_event_loop()
    asyncio.run(main())
    # loop.set_debug(enabled=True)
    print("websocket demo")
    