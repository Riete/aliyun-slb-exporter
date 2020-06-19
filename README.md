### docker build
``` docker build . -t <image>:<tag> ```

### or pull 
``` docker pull riet/aliyun-slb-exporter ```

### run
get all slb metrics
```
docker run \ 
  -d \ 
  --name aliyun-slb-exporter \
  -e ACCESS_KEY_ID=<aliyun ak> \
  -e ACCESS_KEY_SECRET=<aliyun ak sk> \
  -e REGION_ID=<region id> \
  -p 10002:10002 \
  riet/aliyun-slb-exporter 
```

get specified slb metrics
```
docker run \ 
  -d \ 
  --name aliyun-slb-exporter \
  -e ACCESS_KEY_ID=<aliyun ak> \
  -e ACCESS_KEY_SECRET=<aliyun ak sk> \
  -e REGION_ID=<region id> \
  -e INSTANCE_ID=id1,id2,id3 \
  -p 10002:10002 \
  riet/aliyun-slb-exporter 
```

visit http://localhost:10002/metrics