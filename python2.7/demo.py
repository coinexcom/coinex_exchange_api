#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
Created by bu on 2018-01-17
"""
from __future__ import unicode_literals
import time
import hashlib
import json as complex_json
import requests


class RequestClient(object):
    __headers = {
        'Content-Type': 'application/json; charset=utf-8',
        'Accept': 'application/json',
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36'
    }

    def __init__(self, headers={}):
        self.access_id = 'xxx'      # replace
        self.secret_key = 'xxx'     # replace
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
        if method == 'GET':
            self.set_authorization(params)
            print url
            print params
            result = requests.request('GET', url, params=params, headers=self.headers)
        else:
            if data:
                json.update(complex_json.loads(data))
            self.set_authorization(json)
            print url
            print json
            result = requests.request(method, url, json=json, headers=self.headers)
        return result


def get_account():
    request_client = RequestClient()
    response = request_client.request('GET', 'https://api.coinex.com/v1/balance/')
    print response.content


def put_limit():
    request_client = RequestClient()

    data = {
            "amount": "0.4199",
            "price": "1",
            "type": "sell",
            "market": "CDYBCH"
        }

    response = request_client.request(
            'POST',
            'https://api.coinex.com/v1/order/limit',
            json=data,
    )
    print response.content

if __name__ == '__main__':
    print get_account()
