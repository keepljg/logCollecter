# logCollecter
### 架构
 埋点数据/日志上报  ->  kafka  -> logCollecter(多个) -> (ingest pipeline)elasticsearch -> 数据分析（kibana/grafana）
 
 分布式是利用etcd完成
 logCollecter 可随时动态新增和删除
 
 
 
                      
                      
                                  
