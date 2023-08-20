# curl -X POST 127.0.0.1:8989/user/json/render
{
    "code": 10000,
    "msg": "Method Not Allowed",
    "data": "/user/json/render not allowed request POST"
}

# curl -X POST 127.0.0.1:8989/json/render
{"code":10000,"msg":"Not Found","data":"404 not found "}

# curl 127.0.0.1:8989/user/json/render
多出来一个换行，json库的问题，后面看是不是不用标准库
{"code":200,"msg":"OK","data":{"Name":"\u003ch1\u003ewendell\u003c/h1\u003e","Age":25}
}

curl http://127.0.0.1:8989/user/handlerErr
{"code":500,"msg":"数据异常","data":null}

curl 127.0.0.1:8989/user/xml/render
{"code":200,"msg":"OK","data":"<User><name>wendell</name><age>25</age></User>"}

curl 127.0.0.1:8989/user/string/render
{"code":200,"msg":"OK","data":"name: wendell age: 26"}

curl 127.0.0.1:8989/user/success/render
{"code":200,"msg":"OK","data":"成功"}
