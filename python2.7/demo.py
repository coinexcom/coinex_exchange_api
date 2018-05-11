#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
Created by bu on 2018-01-17
"""
from __future__ import unicode_literals
import time
import hashlib
import json as complex_json
import urllib3
from urllib3.exceptions import InsecureRequestWarning

urllib3.disable_warnings(InsecureRequestWarning)
http = urllib3.PoolManager(timeout=urllib3.Timeout(connect=1, read=2))


class RequestClient(object):
    __headers = {
        'Content-Type': 'application/json; charset=utf-8',
        'Accept': 'application/json',
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36'
    }

    def __init__(self, headers={}):
        self.access_id = ''      # replace
        self.secret_key = ''     # replace
        self.url = 'https://api.coinex.com'
        self.headers = self.__headers
        self.headers.update(headers)

    @staticmethod
    def get_sign(params, secret_key):
        sort_params = sorted(params)
        data = []
        for item in sort_params:
            data.append(item + '=' + str(params[item]))
        str_params = "{0}&secret_key={1}".format('&'.join(data), secret_key)
        token = hashlib.md5(str_params).hexdigest().upper()
        return token

    def set_authorization(self, params):
        params['access_id'] = self.access_id
        params['tonce'] = int(time.time()*1000)
        self.headers['AUTHORIZATION'] = self.get_sign(params, self.secret_key)

    def request(self, method, url, params={}, data='', json={}):
        method = method.upper()
        if method in ['GET', 'DELETE']:
            self.set_authorization(params)
            result = http.request(method, url, fields=params, headers=self.headers)
        else:
            if data:
                json.update(complex_json.loads(data))
            self.set_authorization(json)
            encoded_data = complex_json.dumps(json).encode('utf-8')
            result = http.request(method, url, body=encoded_data, headers=self.headers)
        return result


def get_account():
    request_client = RequestClient()
    response = request_client.request('GET', '{url}/v1/balance/'.format(url=request_client.url))
    print response.status


def order_pending(market_type):
    request_client = RequestClient()
    params = {
        'market': market_type
    }
    response = request_client.request(
            'GET',
            '{url}/v1/order/pending'.format(url=request_client.url),
            params=params
    )
    print response.content


def order_finished(market_type, page, limit):
    request_client = RequestClient()
    params = {
        'market': market_type,
        'page': page,
        'limit': limit
    }
    response = request_client.request(
            'GET',
            '{url}/v1/order/finished'.format(url=request_client.url),
            params=params
    )
    print response.status


def put_limit():
    request_client = RequestClient()
    data = {
            "amount": "10",
            "price": "0.0001",
            "type": "sell",
            "market": "CETBCH"
        }

    response = request_client.request(
            'POST',
            '{url}/v1/order/limit'.format(url=request_client.url),
            json=data,
    )
    return response.data


def put_market():
    request_client = RequestClient()

    data = {
            "amount": "1",
            "type": "sell",
            "market": "CETBCH"
        }

    response = request_client.request(
            'POST',
            '{url}/v1/order/market'.format(url=request_client.url),
            json=data,
    )
    print response.content


def cancel_order(id, market):
    request_client = RequestClient()
    data = {
        "id": id,
        "market": market,
    }
    print market

    response = request_client.request(
            'DELETE',
            '{url}/v1/order/pending'.format(url=request_client.url),
            params=data,
    )
    return response.data


if __name__ == '__main__':
    count = 1
    a = time.time() * 1000
    while True:
        b = time.time() * 1000
        order_data = complex_json.loads(put_limit())['data']
        id = order_data['id']
        market = order_data['market']
        cancel_order(id, market)
        print time.time() * 1000 - b
        count += 1
        if count >= 50:
            break

    print 'avg', (time.time() * 1000 - a) / 50.0
