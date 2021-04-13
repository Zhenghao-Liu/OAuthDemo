* redis存储
1. `unordered_map<account,string> appID`
2. `unordered_map<appID,string> account`
3. `[code:${string}]=appID+"_"+account+"_"+scope`
4. `[token:${string}]=appID+"_"+account+"_"+scope`
5. `[refresh:${string}]=appID+"_"+account+"_"+scope`

* 场景
1. `user_info`修改，通过`unordered_map<appID,string> account`把所有`[code:${string}]`、`[token:${string}]`、`[refresh:${string}]`删除
2. `oauth_info`修改，通过`unordered_map<account,string> appID`把所有`[code:${string}]`、`[token:${string}]`、`[refresh:${string}]`删除
3. 请求授权码，验证`account`与`password`与`appID`，验证其他参数，如果原先`appID[account]、account[appID]`有值`s_old`，那么应该取旧值并删除`[code:${s_old}]`、`[token:${s_old}]`、`[refresh:${s_old}]`，生成`s1`，`appID[account]=s1`、`account[appID]=s1`、`[code:${s1}]=appID+"_"+account+"_"+scope`，返回是s1即code
4. 请求令牌，验证`appID`与`appSecret`与`code`，验证其他参数，从`[code:${s1}]=appID+"_"+account+"_"+scope`拿到`appID`、`account`，删除`[code:${s1}]=appID+"_"+account+"_"+scope`，生成`s2`与`s3`，`appID[account]=account[appID]=s2+s3`、`[token:${s2}]=[refresh:${s3}]=appID+"_"+account+"_"+scope`，返回`s2`即`token`，`s3`即`refresh_token`
5. 刷新令牌，验证`appID`与`refresh_token`，从`[refresh:${refresh_token}]`中获取`appID`、`account`，新生成`ss2`与`ss3`，`appID[account]=account[appID]=ss2+ss3`、`[token:${ss2}]=[refresh:${ss3}]=appID+"_"+account+"_"+scope`，返回`ss2`即`token`，`ss3`即`refresh_token`
6. 申请资源，验证`appID`与`token`，从`[token:${string}]=appID+"_"+account+"_"+scope`中拿`scope`，之后去对应的`account`中拿对应的`scope`
7. s1哈希后成为code的key，s2+s3逆哈希后，s2哈希后成为token的key，s3哈希后成为refresh的key，存的都是哈希过后的值，返回给客户端的都是没有加密的值