package es_test

//一个聚合的示例
/*
##楼宇聚合带条件，带分页，带having，带额外字段
POST /delta_dev_backend_t_wide_building/_search
{
  "size": 0,
  "track_total_hits": false,
  "query": {
    "bool": {
      "filter": [
        {
          "match_phrase":{
            "addr_ids":"46"
          }
        }
      ]
    }
  },
  "aggs": {
    "my_addr_ids_count":{
      "cardinality": {
        "field": "addr_ids"
      }
    },
    "my_addr_ids": {
      "terms": {
        "field":"addr_ids",
        "order": {
          "_count": "desc"
        },
        "size": 10
      },
      "aggs": {
        "max_usage":{
          "max": {
            "field": "usage"
          }
        },
        "min_usage":{
          "min": {
            "field": "usage"
          }
        },
        "max_usage_having":{
          "bucket_selector": {
            "buckets_path": {
              "usage_max":"max_usage"
            },
            "script": "params.usage_max <=4"
          }
        },
        "min_usage_having": {
          "bucket_selector": {
            "buckets_path": {
              "usage_min":"min_usage"
            },
            "script": "params.usage_min>=1"
          }
        },
        "count_having": {
          "bucket_selector": {
            "buckets_path": {
              "count_agg":"_count"
            },
            "script": "params.count_agg>=275"
          }
        },
        "agg_page":{
          "bucket_sort": {
            "from": 0,
            "size": 5
          }
        }
      }
    }
  }
}
*/
