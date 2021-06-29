# home24-web scraper

Home24 Go Challenge

## Libraries used

Name  | Description
------------- | -------------
Chi Router  | For routing and middleware
Zap  | For logging
x/net/Html  | For html parsing


## How to run

### With Docker

1. `docker build . -t home24-test`

2. `docker run -p 3000:3000 home24-test`


### without Docker

1. download golang 1.16
2. checkout https://github.com/Aegon95/home24-webscraper
3. run `go mod download`
4. go to cmd/api folder
5. run `go run .`
6. it will start the web server at port 3000


## Screenshots

![image](https://user-images.githubusercontent.com/42437992/123549209-97a92900-d785-11eb-8767-86bc606eca93.png)

![image](https://user-images.githubusercontent.com/42437992/123549232-b7d8e800-d785-11eb-87eb-1ef095482a1b.png)

![image](https://user-images.githubusercontent.com/42437992/123549236-c0312300-d785-11eb-8a72-4d7f759a0c80.png)




