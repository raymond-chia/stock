from pick.list import *
import yfinance as yf
import requests
import time
import json


# 過濾市值太低的
# 要拿掉 ETF
def show_large_market_cap():
    for stock_id in IDs:
        if not is_large_market_cap(stock_id):
            continue
        name = get_chinese_name(stock_id)
        print(f'"{stock_id}": "{name}",')


# 要拿掉 ETF
def is_large_market_cap(stock_id):
    d = yf.Ticker(stock_id)
    d.info["marketCap"]
    m = d.info["marketCap"]
    return m > 10000000000


def get_chinese_name(stock_id):
    stock_id = stock_id.removesuffix(".TWO").removesuffix(".TW")
    query = {
        "v": "1",
        "type": "ta",
        "mkt": "10",
        "sym": stock_id,
        "perd": "d",
        "_": time.time() * 1000,
        "callback": "",
    }
    result = requests.get(f"https://tw.quote.finance.yahoo.net/quote/q", params=query)
    result = result.text.removeprefix("(").removesuffix(");")
    result = json.loads(result)
    return result["mem"]["name"]
