#!/usr/bin/python
# -*- coding: UTF-8 -*-
 
import smtplib
from email import encoders
from email.header import Header
from email.mime.text import MIMEText
from email.utils import parseaddr, formataddr
import sys
def send_mail(dtime,duser,dip,dhostname):
	#基础信息
	# from_addr = input("From:")
	from_addr = "1447678653@qq.com"
	password = "authent"
	#to_addr = from_addr
	to_addr = "1447678653@qq.com"
	# password = raw_input("Password:")
	# to_addr = input("To:")
	
	def _format_addr(s):
    		name, addr = parseaddr(s)
    		return formataddr((Header(name, 'utf-8').encode(), addr))
 
	smtp_server = "smtp.qq.com"
        mimetex = '您的机器:',dhostname,'，于:',dtime,'，被IP:',dip,'以账号',duser,'进行登录,请确认是否为公司员工。'
	#构造邮件
	msg = MIMEText(''.join(mimetex), 'plain', 'utf-8')
	msg['From'] = _format_addr("xr")
	msg['To'] = _format_addr("1447678653@qq.com")
	msg['Subject'] = Header("来自xr", 'utf-8').encode()
	#发送邮件
	server = smtplib.SMTP_SSL(smtp_server, 465)
	server.set_debuglevel(1)
	server.login(from_addr, password)
	server.sendmail(from_addr, [to_addr], msg.as_string())
	server.quit()
 
 
if __name__ == "__main__":
    send_mail(sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4])

