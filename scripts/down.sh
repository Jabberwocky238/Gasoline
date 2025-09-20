iptables -D FORWARD -i wg0 -j ACCEPT; 
iptables -D FORWARD -o wg0 -j ACCEPT;
iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE;