from fastapi import Response

def build_rds_describe_db_instances_response() -> Response:
    with open("./data/rds/DescribeDbInstances.xml") as f:
        content = f.read()
    return Response(content=content, media_type="text/xml")
