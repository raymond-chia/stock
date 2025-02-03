from backtesting import Strategy
from backtesting.lib import crossover
from ta import *


# deepseek 1/27 衝擊 AI 股市


# 3293
# MACD 死亡交叉, 且 CDP 沒收縮 => 跌
# 2025 MACD 死亡交叉, 但 CDP 收縮 => ??


class SMACross(Strategy):
    def init(self):
        fast_n = 5
        slow_n = 20
        self.fast_line = self.I(SMA, self.data.Close, fast_n)
        self.slow_line = self.I(SMA, self.data.Close, slow_n)
        self.hold = 0

    def next(self):
        if crossover(self.fast_line, self.slow_line):
            # print(f"{self.data.index[-1]} Buy")
            self.buy()
        elif crossover(self.slow_line, self.fast_line):
            # print(f"{self.data.index[-1]} Sell")
            self.sell()
