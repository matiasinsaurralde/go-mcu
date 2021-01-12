local pin = 4
local status = gpio.LOW

gpio.mode(pin, gpio.OUTPUT)
gpio.write(pin, status)

local mytimer = tmr.create()
mytimer:register(100, tmr.ALARM_AUTO, function (t)
    if status == gpio.LOW then
        status = gpio.HIGH
    else
        status = gpio.LOW
    end

    gpio.write(pin, status)
end)
mytimer:start()
