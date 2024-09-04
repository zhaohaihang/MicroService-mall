#! /bin/bash

list="user_api goods_api order_api userop_api"

# 使用for循环遍历列表
for item in $list; do
   nohup ./apis/$item/$item --nacosConfig $1 &
done