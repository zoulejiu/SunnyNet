import * as monaco from "monaco-editor";
import axios from "axios";
import {builtCmdWords} from "./builtCmdWords.js";

const SunnyNetConnName = "Conn"
//const requestEventFunc = './src/assets/EventFunc.json'
export const localHost = (document.location.toString().indexOf("https://") > -1 ? "https" : "http") + `://${window.location.host}${window.location.pathname}`
const requestEventFunc = localHost + `/getEventFunc`
const requestEventCodeType = localHost + "/getType?"

export function HoverEvent(funcName, words, target) {
    if (funcName === null) {
        if (words.length === 2 && words[0] === "func") {
            switch (target) {
                case "Event_HTTP":
                    return {
                        contents: [{value: "入口函数:脚本代码 HTTP 请求脚本处理"}]
                    };
                case "Event_WebSocket":
                    return {
                        contents: [{value: "入口函数:脚本代码 WebSocket 请求脚本处理"}]
                    };
                case "Event_TCP":
                    return {
                        contents: [{value: "入口函数:脚本代码 TCP 请求脚本处理"}]
                    };
                case "Event_UDP":
                    return {
                        contents: [{value: "入口函数:脚本代码 UDP 请求脚本处理"}]
                    };
            }
        }
        if (words.length >= 3 && words[0] === "func" && words[2] === SunnyNetConnName) {
            const ret = HoverEventFuncStruct(words[1])
            if (ret !== null) {
                return ret;
            }
        }

    } else if (target === SunnyNetConnName) {
        return HoverEventFuncStruct(funcName);
    } else if (words.length >= 3 && words[0] === "func" && words[2] === SunnyNetConnName) {
        const ret = HoverEventFuncStruct(words[1])
        if (ret !== null) {
            return ret;
        }
    }
    if (funcName !== null) {
        if (words.length >= 2 && words[words.length - 2] === SunnyNetConnName) {
            return HoverEventFuncMembers(funcName, target)

        }
    }
    for (let index = 0; index < builtCmdWords.length; index++) {
        const obj = builtCmdWords[index]
        for (let n = 0; n < obj.name.length; n++) {
            if (obj.name[n] === target) {
                let contents = ""
                for (let i = 0; i < obj.contents.length; i++) {
                    contents += obj.contents[i].value + "\n\n"
                }
                return {
                    contents: [{value: contentsReplaceCode(contents)}]
                };
            }
        }
    }
    console.log("HoverEvent", funcName, words, target)
}

function HoverEventFuncMembers(func, target) {
    let obj = null;
    if (func === "Event_HTTP") {
        obj = exportHTTPInterface[target]
    } else if (func === "Event_WebSocket") {
        obj = exportWebsocketInterface[target]
    } else if (func === "Event_TCP") {
        obj = exportTCPInterface[target]
    } else if (func === "Event_UDP") {
        obj = exportUDPInterface[target]
    }
    if (obj !== null && obj !== undefined) {
        return {
            contents: [{value: obj.HoverValue}, {value: "**使用代码：" + obj.Code + "**"}]
        };
    }
}

function HoverEventFuncStruct(func) {
    if (func === "Event_HTTP") {
        return {
            contents: EventInterfaceHTTP
        };
    }

    if (func === "Event_WebSocket") {
        return {
            contents: EventInterfaceWebsocket
        };
    }

    if (func === "Event_TCP") {
        return {
            contents: EventInterfaceTCP
        };
    }

    if (func === "Event_UDP") {
        return {
            contents: EventInterfaceUDP
        };
    }
}

//---------------------------------------------------------------
//HTTP接口
const exportHTTPInterface = {}
//TCP接口
const exportTCPInterface = {}
//UDP接口
const exportUDPInterface = {}
//Websocket接口
const exportWebsocketInterface = {}
//---------------------------------------------------------------
const EventInterfaceHTTP = [{value: "##  HTTP 事件接口"}]
const EventInterfaceWebsocket = [{value: "##  Websocket 事件接口"}]
const EventInterfaceTCP = [{value: "##  TCP 事件接口"}]
const EventInterfaceUDP = [{value: "##  UDP 事件接口"}]
const split = "\n\t----------------------------------\n"

function init() {
    // 发送 GET 请求
    fetch(requestEventFunc)
        .then(response => response.json()) // 将响应转换为 JSON
        .then(data => {
            //HTTP 接口1
            {
                for (let key = 0; key < data.HTTPEvent.length; key++) {
                    const obj = obj2str(data.HTTPEvent[key])
                    exportHTTPInterface[data.HTTPEvent[key].name] = obj
                    EventInterfaceHTTP.push({value: obj.HoverValue})
                    EventInterfaceHTTP.push({value: "**接口：" + obj.Code + "**"})
                }
            }
            //Websocket 接口
            {
                for (let key = 0; key < data.WebSocketEvent.length; key++) {
                    const obj = obj2str(data.WebSocketEvent[key])
                    exportWebsocketInterface[data.WebSocketEvent[key].name] = obj
                    EventInterfaceWebsocket.push({value: obj.HoverValue})
                    EventInterfaceWebsocket.push({value: "**接口：" + obj.Code + "**"})
                }
            }
            //tcp 接口
            {
                for (let key = 0; key < data.TCPEvent.length; key++) {
                    const obj = obj2str(data.TCPEvent[key])
                    exportTCPInterface[data.TCPEvent[key].name] = obj
                    EventInterfaceTCP.push({value: obj.HoverValue})
                    EventInterfaceTCP.push({value: "**接口：" + obj.Code + "**"})
                }
            }
            //udp 接口
            {
                for (let key = 0; key < data.UDPEvent.length; key++) {
                    const obj = obj2str(data.UDPEvent[key])
                    exportUDPInterface[data.UDPEvent[key].name] = obj
                    EventInterfaceUDP.push({value: obj.HoverValue})
                    EventInterfaceUDP.push({value: "**接口：" + obj.Code + "**"})
                }
            }
        })
        .catch(error => {
            // 处理错误
            console.error('There was an error fetching the data!', error);
        });
}

init()

function obj2str(EventObj) {
    let HoverValue = ""
    let Code = ""
    let title = ""
    let args = ""
    {
        const array = ("" + EventObj.comment).trim().split("\n")
        for (let i = 0; i < array.length; i++) {
            if (array[i].length > 1) {
                title = array[i].trim()
                break
            }
        }
    }
    {
        HoverValue = "```go\n";
        HoverValue += "/*\n";
        HoverValue += "\t接口说明：\n" + EventObj.comment;
        HoverValue += "\n*/\n";
        HoverValue += "/*\n";
        HoverValue += "\t接口名称：" + EventObj.name + "\n";
        if (EventObj.Args === null || EventObj.Args.length < 1) {
            HoverValue += "\t接口参数：无\n";
        } else {
            for (let i = 0; i < EventObj.Args.length; i++) {
                HoverValue += "\t参数" + (i + 1) + "：" + EventObj.Args[i].name + "\t类型:" + EventObj.Args[i].type + "\n";
            }
        }
        if (EventObj.Returns === null || EventObj.Returns.length < 1) {
            HoverValue += "\t接口返回值：无\n";
        } else {
            if (EventObj.Returns.length === 1) {
                HoverValue += "\t接口返回值：" + EventObj.Returns[0] + "\n";
            } else {
                for (let i = 0; i < EventObj.Returns.length; i++) {
                    HoverValue += "\t返回值" + (i + 1) + "：" + EventObj.Returns[i] + "\n";
                }
            }
        }
        HoverValue += "*/\n";
        Code += SunnyNetConnName + "." + EventObj.name + "(";
        let index = 0
        if (EventObj.Args !== null && EventObj.Args.length > 0) {
            for (let i = 0; i < EventObj.Args.length; i++) {
                index++
                if (i !== 0) {
                    Code += " , " + EventObj.Args[i].name;
                    args += " , " + "${" + index + ":" + EventObj.Args[i].name + "}";
                } else {
                    Code += EventObj.Args[i].name;
                    args += "${" + index + ":" + EventObj.Args[i].name + "}";
                }
            }
        }
        Code += ")";
        args += "$0";
        HoverValue += "\n```"
        HoverValue = HoverValue.replaceAll("\n*/\n/*\n", split)
    }
    return {HoverValue: HoverValue, Name: EventObj.name, Code: Code, title: title, Args: args}
}

async function postEventCodeType(code, x, name) {
    let res = {data: ""}
    await axios.post(requestEventCodeType + "x=" + x + "&name=" + name, code, {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(response => {
        res = response
    }).catch(error => {
        console.error('There was an error fetching the data!', error);
    })
    try {
        return JSON.parse(res.data)
    } catch (e) {
        return {}
    }
}

export async function InputEvent(funcName, allWords, words, target, Code, position) {
    let lastWord = ""
    const targetWord = target.toLowerCase()
    if (words.length > 0) {
        lastWord = words[words.length - 1].toLowerCase()
    }
    const list = []
    if (lastWord === SunnyNetConnName.toLowerCase()) {
        let obj = null;
        if (funcName === "Event_HTTP") {
            obj = exportHTTPInterface
        } else if (funcName === "Event_WebSocket") {
            obj = exportWebsocketInterface
        } else if (funcName === "Event_TCP") {
            obj = exportTCPInterface
        } else if (funcName === "Event_UDP") {
            obj = exportUDPInterface
        }
        if (obj === null || obj === undefined) {
            return list
        }
        for (let key in obj) {
            if (target === "." || target === SunnyNetConnName || key.toLowerCase().startsWith(targetWord)) {
                list.push({
                    label: key, // 提示的标签
                    kind: monaco.languages.CompletionItemKind.Method,
                    insertText: key + "(" + obj[key].Args + ")",
                    insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                    detail: obj[key].title,
                    documentation: {
                        value: obj[key].HoverValue + "\n**使用代码：" + obj[key].Code + "**\n\n",
                        isTrusted: true
                    }
                })
            }
        }
    } else if (SunnyNetConnName.toLowerCase().startsWith(targetWord)) {
        let detail = "undefined"
        if (funcName === "Event_HTTP") {
            detail = "HTTP 事件接口属性"
        } else if (funcName === "Event_WebSocket") {
            detail = "WebSocket 事件接口属性"
        } else if (funcName === "Event_TCP") {
            detail = "TCP 事件接口属性"
        } else if (funcName === "Event_UDP") {
            detail = "UDP 事件接口属性"
        }
        list.push({
            label: "Conn", // 提示的标签
            kind: monaco.languages.CompletionItemKind.Event,
            insertText: "Conn.$0",
            insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            detail: detail,
        })
    }
    let name = target
    if (name === ".") {
        if (words.length > 0) {
            name = words[words.length - 1].toLowerCase()
        }
        if (isHeaderType(Code, name, position.lineNumber)) {
            list.push({
                label: "Add     \t-        添加协议头", // 提示的标签
                kind: monaco.languages.CompletionItemKind.Method,
                insertText: "Add(${1:key}$0, ${2:value})",
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            })
            list.push({
                label: "Set     \t-        设置协议头", // 提示的标签
                kind: monaco.languages.CompletionItemKind.Method,
                insertText: "Set(${1:key}$0, ${2:value})",
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            })
            list.push({
                label: "Get     \t-        获取一个协议头值", // 提示的标签
                kind: monaco.languages.CompletionItemKind.Method,
                insertText: "Get(${2:key}$0)",
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            })
            list.push({
                label: "GetArray     \t-        获取多个同名的协议头值", // 提示的标签
                kind: monaco.languages.CompletionItemKind.Method,
                insertText: "GetArray(${1:key}$0)",
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            })
            list.push({
                label: "SetArray     \t-        设置多个同名的协议头值", // 提示的标签
                kind: monaco.languages.CompletionItemKind.Method,
                insertText: "SetArray(${1:key}$0, ${2:ArrayValue})",
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            })
            list.push({
                label: "Del     \t-        (删除一个协议头)", // 提示的标签
                kind: monaco.languages.CompletionItemKind.Method,
                insertText: "Del(${1:key}$0)",
                insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
            })
        }
    } else {
        for (let index = 0; index < builtCmdWords.length; index++) {
            const obj = builtCmdWords[index]
            for (let n = 0; n < obj.name.length; n++) {
                if (obj.name[n].toLowerCase().includes(targetWord)) {
                    let contents = ""
                    for (let i = 0; i < obj.contents.length; i++) {
                        contents += obj.contents[i].value + "\n\n"
                    }
                    list.push({
                        label: obj.name[n], // 提示的标签
                        kind: monaco.languages.CompletionItemKind.Method,
                        insertText: obj.insertText,
                        insertTextRules: monaco.languages.CompletionItemInsertTextRule.InsertAsSnippet,
                        detail: obj.detail,
                        documentation: {
                            value: contentsReplaceCode(contents)
                        }
                    })
                }
            }
        }
        //console.log(funcName, "-", JSON.stringify(words), "[" + target + "]", list, target, position)
    }
    return list
}

function isHeaderType(code, name, line) {
    const arr = code.split("\n")
    for (let i = 0; i < arr.length; i++) {
        if (i + 1 < line) {
            const s = (arr[i] + "").trim().replaceAll("\t", " ").replaceAll(" ", "")
            if (s.startsWith("var" + name + "=") && (s.indexOf("Header") !== -1)) {
                return true
            }
            if (s.startsWith(name + "=") && (s.indexOf("Header") !== -1 || (s.indexOf("make") !== -1))) {
                return true
            }
            if (s.startsWith(name + ":=") && (s.indexOf("Header") !== -1 || (s.indexOf("make") !== -1))) {
                return true
            }
            continue
        }
        break
    }
}

function contentsReplaceCode(code) {

// 示例用法
    const result = getMiddleString(code, '```go', '```');
    if (result !== "") {
        const result2 = result.replaceAll("\n\t", "\n").trim()
        let res = ""
        const array = result2.split("\n")
        for (let i = 0; i < array.length; i++) {
            res += "\t" + array[i].trim().replaceAll("\t", "") + "\n"
        }
        if (!res.endsWith("\n")) {
            res += "\n"
        }
        res = code.replaceAll(result, "\n" + res)
        return res.replaceAll("\n\n\n", "\n\n")
    }
    return code.replaceAll("\n\n\n", "\n\n")
}

function getMiddleString(str, startStr, endStr) {
    // 找到 startStr 的起始位置
    let startIndex = str.indexOf(startStr);
    if (startIndex === -1) {
        return ""; // 如果 startStr 不存在，返回 null
    }
    // 找到 startStr 的结束位置
    startIndex += startStr.length;

    // 找到 endStr 的起始位置
    let endIndex = str.indexOf(endStr, startIndex);
    if (endIndex === -1) {
        return ""; // 如果 endStr 不存在，返回 null
    }

    // 截取 startStr 和 endStr 之间的字符串
    return str.substring(startIndex, endIndex);
}
