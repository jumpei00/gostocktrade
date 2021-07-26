![file-1587959572966-golang_eyecatch](https://user-images.githubusercontent.com/61613314/124248228-98630600-db5d-11eb-8088-97353ee7c78b.png)

![stock](https://user-images.githubusercontent.com/61613314/126995400-581700e0-1e59-41dd-a802-6a4bb35c5ce5.gif)

# GO Stock Trade Tool
This application enables you to judge whether you should buy stock, sell stock, or not.

# Overveiw
What this application implements is as follows.
- candle stick of stock data(daily data, using yahoo api)
- indicators(using HighChrats)
- backtest of EMA, BollingerBand, MACD, RSI, WilliamR
- display trade timing of past
- display whether today is BUY, or SELL, or not

# Usage
## generate
```
$ go mod tidy
$ go run main.go
```
or
```
$ docker-compose up --build
```
and access 127.0.0.1:8080
## test
```
$ go mod tidy
$ go -v -cover ./...
```

# Requirements
- GO 1.16.3
- go-quote latest
- go-talib latest
- logrus 1.8.1
- goconvey 1.6.4
- testify 1.7.0
- ini 1.62.0
- sqlite 1.1.4
- gorm 1.21.10

# Author
Jumpei Motohashi

# Licence
no licence, but due to use highchats library, NOT use as commercial.
