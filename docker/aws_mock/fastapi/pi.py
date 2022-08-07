import json
from typing import List

from fastapi import Request
from datetime import datetime, timedelta, time

DB_ID = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
EVERY_PERIOD_TOKENIZED_ID = "TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT"
EVERY_PERIOD_VALUE = 0.01

async def build_pi_get_resource_metrics_response(request: Request) -> {}:
    j = await request.json()
    start = get_next_5min_period(datetime.fromtimestamp(j["StartTime"]))
    end = get_next_5min_period(datetime.fromtimestamp(j["EndTime"]))
    identifier = j["Identifier"]

    group = j["MetricQueries"][0]["GroupBy"]["Group"]
    if group == "db.sql":
        metric_list = __build_sql_metric_list(start, end)
    elif group == "db.sql_tokenized":
        metric_list = __build_tokenized_sql_metric_list(start, end)
    else:
        metric_list = []

    data = {
        "AlignedStartTime": start.timestamp(),
        "AlignedEndTime": end.timestamp(),
        "Identifier": identifier,
        "MetricList": metric_list
    }
    return data

def get_next_5min_period(date: datetime) -> datetime:
    new_date = date.replace(second=0, microsecond=0)
    if date == new_date:
        return date
    diff_minute = 5-new_date.minute%5
    return new_date + timedelta(minutes=diff_minute)

def __build_sql_metric_list(start: datetime, end: datetime):
    step_count = int((end-start) / timedelta(minutes=5))
    target_time_range = [start + timedelta(minutes=5*x) for x in range(step_count)]
    base_metric = {
        "Key": {"Metric": "db.load.avg"},
        "DataPoints": [
            {
                "Timestamp": d.timestamp(),
                "Value": 0.0
            } for d in target_time_range
        ]
    }

    every_minute_metric = {
        "Key": {
            "Metric": "db.load.avg",
            "Dimensions": {
                "db.sql.DB_ID": DB_ID,
                "db.sql.id": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
                "db.sql.statement": "select * from users where email = 'test1@example.com'",
                "db.sql.tokenized_id": EVERY_PERIOD_TOKENIZED_ID
            }
        },
        "DataPoints": [
            {
                "Timestamp": d.timestamp(),
                "Value": EVERY_PERIOD_VALUE
            } for d in target_time_range
        ]
    }
    metric_list = [base_metric, every_minute_metric]

    with open("./data/pi/GetResourceMetrics_data.json") as f:
        sqls = json.load(f)

    for sql in sqls:
        sql_time = time.fromisoformat(sql["Time"])
        metric_list.append({
            "Key": {
                "Metric": "db.load.avg",
                "Dimensions": {
                    "db.sql.DB_ID": DB_ID,
                    "db.sql.id": sql["db.sql.id"],
                    "db.sql.statement": sql["db.sql.statement"],
                    "db.sql.tokenized_id": sql["db.sql.tokenized_id"],
                }
            },
            "DataPoints": [
                {
                    "Timestamp": d.timestamp(),
                    "Value": sql["Value"] if d.time() <= sql_time < (d+timedelta(minutes=5)).time() else 0
                } for d in target_time_range
            ]
        })

    return metric_list

def __build_tokenized_sql_metric_list(start: datetime, end: datetime):
    step_count = int((end-start) / timedelta(minutes=5))
    target_time_range = [start + timedelta(minutes=5*x) for x in range(step_count)]
    base_metric = {
        "Key": {"Metric": "db.load.avg"},
        "DataPoints": [
            {
                "Timestamp": d.timestamp(),
                "Value": 0.0
            } for d in target_time_range
        ]
    }

    every_period_metric = {
        "Key": {
            "Metric": "db.load.avg",
            "Dimensions": {
                "db.sql_tokenized.DB_ID": DB_ID,
                "db.sql_tokenized.statement": "select * from `users` where `email` = ?",
                "db.sql_tokenized.id": EVERY_PERIOD_TOKENIZED_ID
            }
        },
        "DataPoints": [
            {
                "Timestamp": d.timestamp(),
                "Value": EVERY_PERIOD_VALUE
            } for d in target_time_range
        ]
    }
    metric_list = [base_metric, every_period_metric]

    with open("./data/pi/GetResourceMetrics_tokenized_data.json") as f:
        sqls = json.load(f)

    for sql in sqls:
        sql_timestamp = time.fromisoformat(sql["Time"])
        metric_list.append({
            "Key": {
                "Metric": "db.load.avg",
                "Dimensions": {
                    "db.sql_tokenized.DB_ID": DB_ID,
                    "db.sql_tokenized.statement": sql["db.sql_tokenized.statement"],
                    "db.sql_tokenized.id": sql["db.sql_tokenized.id"],
                }
            },
            "DataPoints": [
                {
                    "Timestamp": d.timestamp(),
                    "Value": sql["Value"] if d.time() <= sql_timestamp < (d+timedelta(minutes=5)).time() else 0
                } for d in target_time_range
            ]
        })

    return metric_list
