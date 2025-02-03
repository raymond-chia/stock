import yfinance as yf
import ta
import scipy.signal as sci
from pick.list import *
from pick.list_generator import *


# 日
KDJDailyFilter = 30.0
BBIDailyFilter = 2.0
# 周
KDJWeeklyFilter = 45.0
BBIWeeklyFilter = 2.0
# 月 算出來跟 yahoo 差距過大

# 成交量不低
# rsi 替代 kdj

VolumeThresholdWeekly = 2
VolumeThresholdMonthly = 3


class Indicator:
    Volume = 1
    KDJ = 2
    M3_BBI = 4
    DI_Plus = 8
    DI_Minus = 16


# 每股盈餘 (奇摩基本) x 營收年增% x 本益比 = 預期股價 ??
# ROE 15%+ (財報狗)
# 營業利益率 10%+ (奇摩基本)
# RSI 當 RSI 高於 70, 表示股價處於超買區, 可能會回調; 當 RSI 低於 30, 表示股價處於超賣區, 可能會反彈
# DMI +DI尖賣; -DI尖買
# DMI +DI 上漲程度; -DI 下降程度; AD https://tw.stock.yahoo.com/news/%E6%8A%80%E8%A1%93%E5%88%86%E6%9E%90-dmi%E6%8C%87%E6%A8%99-dmi%E5%8B%95%E5%90%91%E6%8C%87%E6%A8%99-%E5%A4%9A%E7%A9%BA%E6%96%B9%E5%90%91-%E8%B6%A8%E5%8B%A2%E5%8B%95%E8%83%BD-130256849.html
# MACD 長期趨勢. 周線 負正黃金交叉; 正負死亡交叉
# 多空指標乖離 周線 不要太高, 可以很負
# CDP 短線趨勢. 收縮賣
# Yahoo 股市: 本益比, 股利, 財務, 基本
# Goodinfo 財務評分表
# 股價預估值 ?? https://www.findbillion.com/twstock/1231
# 今年營收成長大於 0% ?


def pick():
    result = {"d": {}, "w": {}, "m": {}}
    for stock_id in IDs:
        stock = check_one(stock_id)
        for interval in ["d", "w", "m"]:
            status = stock[interval]
            if status > 0:
                result[interval][stock_id] = status
    daily = set()
    weekly = set()
    status = result["w"]
    for stock_id in status:
        v = status[stock_id]
        if v & Indicator.KDJ and v & Indicator.M3_BBI:
            weekly.add(stock_id)
    status = result["d"]
    for stock_id in status:
        v = status[stock_id]
        if (v & Indicator.KDJ and v & Indicator.M3_BBI) and stock_id not in weekly:
            daily.add(stock_id)
    print("daily")
    for stock_id in daily:
        print(f"{stock_id} {IDs[stock_id]}")
    print("weekly")
    for stock_id in weekly:
        print(f"{stock_id} {IDs[stock_id]}")
    print("DMI -DI peak")
    status = result["w"]
    for stock_id in status:
        v = status[stock_id]
        if v & Indicator.DI_Minus:
            print(f"{stock_id} {IDs[stock_id]}")
    print("week volume")
    status = result["w"]
    for stock_id in status:
        v = status[stock_id]
        if v & Indicator.Volume:
            print(f"{stock_id} {IDs[stock_id]}")
    print("month volume")
    status = result["m"]
    for stock_id in status:
        v = status[stock_id]
        if v & Indicator.Volume:
            print(f"{stock_id} {IDs[stock_id]}")


def check_one(id):
    d = yf.Ticker(id)
    r = {}
    r["d"] = _daily(d)
    r["w"] = _weekly(d)
    r["m"] = _monthly(d)
    return r


def _daily(d):
    df = d.history(period="max")
    result = 0

    k, _, _ = ta.KDJ(df.High, df.Low, df.Close)
    if k.iloc[-1] <= KDJDailyFilter:
        result += Indicator.KDJ
    bbi = ta.BullBearIndex(df.Close)
    if ta.SMA(df.Close, timeperiod=3).iloc[-1] - bbi.iloc[-1] <= BBIDailyFilter:
        result += Indicator.M3_BBI
    return result


def _weekly(d):
    df = d.history(period="max", interval="1wk")
    result = 0

    k, _, _ = ta.KDJ(df.High, df.Low, df.Close)
    if k.iloc[-1] <= KDJWeeklyFilter:
        result += Indicator.KDJ
    bbi = ta.BullBearIndex(df.Close)
    if ta.SMA(df.Close, timeperiod=3).iloc[-1] - bbi.iloc[-1] <= BBIWeeklyFilter:
        result += Indicator.M3_BBI

    weeks = 3

    di_plus = ta.PLUS_DI(df.High, df.Low, df.Close, timeperiod=14)
    di_minus = ta.MINUS_DI(df.High, df.Low, df.Close, timeperiod=14)
    di_plus = di_plus.iloc[-weeks:]
    di_minus = di_minus.iloc[-weeks:]
    peak_plus, _ = sci.find_peaks(di_plus)
    peak_minus, _ = sci.find_peaks(di_minus)
    if len(peak_plus) > 0 and len(peak_minus) > 0:
        # TODO
        # if peak_plus[-1] > peak_minus[-1]:
        #     result += Indicator.DI_Plus
        if peak_plus[-1] < peak_minus[-1]:
            result += Indicator.DI_Minus
    # TODO
    # elif len(peak_plus) > 0:
    #     result += Indicator.DI_Plus
    elif len(peak_minus) > 0:
        result += Indicator.DI_Minus

    if df.Volume.iloc[-2] * VolumeThresholdWeekly < df.Volume.iloc[-1]:
        result += Indicator.Volume
    # TODO
    # print(df.Volume.iloc[-2] * VolumeThresholdWeekly, df.Volume.iloc[-1])
    return result


def _monthly(d):
    df = d.history(period="max", interval="1mo")
    result = 0

    if df.Volume.iloc[-2] * VolumeThresholdMonthly < df.Volume.iloc[-1]:
        result += Indicator.Volume
    # TODO
    # print(df.Volume.iloc[-2] * VolumeThresholdMonthly, df.Volume.iloc[-1])
    return result
