import yfinance as yf
import ta
from backtesting import Backtest
import strategy
import pick

# https://ithelp.ithome.com.tw/articles/10314112
# 成交量大 => 不容易被主力控制
# 50 億以下 (小型股), 50 億 ~ 100 億 (中型股), 100 億以上 (大型股)
# - 中小型股
#   - 大量: 10日均量的3倍
#   - 爆量: 10日均量的5倍
#   - 天量: 10日均量的10倍
# - 大型股
#   - 大量: 比10日均量多出30%
#   - 爆量: 比10日均量多出50%
#   - 天量: 比10日均量多出100%


# 台灣景氣對策信號 亮紅燈後, 平均約 9.5 個月指數達到高點
#   無視極值可能 11 個月
#   2025/3 ~ 6 達到高點 ??
#   選運輸 ??
#   https://www.youtube.com/watch?v=f6y-OM3l1OI
# 航運 可以參考 bdi ??


# # 鈊象
# stock_id = "3293.TWO"

# d = yf.Ticker(stock_id)
# # df = d.history(period="max")
# df = d.history(start="2020-01-01")
# # print("市值:", d.info["marketCap"])
# # print(df)

# k, d, j = ta.KDJ(df.High, df.Low, df.Close)
# # print("kdj")
# # print(k)
# # print(d)
# # print(j)

# # macd, macd_signal, macd_hist = ta.MACD(
# #     df.Close, fastperiod=12, slowperiod=26, signalperiod=9
# # )
# # macd, macd_signal, macd_hist = ta.MACDFIX(df.Close, signalperiod=9)
# # TODO 確認不同是否造成影響
# # print("macd")
# # print(macd)
# # print(macd_signal)
# # print(macd_hist)

# rsi5 = ta.RSI(df.Close, timeperiod=5)
# # print("rsi5")
# # print(rsi5)

# bias10 = ta.BIAS(df.Close, timeperiod=10)
# # print("乖離率")
# # print(bias10)

# will = ta.WILLR(df.High, df.Low, df.Close, timeperiod=9)
# will = 0 - will
# # print("威廉指標")
# # print(will)

# bbi = ta.BullBearIndex(df.Close)
# # https://tayu.tripod.com/fashion2-19.htm
# # print("多空指標乖離")
# # print(ta.SMA(df.Close, timeperiod=3) - bbi)

# # https://www.wantgoo.com/blog/98845/post/72
# # https://blog.csdn.net/wjl__ai__/article/details/118075594
# CDP, AH, NH, NL, AL = ta.CDP(df.High, df.Low, df.Close)
# # print("CDP")
# # print(CDP)
# # print(AH)
# # print(NH)
# # print(NL)
# # print(AL)

# adx = ta.ADX(df.High, df.Low, df.Close, timeperiod=14)
# adxr = ta.ADXR(df.High, df.Low, df.Close, timeperiod=14)
# di_plus = ta.PLUS_DI(df.High, df.Low, df.Close, timeperiod=14)
# di_minus = ta.MINUS_DI(df.High, df.Low, df.Close, timeperiod=14)
# # print("DMI")
# # print(adx)
# # print(adxr)
# # print(di_plus)
# # print(di_minus)


pick.pick()

# test = Backtest(
#     df,
#     strategy.SMACross,
#     cash=10000,
#     commission=0.004,
#     trade_on_close=True,
#     exclusive_orders=True,
# )
# result = test.run()
# test.plot(filename=f"./py/backtest/{stock_id}.html")
# print(result)
