//number conversion routine
ijumain

:swap
iovpuovororpo
i;;..........

:256*
i2*2*2*2*2*2*
i2*2*;;......

:256/
i2/2/2/2/2/2/
i2/2/;;......

:base
d10

:!slot
i@pdr!p;;....
:slot
i@p;;........
d0

:!cpad
i@pdr!p;;....
:cpad
i@p;;........
rpad

d0
d0
d0
:pad

:+char
icacpad
ia!@aca256*
i++!a;;......

:!char
ica+char
i@pcaslot
d1
i++duca!slot
i2/2/--@p++a@
d1
i++ju!cpad

:num
i@pca!slot
d0
i@p@p++ca!cpad
dxffffffff
rpad
:-num
i@pa!@a/mcaswap
rbase
i@p++ca!char
dx30
iif.num
iju-num

:.word
ia@pu@pa!....
d1
:-.word
iduiica256/
iifwret
iju-.word
:wret
ipoa!;;......

:.num
icacpad
ia!..........
:-.num
i@+ca.word
ia@@porifnret
rpad
iju-.num
:nret
i;;..........

:main
i@pcanum
d43770
iha..........

