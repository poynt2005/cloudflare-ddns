# cloudflare-ddns

## 自動使用 cloudflare 提供的 api 將本機端 ip 更新 domain name
## 如果你有個人 domain name 可以在主機上設置一個 cronjob 來跑這個 image，當作是 DDNS 來用
## 使用方法
首先要先設置一個 cloudflare api key, permission 選 ZONE/DNS/Edit.  
`
docker run --rm -e CLOUDFLARE_DNS_API_KEY=<你的 api key> -e CLOUDFLARE_DNS_ZONE_NAME=<你的域名> -e CLOUDFLARE_DNS_SUBDOMAIN_NAME=<你想叫的子網域名稱> -it poynt2005/cloudflare-ddns
`
