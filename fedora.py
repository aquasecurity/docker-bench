#!/usr/bin/python3
import cgi
import subprocess
import cgitb
import random
cgitb.enable()
print("Content-type: text/html\r\n\r\n")
form = cgi.FieldStorage()
no=form.getvalue("x")
print("<center>")
print("<h1>Created Containers</h1></br>")
for i in range(int(no)):
        ports=[line.split(":")[-1] for line in subprocess.getoutput("netstat -tunlep | grep LISTEN | awk '{print $4}'")]
        while True:
                port=random.randrange(2000,65535)
                if port not in ports:
                        break
        port=str(port)
        try:
                subprocess.getoutput("sudo docker run -d -p "+port+":4200 fed_in_box4")
        except:
                print("Can't start conatiner sorry")
                break
        tag='<a href="https://13.235.104.102:'+port+'"> container'+str(i)+'</a>'
        print(tag)
        print("</br>")
print("</center>")