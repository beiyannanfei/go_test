-- example HTTP POST script which demonstrates setting the
-- HTTP method, body, and adding a header

wrk.method = "POST"
wrk.body   = [[{
        "Token":"f3c302a3f4dccf628e883557acbbf433",
        "GlobalId": 949956626188603392
        }]]
wrk.headers["Content-Type"] = "application/json"