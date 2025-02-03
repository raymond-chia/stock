import numpy as np
import pandas as pd
from talib import *


# RSV = (Ct - Ln) ÷ (Hn - Ln) × 100
# Ct 為第 n 日收盤價
# Ln 為 n 日內最低價
# Hn 為 n 日內最高價
def RSV(H, L, C):
    l = L.rolling(9).min()
    h = H.rolling(9).max()
    return 100 * ((C - l) / (h - l)).values


# https://www.futunn.com/hk/learn/detail-what-is-kdj-64858-220831019
# n 請見 RSV
# 其他實作方式: https://github.com/TA-Lib/ta-lib-python/issues/99#issuecomment-653953212
def KDJ(H, L, C):
    # TODO
    # 有些股票的 RSV 會變成 nan
    # 先縮小計算範圍, 降低遇到的機率
    H = H[-100:]
    L = L[-100:]
    C = C[-100:]
    rsv = RSV(H, L, C)

    k = 50
    d = 50
    k_out = []
    d_out = []
    for j in range(len(rsv)):
        if rsv[j] == rsv[j]:  # check for nan
            k = 1 / 3 * rsv[j] + 2 / 3 * k
            k_out.append(k)
            d = 1 / 3 * k + 2 / 3 * d
            d_out.append(d)
        else:
            k_out.append(np.nan)
            d_out.append(np.nan)

    j_out = (3 * np.array(k_out)) - (2 * np.array(d_out))
    return (
        pd.Series(k_out, H.index),
        pd.Series(d_out, H.index),
        pd.Series(j_out, H.index),
    )


# 乖離
def BIAS(C, timeperiod):
    ma = SMA(C, timeperiod=timeperiod)
    return (C - ma) / ma * 100


# https://peggy0501.pixnet.net/blog/post/16240317
def BullBearIndex(C):
    # return (
    #     SMA(C, timeperiod=3)
    #     + SMA(C, timeperiod=6)
    #     + SMA(C, timeperiod=12)
    #     + SMA(C, timeperiod=24)
    # ) / 4
    return (
        SMA(C, timeperiod=3)
        + SMA(C, timeperiod=6)
        + SMA(C, timeperiod=9)
        + SMA(C, timeperiod=12)
    ) / 4


def CDP(H, L, C):
    CDP = (H + L + C * 2) / 4
    AH = CDP + (H - L)
    NH = 2 * CDP - L
    NL = 2 * CDP - H
    AL = CDP - (H - L)
    return CDP, AH, NH, NL, AL
