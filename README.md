# Go_Concurrency_Project

### Goal of the code
- first explain what this code does
- then probably put the original project submission, with more nuance and with more sense to be made 

### Explanation of why the project was created and the stages
- explain what our teach required, explain the stages in more depth here

### How to run the code with the necessary dependencies and such (need some mingwx86 thingy and go,)
    - you need to install Go (this is for the main code)
    - you need to install MinGW-W64 x86_64-msvcrt-posix-seh 15.1.0 (this is for the database)
    - you need to do 'go get' and then any imports that start with 'golang.org/',include the entire path as is in the import, (you may also need to do this for the SQLite driver)

### How to connect from a different machine 
- ensure you are one the same wifi, this project does not expand to being able to connect through the internet from different wifi's and different machines
- 
#### 1. Ensure the Server is listening on the right interface
- for the main of the server machine, make sure it is contains 
`bash
    listener, err:= net.Listen("tcp", ":6000")`
- tcp:6000 means you are listening on all interfaces, including your LAN IP address

#### 2. Get the Server's IP address
- On the machine running the server, you need to get the ip address

###### On Linux/macOS: 
    `bash  
        ip a`
    or
    `bash
        ifconfig`
###### On Windows:
    `bash
        ipconfig`

- You'll want to look for you IPv4 address, e.g:
        IPv4 Address: 000.000.0.000
- This is what clients on your local network will use to connect to the server from their machine

#### 3. Allow incoming connetions (firewall)
###### On Linux:
- If you're using ufw (Uncomplicated Firewall):
    ```bash
    sudo ufw allow 6000/tcp```

- Or for iptables:
    `bash
    sudo iptables -A INPUT -p tcp --dport 6000 -j ACCEPT`

###### On Windows:
- 1. Make sure port 6000 is allowed through the Windows Defender Firewall:

- 2. Open Windows Defender Firewall with Advanced Security

- 3. Inbound Rules → New Rule

- 4. Rule Type: Port

- 5. Protocol: TCP, Port: 6000

- 6. Allow the connection

- 7. Apply it for your profiles (private, public)

- 8. Name it, e.g., "Chat Server Port 6000"

#### 4. Client connects using the server's IP address

- On the client machine, in client.go, change this line:
    `go
        conn, err := net.Dial("tcp", "localhost:6000")

- To this, that way you're using the server's IP address (e.g. 000.000.0.000):
    `go 
        conn, err := net.Dial("tcp", "000.000.0.000:6000")

- Then rebuild the client app (do 'run go .') and run it on the other machine -- it should connect to the server across the network


