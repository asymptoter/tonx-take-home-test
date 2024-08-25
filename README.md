# Tonx Take Home Assignment
### Requirements
1.  每天定點 23 點發送優惠券，需要用戶提前 1-5 分鐘前先預約，優惠券的數量為預約用戶數量的 20%。

    **- 用戶每天只有 22:55:00 到 22:58:59 (GMT+8) 可以預約（打預約搶購優惠券的API）**
    **- 假設預約用戶有 300 人，則優惠券數量有 15 份**
2. 搶購優惠券時，每個用戶只能搶一次，並且只有 1 分鐘的搶購時間，如何確保用戶搶到的概率接近 20%。

    **- 假設所有預約用戶都會搶購**
    **- 在用戶預約的時候就決定該用戶是否獲得優惠券**
    **- 用戶打 API 搶購只會讀取 Database 的一個 row**
3. 需要考慮 300 人搶購和 30000 人搶購的場景。
    
    **- 300 人的情況不管是預約或是搶購應該可以不用做任何 cache，搶購 API 直接 access database**
    **- 30000 人的情況預約的 API 可以用非同步的方式進行，用 message queue 的方式將寫入 database 的操作減少，原本要寫入 30000次的操作可以一次寫入多筆（可以每收集到 100 筆寫入一次，或是每 5 秒寫入一次），搶購的 API 則可以先把預約時已經得到的結果從 database 讀到 redis 或/和 local memory 做 cache**
         
4. 設計表結構和索引，並編寫主要程式碼。
    - Campaigns
    
    |Column Name|Type|Index|
    |-----------|----|-----|
    |id|unsigned int|Primary Key|
    |created_at|timestamp|Index|
        campaign_id 在這張表必須是 unique，否則重複的 campaign_id 會導致查詢 Reservations 會出錯
        campaign_id 用日期最簡單，但如果未來需求變更成每天會發送多次優惠券的話就很難改動
    
    - Coupon_Reservations
    
    |Column Name|Type|Index|
    |-----------|----|-----|
    |user_id|uuid|Primary Key=campaign_id+user_id|
    |campaign_id|unsigend int|Primary Key=campaign_id+user_id,Foreign Key Reference Campaigns.id|
    |coupon_code|text||
        user_id 假設為系統指定的 uuid 
