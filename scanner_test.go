// Copyright (c) 2014 Dataence, LLC. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sequence

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	sigtests = []struct {
		data, sig string
	}{
		{
			"jan 12 06:49:41 irc sshd[7034]: pam_unix(sshd:auth): authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost=218-161-81-238.hinet-ip.hinet.net  user=root",
			"%time%[%integer%]:(:):;==%integer%=%integer%====",
		},
		{
			"jan 12 06:49:42 irc sshd[7034]: failed password for root from 218.161.81.238 port 4228 ssh2",
			"%time%[%integer%]:%ipv4%%integer%",
		},
		{
			"9.26.157.45 - - [16/jan/2003:21:22:59 -0500] \"get /wssamples/ http/1.1\" 200 1576",
			"%ipv4%--[%time%]\"\"%integer%%integer%",
		},
		{
			"209.36.88.3 - - [03/may/2004:01:19:07 +0000] \"get http://npkclzicp.xihudohtd.ngm.au/abramson/eiyscmeqix.ac;jsessionid=b0l0v000u0?sid=00000000&sy=afr&kw=goldman&pb=fin&dt=selectrange&dr=0month&so=relevance&st=nw&ss=afr&sf=article&rc=00&clspage=0&docid=fin0000000r0jl000d00 http/1.0\" 200 27981",
			"%ipv4%--[%time%]\"%url%\"%integer%%integer%",
		},
		{
			"4/5/2012 17:55,172.23.1.101,1101,172.23.0.10,139, generic protocol command decode,3, [1:2100538:17] gpl netbios smb ipc$ unicode share access ,tcp ttl:128 tos:0x0 id:1643 iplen:20 dgmlen:122 df,***ap*** seq: 0xcef93f32  ack: 0xc40c0bb  n: 0xfc9c  tcplen: 20,",
			"%time%,%ipv4%,%integer%,%ipv4%,%integer%,,%integer%,[%integer%:%integer%:%integer%],:%integer%::%integer%:%integer%:%integer%,::n::%integer%,",
		},
		{
			"2012-04-05 17:54:47     local4.info     172.23.0.1      %asa-6-302015: built outbound udp connection 1315679 for outside:193.0.14.129/53 (193.0.14.129/53) to inside:172.23.0.10/64048 (10.32.0.1/52130)",
			"%time%%ipv4%:%integer%:%ipv4%/%integer%(%ipv4%/%integer%):%ipv4%/%integer%(%ipv4%/%integer%)",
		},
		{
			"may  2 19:00:02 dlfssrv sendmail[18980]: taa18980: from user daemon: size is 596, class is 0, priority is 30596, and nrcpts=1, message id is <200305021400.taa18980@dlfssrv.in.ibm.com>, relay=daemon@localhost",
			"%time%[%integer%]:::%integer%,%integer%,%integer%,=%integer%,<>,=",
		},
		{
			"jan 12 06:49:56 irc last message repeated 6 times",
			"%time%%integer%",
		},
		{
			"9.26.157.44 - - [16/jan/2003:21:22:59 -0500] \"get http://wssamples http/1.1\" 301 315",
			"%ipv4%--[%time%]\"%url%\"%integer%%integer%",
		},
		{
			"2012-04-05 17:51:26     local4.info     172.23.0.1      %asa-6-302016: teardown udp connection 1315632 for inside:172.23.0.2/514 to identity:172.23.0.1/514 duration 0:09:23 bytes 7999",
			"%time%%ipv4%:%integer%:%ipv4%/%integer%:%ipv4%/%integer%%integer%",
		},
		{
			"id=firewall time=\"2005-03-18 14:01:43\" fw=topsec priv=4 recorder=kernel type=conn policy=504 proto=tcp rule=deny src=210.82.121.91 sport=4958 dst=61.229.37.85 dport=23124 smac=00:0b:5f:b2:1d:80 dmac=00:04:c1:8b:d8:82",
			"==\"%time%\"==%integer%===%integer%===%ipv4%=%integer%=%ipv4%=%integer%=%mac%=%mac%",
		},
		{
			"mar 01 09:42:03.875 pffbisvr smtp[2424]: 334 warning: denied access to command 'ehlo vishwakstg1.msn.vishwak.net' from [209.235.210.30]",
			"%time%[%integer%]:%integer%:''[%ipv4%]",
		},
		{
			"mar 01 09:45:02.596 pffbisvr smtp[2424]: 121 statistics: duration=181.14 user=<egreetings@vishwak.com> id=zduqd sent=1440 rcvd=356 srcif=d45f49a2-b30 src=209.235.210.30/61663 cldst=192.216.179.206/25 svsrc=172.17.74.195/8423 dstif=fd3c875c-064 dst=172.17.74.52/25 op=\"to 1 recips\" arg=<vishwakstg1ojte15fo000033b4@vishwakstg1.msn.vishwak.net> result=\"250 m2004030109385301402 message accepted for delivery\" proto=smtp rule=131 (denied access to command 'ehlo vishwakstg1.msn.vishwak.net' from [209.235.210.30])",
			"%time%[%integer%]:%integer%:=%float%=<>==%integer%=%integer%==%ipv4%/%integer%=%ipv4%/%integer%=%ipv4%/%integer%==%ipv4%/%integer%=\"\"=<>=\"\"==%integer%(''[%ipv4%])",
		},
	}

	seqtests = []struct {
		data string
		seq  Sequence
	}{
		{
			"Jan 12 06:49:41 irc sshd[7034]: pam_unix(sshd:auth): authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost=218-161-81-238.hinet-ip.hinet.net  user=root", Sequence{
				Token{Type: TokenTime, Field: FieldUnknown, Value: "Jan 12 06:49:41"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "irc"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "sshd"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "7034"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "pam_unix"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "("},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "sshd"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "auth"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ")"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "authentication"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "failure"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ";"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "logname"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "uid"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "0"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "euid"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "0"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "tty"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "ssh"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "ruser"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "rhost"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "218-161-81-238.hinet-ip.hinet.net"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "user"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "root"},
			},
		},

		{
			"Jan 12 06:49:42 irc sshd[7034]: Failed password for root from 218.161.81.238 port 4228 ssh2", Sequence{
				Token{Type: TokenTime, Field: FieldUnknown, Value: "Jan 12 06:49:42"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "irc"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "sshd"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "7034"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Failed"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "password"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "for"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "root"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "from"},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "218.161.81.238"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "port"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "4228"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "ssh2"},
			},
		},

		//"Jan 13 17:25:59 jlz sshd[19322]: Accepted password for jlz from 108.61.8.124 port 56731 ssh2",
		//"Jan 12 14:44:48 irc sshd[11084]: Accepted publickey for jlz from 76.21.0.16 port 36609 ssh2",
		{
			"Jan 12 06:49:56 irc last message repeated 6 times", Sequence{
				Token{Type: TokenTime, Field: FieldUnknown, Value: "Jan 12 06:49:56"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "irc"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "last"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "message"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "repeated"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "6"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "times"},
			},
		},

		{
			"9.26.157.44 - - [16/Jan/2003:21:22:59 -0500] \"GET http://WSsamples HTTP/1.1\" 301 315", Sequence{
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "9.26.157.44"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "-"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "-"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenTime, Field: FieldUnknown, Value: "16/Jan/2003:21:22:59 -0500"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "GET"},
				Token{Type: TokenURL, Field: FieldUnknown, Value: "http://WSsamples"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "HTTP/1.1"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "301"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "315"},
			},
		},

		{
			"9.26.157.45 - - [16/Jan/2003:21:22:59 -0500] \"GET /WSsamples/ HTTP/1.1\" 200 1576", Sequence{
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "9.26.157.45"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "-"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "-"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenTime, Field: FieldUnknown, Value: "16/Jan/2003:21:22:59 -0500"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "GET"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/WSsamples/"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "HTTP/1.1"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "200"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "1576"},
			},
		},

		{
			"209.36.88.3 - - [03/May/2004:01:19:07 +0000] \"GET http://npkclzicp.xihudohtd.ngm.au/abramson/eiyscmeqix.ac;jsessionid=b0l0v000u0?sid=00000000&sy=afr&kw=goldman&pb=fin&dt=selectRange&dr=0month&so=relevance&st=nw&ss=AFR&sf=article&rc=00&clsPage=0&docID=FIN0000000R0JL000D00 HTTP/1.0\" 200 27981", Sequence{
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "209.36.88.3"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "-"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "-"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenTime, Field: FieldUnknown, Value: "03/May/2004:01:19:07 +0000"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "GET"},
				Token{Type: TokenURL, Field: FieldUnknown, Value: "http://npkclzicp.xihudohtd.ngm.au/abramson/eiyscmeqix.ac;jsessionid=b0l0v000u0?sid=00000000&sy=afr&kw=goldman&pb=fin&dt=selectRange&dr=0month&so=relevance&st=nw&ss=AFR&sf=article&rc=00&clsPage=0&docID=FIN0000000R0JL000D00"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "HTTP/1.0"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "200"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "27981"},
			},
		},

		{
			"4/5/2012 17:55,172.23.1.101,1101,172.23.0.10,139, Generic Protocol Command Decode,3, [1:2100538:17] GPL NETBIOS SMB IPC$ unicode share access ,TCP TTL:128 TOS:0x0 ID:1643 IpLen:20 DgmLen:122 DF,***AP*** Seq: 0xCEF93F32  Ack: 0xC40C0BB  n: 0xFC9C  TcpLen: 20,", Sequence{
				Token{Type: TokenTime, Field: FieldUnknown, Value: "4/5/2012 17:55"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "172.23.1.101"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "1101"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "172.23.0.10"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "139"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Generic"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Protocol"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Command"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Decode"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "3"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "1"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "2100538"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "17"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "GPL"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "NETBIOS"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "SMB"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "IPC$"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "unicode"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "share"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "access"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "TCP"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "TTL"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "128"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "TOS"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "0x0"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "ID"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "1643"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "IpLen"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "20"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "DgmLen"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "122"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "DF"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "***AP***"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Seq"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "0xCEF93F32"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Ack"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "0xC40C0BB"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "n"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "0xFC9C"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "TcpLen"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "20"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
			},
		},

		{
			"2012-04-05 17:51:26     Local4.Info     172.23.0.1      %ASA-6-302016: Teardown UDP connection 1315632 for inside:172.23.0.2/514 to identity:172.23.0.1/514 duration 0:09:23 bytes 7999", Sequence{
				Token{Type: TokenTime, Field: FieldUnknown, Value: "2012-04-05 17:51:26"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Local4.Info"},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "172.23.0.1"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "%ASA-6-302016"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Teardown"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "UDP"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "connection"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "1315632"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "for"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "inside"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "172.23.0.2"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "514"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "to"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "identity"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "172.23.0.1"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "514"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "duration"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "0:09:23"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "bytes"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "7999"},
			},
		},

		{
			"2012-04-05 17:54:47     Local4.Info     172.23.0.1      %ASA-6-302015: Built outbound UDP connection 1315679 for outside:193.0.14.129/53 (193.0.14.129/53) to inside:172.23.0.10/64048 (10.32.0.1/52130)", Sequence{
				Token{Type: TokenTime, Field: FieldUnknown, Value: "2012-04-05 17:54:47"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Local4.Info"},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "172.23.0.1"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "%ASA-6-302015"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "Built"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "outbound"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "UDP"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "connection"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "1315679"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "for"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "outside"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "193.0.14.129"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "53"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "("},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "193.0.14.129"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "53"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ")"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "to"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "inside"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "172.23.0.10"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "64048"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "("},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "10.32.0.1"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "52130"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ")"},
			},
		},

		{
			"id=firewall time=\"2005-03-18 14:01:43\" fw=TOPSEC priv=4 recorder=kernel type=conn policy=504 proto=TCP rule=deny src=210.82.121.91 sport=4958 dst=61.229.37.85 dport=23124 smac=00:0b:5f:b2:1d:80 dmac=00:04:c1:8b:d8:82", Sequence{
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "id"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "firewall"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "time"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenTime, Field: FieldUnknown, Value: "2005-03-18 14:01:43"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "fw"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "TOPSEC"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "priv"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "4"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "recorder"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "kernel"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "type"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "conn"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "policy"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "504"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "proto"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "TCP"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "rule"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "deny"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "src"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "210.82.121.91"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "sport"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "4958"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "dst"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "61.229.37.85"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "dport"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "23124"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "smac"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenMac, Field: FieldUnknown, Value: "00:0b:5f:b2:1d:80"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "dmac"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenMac, Field: FieldUnknown, Value: "00:04:c1:8b:d8:82"},
			},
		},

		{
			"mar 01 09:42:03.875 pffbisvr smtp[2424]: 334 warning: denied access to command 'ehlo vishwakstg1.msn.vishwak.net' from [209.235.210.30]", Sequence{
				Token{Type: TokenTime, Field: FieldUnknown, Value: "mar 01 09:42:03.875"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "pffbisvr"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "smtp"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "2424"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "334"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "warning"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "denied"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "access"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "to"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "command"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "'"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "ehlo vishwakstg1.msn.vishwak.net"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "'"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "from"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "209.235.210.30"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
			},
		},

		{
			"may  2 19:00:02 dlfssrv sendmail[18980]: taa18980: from user daemon: size is 596, class is 0, priority is 30596, and nrcpts=1, message id is <200305021400.taa18980@dlfssrv.in.ibm.com>, relay=daemon@localhost", Sequence{
				Token{Type: TokenTime, Field: FieldUnknown, Value: "may  2 19:00:02"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "dlfssrv"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "sendmail"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "18980"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "taa18980"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "from"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "user"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "daemon"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "size"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "is"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "596"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "class"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "is"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "0"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "priority"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "is"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "30596"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "and"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "nrcpts"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "1"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "message"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "id"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "is"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "<"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "200305021400.taa18980@dlfssrv.in.ibm.com"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ">"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ","},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "relay"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "daemon@localhost"},
			},
		},

		{
			"mar 01 09:45:02.596 pffbisvr smtp[2424]: 121 statistics: duration=181.14 user=<egreetings@vishwak.com> id=zduqd sent=1440 rcvd=356 srcif=d45f49a2-b30 src=209.235.210.30/61663 cldst=192.216.179.206/25 svsrc=172.17.74.195/8423 dstif=fd3c875c-064 dst=172.17.74.52/25 op=\"to 1 recips\" arg=<vishwakstg1ojte15fo000033b4@vishwakstg1.msn.vishwak.net> result=\"250 m2004030109385301402 message accepted for delivery\" proto=smtp rule=131 (denied access to command 'ehlo vishwakstg1.msn.vishwak.net' from [209.235.210.30])", Sequence{
				Token{Type: TokenTime, Field: FieldUnknown, Value: "mar 01 09:45:02.596"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "pffbisvr"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "smtp"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "2424"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "121"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "statistics"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "duration"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenFloat, Field: FieldUnknown, Value: "181.14"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "user"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "<"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "egreetings@vishwak.com"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ">"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "id"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "zduqd"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "sent"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "1440"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "rcvd"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "356"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "srcif"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "d45f49a2-b30"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "src"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "209.235.210.30"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "61663"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "cldst"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "192.216.179.206"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "25"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "svsrc"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "172.17.74.195"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "8423"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "dstif"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "fd3c875c-064"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "dst"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "172.17.74.52"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "/"},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "25"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "op"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "to 1 recips"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "arg"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "<"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "vishwakstg1ojte15fo000033b4@vishwakstg1.msn.vishwak.net"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ">"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "result"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "250 m2004030109385301402 message accepted for delivery"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "\""},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "proto"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "smtp"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "rule"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
				Token{Type: TokenInteger, Field: FieldUnknown, Value: "131"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "("},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "denied"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "access"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "to"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "command"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "'"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "ehlo vishwakstg1.msn.vishwak.net"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "'"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "from"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "["},
				Token{Type: TokenIPv4, Field: FieldUnknown, Value: "209.235.210.30"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: "]"},
				Token{Type: TokenLiteral, Field: FieldUnknown, Value: ")"},
			},
		},

		{
			"2015-02-11 11:04:40 H=(amoricanexpress.com) [64.20.195.132]:10246 F=<fxC4480@amoricanexpress.com> rejected RCPT <SCRUBBED@SCRUBBED.com>: Sender verify failed", Sequence{
				Token{Field: FieldUnknown, Type: TokenTime, Value: "2015-02-11 11:04:40", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "H", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "=", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "(", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "amoricanexpress.com", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ")", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "[", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenIPv4, Value: "64.20.195.132", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "]", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "10246", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "F", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "=", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "<", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "fxC4480@amoricanexpress.com", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ">", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "rejected", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "RCPT", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "<", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "SCRUBBED@SCRUBBED.com", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ">", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "Sender", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "verify", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "failed", isKey: false, isValue: false},
			},
		},

		{
			"Jan 31 21:42:59 mail postfix/anvil[14606]: statistics: max connection rate 1/60s for (smtp:5.5.5.5) at Jan 31 21:39:37", Sequence{
				Token{Field: FieldUnknown, Type: TokenTime, Value: "Jan 31 21:42:59", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "mail", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "postfix/anvil", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "[", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "14606", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "]", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "statistics", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "max", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "connection", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "rate", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "1/60s", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "for", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "(", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "smtp", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenIPv4, Value: "5.5.5.5", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ")", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "at", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenTime, Value: "Jan 31 21:39:37", isKey: false, isValue: false},
			},
		},

		{
			"Jan 31 21:42:59 mail postfix/anvil[14606]: statistics: max connection count 1 for (smtp:5.5.5.5) at Jan 31 21:39:37", Sequence{
				Token{Field: FieldUnknown, Type: TokenTime, Value: "Jan 31 21:42:59", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "mail", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "postfix/anvil", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "[", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "14606", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "]", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "statistics", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "max", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "connection", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "count", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "1", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "for", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "(", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "smtp", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenIPv4, Value: "5.5.5.5", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ")", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "at", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenTime, Value: "Jan 31 21:39:37", isKey: false, isValue: false},
			},
		},

		{
			"Jan 31 21:42:59 mail postfix/anvil[14606]: statistics: max cache size 1 at Jan 31 21:39:37", Sequence{
				Token{Field: FieldUnknown, Type: TokenTime, Value: "Jan 31 21:42:59", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "mail", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "postfix/anvil", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "[", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "14606", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "]", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "statistics", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "max", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "cache", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "size", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "1", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "at", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenTime, Value: "Jan 31 21:39:37", isKey: false, isValue: false},
			},
		},

		// relates to #2
		{
			"Feb 06 13:37:00 box sshd[4388]: Accepted publickey for cryptix from dead:beef:1234:5678:223:32ff:feb1:2e50 port 58251 ssh2: RSA de:ad:be:ef:74:a6:bb:45:45:52:71:de:b2:12:34:56", Sequence{
				Token{Field: FieldUnknown, Type: TokenTime, Value: "Feb 06 13:37:00", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "box", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "sshd", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "[", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "4388", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "]", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "Accepted", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "publickey", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "for", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "cryptix", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "from", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenIPv6, Value: "dead:beef:1234:5678:223:32ff:feb1:2e50", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "port", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "58251", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "ssh2", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "RSA", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "de:ad:be:ef:74:a6:bb:45:45:52:71:de:b2:12:34:56", isKey: false, isValue: false},
			},
		},

		// relates to #6
		{
			"2015-01-21 21:41:27 4515 [Note] - '::' resolves to '::';", Sequence{
				Token{Field: FieldUnknown, Type: TokenTime, Value: "2015-01-21 21:41:27", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "4515", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "[", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "Note", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "]", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "-", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "'", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenIPv6, Value: "::", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "'", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "resolves", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "to", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "'", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenIPv6, Value: "::", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "'", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ";", isKey: false, isValue: false},
			},
		},

		// relates to #6,
		{
			"2015-01-21 21:41:27 4515 [Note] Server socket created on IP: '::'.", Sequence{
				Token{Field: FieldUnknown, Type: TokenTime, Value: "2015-01-21 21:41:27", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenInteger, Value: "4515", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "[", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "Note", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "]", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "Server", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "socket", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "created", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "on", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "IP", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ":", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "'", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenIPv6, Value: "::", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: "'", isKey: false, isValue: false},
				Token{Field: FieldUnknown, Type: TokenLiteral, Value: ".", isKey: false, isValue: false},
			},
		},

		// {
		// 	"%msgtime% %apphost% %appname% : %srcuser% : tty = %string% ; pwd = %string% ; user = %dstuser% ; command = %method/10%", Sequence{
		// 		Token{Type: TokenTime, Field: FieldMsgTime, Value: "%msgtime%"},
		// 		Token{Type: TokenLiteral, Field: FieldAppHost, Value: "%apphost%"},
		// 		Token{Type: TokenLiteral, Field: FieldAppName, Value: "%appname%"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
		// 		Token{Type: TokenLiteral, Field: FieldSrcUser, Value: "%srcuser%"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: ":"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "tty"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "%string%"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: ";"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "pwd"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "%string%"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: ";"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "user"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
		// 		Token{Type: TokenLiteral, Field: FieldDstUser, Value: "%dstuser%"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: ";"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "command"},
		// 		Token{Type: TokenLiteral, Field: FieldUnknown, Value: "="},
		// 		Token{Type: TokenLiteral, Field: FieldMethod, Value: "%method/10%"},
		// 	},
		// },
	}

	hextests = []struct {
		data  string
		valid bool
	}{
		{"f0::1", true},
		{"f0f0::1", true},
		{"1:2:3:4:5:6:1:2", true},
		{"0:0:0:0:0:0:0:0", true},
		{"1:2:3:4:5::7:8", true},
		{"f0f0::f:1", true},
		{"f0f0:f::1", true},
		{"f0::1", true},
		{"::2:3:4", true},
		{"0:0:0:0:0:0:0:5", true},
		{"::5", true},
		{"::", true},
		{"ABC:567:0:0:8888:9999:1111:0", true},
		{"ABC:567::8888:9999:1111:0", true},
		{"ABC:567::8888:9999:1111:0 ", true}, // space at the end
		{"ABC::567::891::00", false},
		{":::00", false},
		{"00:04:c1:8b:d8:82", true},
		{"de:ad:be:ef:74:a6:bb:45:45:52:71:de:b2:12:34:56", true},
		{"00:0b:5f:b2:1d:80", true},
		{"00:04:c1:8b:d8:82", true},
		{"00:04:c1:8b:d8:82 ", true}, // space at end
		{"0:09:23 ", true},
		{"g:09:23 ", false},
		{"dead:beef:1234:5678:223:32ff:feb1:2e50", true},
		{"12345:32432:3232", false},
	}
)

func TestMessageScanHexString(t *testing.T) {
	msg := &message{}

	for _, tc := range hextests {
		var valid, stop bool

		msg.resetHexStates()

		for i, r := range tc.data {
			valid, stop = msg.hexStep(i, r)
			if stop {
				break
			}
		}
		valid = valid && msg.state.hexSuccColonsSeries < 2 && msg.state.hexMaxSuccColons < 3
		require.Equal(t, tc.valid, valid, tc.data)
	}
}

func TestGeneralScannerSignature(t *testing.T) {
	seq := make(Sequence, 0, 20)
	for _, tc := range sigtests {
		seq = seq[:0]
		seq, err := DefaultScanner.Tokenize(tc.data, seq)
		require.NoError(t, err)
		require.Equal(t, tc.sig, seq.Signature(), tc.data+"\n"+seq.PrintTokens())
	}
}

func TestGeneralScannerTokenize(t *testing.T) {
	seq := make(Sequence, 0, 20)
	for _, tc := range seqtests {
		seq = seq[:0]
		seq, err := DefaultScanner.Tokenize(tc.data, seq)
		require.NoError(t, err)
		// for i, tok := range seq {
		// 	require.Equal(t, tc.seq[i], tok)
		// }
		require.Equal(t, tc.seq, seq, tc.data+"\n"+seq.PrintTokens())
	}
}

func BenchmarkGeneralScannerOne(b *testing.B) {
	seq := make(Sequence, 0, 20)
	data := sigtests[0].data
	for i := 0; i < b.N; i++ {
		seq = seq[:0]
		DefaultScanner.Tokenize(data, seq)
	}
}
