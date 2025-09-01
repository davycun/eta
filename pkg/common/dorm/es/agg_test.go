package es_test

/*

GET /niuta_dev_339240270212108288_t_enterprise/_search
{
  "size": 0,
  "query": {
    "bool": {
      "filter": [
        {
          "exists": {
            "field": "office_addr_ids"
          }
        },
        {
          "exists": {
            "field": "enroll_addr_ids"
          }
        },
        {
          "nested": {
            "query": {
              "bool": {
                "filter": [
                  {
                    "term": {
                      "address_list.addr_type.keyword": "subdistrict"
                    }
                  },
                  {
                    "term": {
                      "address_list.parent_id.keyword": "74415626332016640"
                    }
                  }
                ]
              }
            },
            "path": "address_list"
          }
        }
      ]
    }
  },
  "aggs":{
    "nested_buckets":{
      "nested": {
        "path": "address_list"
      },
      "aggs": {
        "filter_buckets": {
          "filter": {
            "bool": {
              "filter":[
                {
                  "term": {
                    "address_list.addr_type.keyword":"subdistrict"
                  }
                }
              ]
            }
          },
          "aggs": {
            "composite_buckets":{
              "composite": {
                "size": 20,
                "after": {
                  "address_list_id":"74417133714542603",
                  "address_addr_type":"subdistrict"
                },
                "sources": [
                  {
                    "address_list_id": {
                      "terms": {
                        "field": "address_list.id.keyword",
                        "order": "asc"
                      }
                    }
                  },
                  {
                    "address_addr_type": {
                      "terms": {
                        "field": "address_list.addr_type.keyword",
                        "order": "asc"
                      }
                    }
                  }
                ]
              },
              "aggs":{
                "my_count":{
                  "cardinality": {
                    "field": "address_list.name.keyword"
                  }
                },
                "my_max":{
                  "max": {
                    "field": "address_list.level"
                  }
                },
                "my_sum":{
                  "sum": {
                    "field": "address_list.level"
                  }
                },
                "my_having":{
                  "bucket_selector": {
                    "buckets_path": {
                      "myBuckets":"_count"
                    },
                    "script": "params.myBuckets > 70"
                  }
                },
                "my_having2":{
                  "bucket_selector": {
                    "buckets_path": {
                      "myBuckets":"my_max"
                    },
                    "script": "params.myBuckets > 2"
                  }
                }
              }
            }
          }
        }
      }
    },
    "group_count": {
      "nested": {
        "path": "address_list"
      },
      "aggs": {
        "nested_filtered": {
          "filter": {
            "bool": {
              "filter":[
                {
                  "term": {
                    "address_list.addr_type.keyword": "subdistrict"
                  }
                },
                {
                  "term": {
                    "address_list.parent_id.keyword": "74415626332016640"
                  }
                }
              ]
            }
          },
          "aggs": {
            "my_count": {
              "cardinality": {
                "field": "address_list.id.keyword",
                "precision_threshold": 10000
              }
            }
          }
        }
      }
    }
  }
}

*/

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
