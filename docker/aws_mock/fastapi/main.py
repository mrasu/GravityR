from fastapi import FastAPI, Header, Response, Request
from starlette.responses import JSONResponse

from pi import build_pi_get_resource_metrics_response
from rds import build_rds_describe_db_instances_response

app = FastAPI()

@app.get("/")
def read_root():
    return {"pong": "hello"}

@app.post("/")
async def read_aws(request: Request, content_type: str = Header(default="")):
    action = await get_action(request)
    if action == "PerformanceInsightsv20180227.GetResourceMetrics":
        return await build_pi_get_resource_metrics_response(request)
    elif action == "DescribeDBInstances":
        return build_rds_describe_db_instances_response()
    else:
        return build_unknown_action_response(content_type, action)

async def get_action(request: Request) -> str:
    x_amz_target = request.headers.get("X-Amz-Target", "")
    form = await request.form()
    return form.get("Action", x_amz_target)

def build_unknown_action_response(content_type: str, action: str) -> Response:
    if content_type == "application/x-www-form-urlencoded":
        content = f"""
<ErrorResponse xmlns="http://rds.amazonaws.com/doc/2014-10-31/">
  <Error>
    <Type>Sender</Type>
    <Code>UnknownAction</Code>
    <Message>Unknown Action: {action}</Message>
  </Error>
  <RequestId>799fa280-3701-465f-8b26-49edcf814748</RequestId>
</ErrorResponse>
"""
        return Response(status_code=400, content=content, media_type="text/xml")
    else:
        detail = {"__type": "UnknownTargetException","Message": f"Unknown X-Amz-Target: {action}"}
        return JSONResponse(status_code=400, content=detail)
