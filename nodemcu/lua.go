package nodemcu

const (
	recvCode = `function recv()
    local on,w,ack,nack=uart.on,uart.write,'\6','\21'
    local fd
    local function recv_block(d)
        local t,l = d:byte(1,2)
        if t ~= 1 then w(0, nack); fd:close(); return on('data') end
        if l >= 0  then fd:write(d:sub(3, l+2)); end
        if l == 0 then fd:close(); w(0, ack); return on('data') else w(0, ack) end
    end
    local function recv_name(d) d = d:gsub('%z.*', '') d:sub(1,-2) file.remove(d) fd=file.open(d, 'w') on('data', 130, recv_block, 0) w(0, ack) end
    on('data', '\0', recv_name, 0)
    w(0, 'C')
  end
function shafile(f) print(crypto.toHex(crypto.fhash('sha1', f))) end
`
	listFilesCode    = "for key,value in pairs(file.list()) do print(key,\"|\",value) end\r\n"
	hardwareInfoCode = "for key,value in pairs(node.info('hw')) do k=tostring(key) print(k, '|', tostring(value)) end\r\n"
)
