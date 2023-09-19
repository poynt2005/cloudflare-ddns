# cloudflare-ddns

### 更新說明

原專案的 [commit](https://github.com/poynt2005/cloudflare-ddns/commit/1fdf7f14a9c2d867f210a32bbcd859168aa29d70)  
在 2023-09-19 經過一些調整、代碼重構之後的專案  
調整重點:

1. 完全去除 python ，採用 cloudflare 原生 http api
2. docker image 鏡像容量: 64.1 兆字節 -> 7.74 兆字節，節省容量達 88%

### 說明

自動使用 cloudflare 提供的 api 將本機端 ip 更新 domain name  
如果你有個人 domain name 可以在主機上設置一個 cronjob 來跑這個 image，當作是 DDNS 來用

## 使用方法

首先要去 cloudflare 的 domain 管理面板 設置一個 cloudflare api key, permission 選 ZONE/DNS/Edit.  
再來請參考 runscript.txt，其中

1. Your_Cloudflare_Api_Key: 你獲取的 api key
2. Your_Cloudflare_Main_Domain_Name: 你的主要 domain，例如 example.com
3. Your_Cloudflare_Subdomain_Name: 你要映射的子網域名稱，例如你要映射 abcd.example.com 這個值就寫 abcd

以上三個值缺一不可，缺了就會報錯  
可搭配 cronjob 進行定期的域名映射更新
