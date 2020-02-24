package websocket

import (
	"bufio"
	"io"
)

/*
一、帧结构图及含义


0                   1                   2                   3
0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |            (16/64)            |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|    Extended payload length continued, if payload len == 127   |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               | Masking-key, if MASK set to 1 |
+-------------------------------+-------------------------------+
|    Masking-key (continued)    |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                   Payload Data continued ...                  :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                   Payload Data continued ...                  |
+---------------------------------------------------------------+

0Bit：
FIN 结束标识位，如果FIN为1，代表该帧为结束帧（如果一条消息过长可以将其拆分为多个帧，这时候FIN可以置为0，表示后面还有数据帧，服务器需要将该帧内容缓存起来，待所有帧都接收后再拼接到一起。控制帧不可拆分为多帧）。

1~3Bit：
RSV1~RSV3 保留标识位，以后做协议扩展时才会用到，目前该3位都为0

4~7Bit：
opcode 操作码，用于标识该帧负载的类型，如果收到了未知的操作码，则根据协议，需要断开WebSocket连接。操作码含义如下：
0x00 连续帧，浏览器的WebSocket API一般不会收到该类型的操作码
0x01 文本帧，最常用到的数据帧类别之一，表示该帧的负载是一段文本(UTF-8字符流)
0x02 二进制帧，较常用到的数据帧类别之一，表示该帧的负载是二进制数据
0x03-0x07 保留帧，留作未来非控制帧扩展使用
0x08 关闭连接控制帧，表示要断开WebSocket连接，浏览器端调用close方法会发送0x08控制帧
0x09 ping帧，用于检测端点是否可用，暂未发现浏览器可以通过何种方法发送该帧
0x0A pong帧，用于回复ping帧，暂未发现浏览器可以发送此种类型的控制帧
0x0B-0x0F 保留帧，留作未来控制帧扩展使用

8Bit：
MASK 掩码标识位，用来表明负载是否经过掩码处理，浏览器发送的数据都是经过掩码处理(浏览器自动处理，无需开发者编码)，服务器发送的帧必须不经过掩码处理。所以此处浏览器发送的帧必为1，服务器发送的帧必为0，否则应断开WebSocket连接

9~15Bit：
payload length 负载长度，单位字节如果负载长度0~125字节，则此处就是负载长度的字节数，如果负载长度在126~65535之间，则此处的值为126，16~32Bit表示负载的真实长度。如果负载长度在65536~2的64次方-1时，16~80Bit表示负载的真实长度。其中负载长度包括应用数据长度和扩展数据的长度

payload length 后面4个字节可能是掩码的key(如果掩码位是1则有这4个字节的key，否则没有)，掩码计算方法将在后面给出。

接下来就是负载的数据了，他们可能需要根据掩码的key进行编码（仅浏览器需要掩码），如果存在扩展数据，需要放在应用数据之前
*/

type Conn struct {
	rwc io.ReadWriteCloser
	br  *bufio.Reader
	bw  *bufio.Writer
}

func newConn(rwc io.ReadWriteCloser, br *bufio.Reader, bw *bufio.Writer) *Conn {
	return &Conn{rwc: rwc, br: br, bw: bw}
}

func (c *Conn) Write(msgType int, msg []byte) {
}

func (c *Conn) WriteHeader(msgType int, length int) {
}

func (c *Conn) Read() []byte {
	return nil
}

func (c *Conn) Close() {

}
